package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sigolang/config"
	"sigolang/internal/handler"
	"sigolang/internal/service"
	"sigolang/lib/cache"
	"sigolang/lib/db"
	"sigolang/lib/httpclient"
	"sigolang/lib/transport"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/peruri-dev/inalog"
	"github.com/peruri-dev/inalog/integrations/logtint"
	"github.com/peruri-dev/inatrace/integrations/uptrace"
)

type App struct {
	f *fiber.App
	c *config.Config
}

func NewApp(c *config.Config) *App {
	os.Setenv("INALOG_LOG_LEVEL", c.AppVersion)
	os.Setenv("INALOG_SERVICE_NAME", "sigolang")
	os.Setenv("INALOG_SERVICE_ENV", c.Env)
	os.Setenv("INALOG_SERVICE_VERSION", c.AppVersion)

	cfg := inalog.Cfg{
		Source:  true,
		TextLog: !c.JsonLog,
	}
	if !c.JsonLog {
		cfg.CustomFunc = logtint.CreateTintHandler()
	}
	inalog.Init(cfg)
	inalog.AddHook(uptrace.ExtractTraceSpanID)

	return &App{
		c: c,
	}
}

func (app *App) Start() {
	f := transport.InitFiber(app.c)
	app.f = f

	svc := &service.Services{}

	dbConn, err := db.Open(&app.c.DB)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	svc.DB = dbConn

	cache, err := cache.NewCache(&app.c.Cache)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	svc.Cache = cache

	svc.Resty = httpclient.InitRestyClient()

	handler.RegisterRoutes(app.f, svc)

	tracer := uptrace.InitTracer("inagov-be", app.c.AppVersion)

	defer func() {
		if err := tracer.Shutdown(context.Background()); err != nil {
			slog.Error("Error shutting down tracer provider", slog.Any("error", err))
		}
	}()

	// Start your server here
	err = f.Listen(fmt.Sprintf("%s:%d", app.c.Host, app.c.Port))
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
}

func (app *App) Stop() {
	app.f.ShutdownWithTimeout(5 * time.Second)
}
