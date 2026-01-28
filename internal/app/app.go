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
	f   *fiber.App
	c   *config.Config
	svc *service.Services
}

func NewApp(c *config.Config) *App {
	os.Setenv("INALOG_LOG_LEVEL", c.AppVersion)
	os.Setenv("INALOG_ACCESS_LOG", "DEBUG")
	os.Setenv("INALOG_SERVICE_NAME", c.ServiceName)
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

func (app *App) Init() {
	f := transport.InitFiber(app.c)
	app.f = f

	svc := &service.Services{}

	dbConn, err := db.Open(&app.c.DB)
	if err != nil {
		slog.Error("error opening database", slog.String("DB", app.c.DB.DatabaseUri), slog.Any("error", err))
	} else if dbConn == nil {
		slog.Warn("not using database")
	}
	svc.DB = dbConn

	cache, err := cache.NewCache(&app.c.Cache)
	if err != nil {
		slog.Error("error opening cache", slog.String("Cache", app.c.Cache.CacheUri), slog.Any("error", err))
	} else if cache == nil {
		slog.Warn("not using cache")
	}
	svc.Cache = cache

	svc.Resty = httpclient.InitRestyClient()

	app.svc = svc
}

func (app *App) Routes() {
	handler.RegisterRoutes(app.f, app.svc)
}

func (app *App) Start() {
	var err error
	app.Init()

	tracer := uptrace.InitTracer(app.c.ServiceName, app.c.AppVersion)

	defer func() {
		if err := tracer.Shutdown(context.Background()); err != nil {
			slog.Error("Error shutting down tracer provider", slog.Any("error", err))
		}
	}()

	// Init Routers
	app.Routes()

	// Start your server here
	err = app.f.Listen(fmt.Sprintf("%s:%d", app.c.Host, app.c.Port))
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
}

func (app *App) Stop() {
	app.f.ShutdownWithTimeout(5 * time.Second)
}
