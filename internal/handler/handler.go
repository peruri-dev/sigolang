package handler

import (
	"fmt"
	"log/slog"
	"runtime/debug"
	"sigolang/config"
	"sigolang/internal/service"
	"sigolang/lib/errs"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/peruri-dev/inalog"
	"github.com/pkg/errors"

	mCORS "github.com/gofiber/fiber/v2/middleware/cors"
	mHelmet "github.com/gofiber/fiber/v2/middleware/helmet"
	mRecover "github.com/gofiber/fiber/v2/middleware/recover"
	mRequestId "github.com/gofiber/fiber/v2/middleware/requestid"
)

const (
	BEARER_AUTH = "bearerAuth"
	COOKIE_AUTH = "cookieAuth"
)

type Handler struct {
	svc service.AllServices
}

func UnwrapFiberUserContextMiddleware(ctx huma.Context, next func(huma.Context)) {
	ctx = huma.WithContext(ctx, inalog.WithFiberCtx(ctx.Context()))
	next(ctx)
}

func RegisterRoutes(f *fiber.App, svc service.AllServices) huma.API {
	c := config.Get()
	f.Use(otelfiber.Middleware())
	f.Use(mCORS.New())
	f.Use(mHelmet.New(mHelmet.Config{
		XFrameOptions:         "DENY",
		ReferrerPolicy:        "no-referrer",
		PermissionPolicy:      "fullscreen=(), camera=(), microphone=()",
		ContentSecurityPolicy: "frame-ancestors 'self';",
		XSSProtection:         "1; mode=block",
	}))

	f.Use(mRequestId.New(mRequestId.Config{ContextKey: inalog.CtxKeyRequestID}))
	f.Use(inalog.NewFiberMiddleware())
	f.Use(mRecover.New(mRecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e any) {
			if err, ok := e.(error); ok {
				fmt.Println(errs.ErrToStack(errors.WithStack(err)))

				inalog.LogWith(inalog.WithCfg{Ctx: inalog.WithFiberCtx(ctx.Context())}).Error(err.Error(), errs.ErrorField(errors.WithStack(err)))
			} else {
				ctx.Locals("panic", true)
				fmt.Printf("%s", debug.Stack())

				slog.Error("panic", slog.Any("e", e))
			}
		},
	}))

	cfg := huma.DefaultConfig(c.ServiceName, c.AppVersion)
	if c.IsProduction() {
		cfg.DocsPath = ""
	}
	cfg.Servers = []*huma.Server{
		{URL: c.PublishUrl},
	}

	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		BEARER_AUTH: {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
		COOKIE_AUTH: {
			Type: "apiKey",
			In:   "cookie",
			Name: "session",
		},
	}

	api := humafiber.New(f, cfg)
	api.UseMiddleware(UnwrapFiberUserContextMiddleware)

	h := &Handler{
		svc,
	}
	h.RoutesStatus(api)
	h.RoutesUser(api)

	feEnabled := RegisterFrontend(f)
	if !feEnabled {
		f.Use(NotFound)
	}

	return api
}
