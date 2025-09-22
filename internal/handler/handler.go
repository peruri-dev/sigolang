package handler

import (
	"fmt"
	"log/slog"
	"net/url"
	"runtime/debug"
	"sigolang/config"
	"sigolang/internal/middleware"
	"sigolang/internal/service"
	"sigolang/lib/errs"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"github.com/peruri-dev/inalog"
	"github.com/pkg/errors"

	fiberotel "github.com/gofiber/contrib/v3/otel"
	slogfiber "github.com/samber/slog-fiber"

	mCORS "github.com/gofiber/fiber/v3/middleware/cors"
	mHelmet "github.com/gofiber/fiber/v3/middleware/helmet"
	mLogger "github.com/gofiber/fiber/v3/middleware/logger"
	mRecover "github.com/gofiber/fiber/v3/middleware/recover"
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

func ParseAllowDomains(c *config.Config) []string {
	allowDomains := []string{}
	for _, s := range c.CorsAllowDomains {
		s = strings.Trim(s, " ")
		if s != "" {
			sp := strings.Split(s, ":")
			allowDomains = append(allowDomains, sp[0])
		}
	}

	if len(allowDomains) > 0 {
		slog.Info("CORS allow domain", slog.Any("domains", allowDomains))
	}

	return allowDomains
}

func RegisterRoutes(f *fiber.App, svc service.AllServices) huma.API {
	c := config.Get()
	allowDomains := ParseAllowDomains(c)

	if c.FiberLog {
		mlogCfg := mLogger.Config{
			Format:     "[Fiber] " + mLogger.ConfigDefault.Format,
			TimeFormat: time.RFC3339Nano,
		}

		if c.JsonLog {
			mlogCfg.Format = `{"reporter":"fiber","time":"${time}","level":"INFO","status":${status},"latency":"${latency}","ip":"${ip}","method":"${method}","path":"${path}","error":"${error}"}` + "\n"
			mlogCfg.DisableColors = true
		}

		f.Use(mLogger.New(mlogCfg))
	}

	fiberOpts := []fiberotel.Option{}

	// Configure OpenTelemetry to use Datadog RUM propagator along with standard propagators
	f.Use(fiberotel.New(fiberOpts...))

	f.Use(mRecover.New(mRecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx fiber.Ctx, e interface{}) {
			if err, ok := e.(error); ok {
				fmt.Println(errs.ErrToStack(errors.WithStack(err)))

				inalog.LogWith(inalog.WithCfg{Ctx: ctx.Context()}).Error(err.Error(), errs.ErrorField(errors.WithStack(err)))
			} else {
				ctx.Locals("panic", true)
				fmt.Printf("%s", debug.Stack())

				slog.Error("panic", slog.Any("e", e))

			}
		},
	}))

	f.Use(slogfiber.New(slog.Default()))
	// f.Use(inalog.NewFiberMiddleware()) for V2

	f.Use(mCORS.New(mCORS.Config{
		AllowOriginsFunc: func(origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			for _, d := range allowDomains {
				if strings.HasSuffix(u.Hostname(), d) {
					return true
				}
			}
			return false
		},
		AllowCredentials: true,
	}))
	f.Use(mHelmet.New(mHelmet.Config{
		XFrameOptions:         "DENY",
		ReferrerPolicy:        "no-referrer",
		PermissionPolicy:      "fullscreen=(), camera=(), microphone=()",
		ContentSecurityPolicy: "frame-ancestors 'self';",
		XSSProtection:         "1; mode=block",
	}))

	fiberSessionMiddleware := middleware.NewSessionFiberMiddleware()
	f.Use(fiberSessionMiddleware)

	cfg := huma.DefaultConfig(c.ServiceName, c.AppVersion)
	cfg.FieldsOptionalByDefault = true

	if !c.ApiDocs || !c.IsDevelopment() {
		cfg.DocsPath = ""
		cfg.SchemasPath = ""
		cfg.OpenAPIPath = ""
	}

	cfg.CreateHooks = nil
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
