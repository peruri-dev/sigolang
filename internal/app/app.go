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
	"sigolang/lib/cache/redisc"
	"sigolang/lib/db"
	"sigolang/lib/httpclient"
	"sigolang/lib/storage/s3minio"
	"sigolang/lib/transport"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/peruri-dev/inalog"
	"github.com/peruri-dev/inalog/integrations/logtint"
	"github.com/peruri-dev/inatrace"
	"github.com/peruri-dev/inatrace/integrations/ddtrace"

	inatraceModels "github.com/peruri-dev/inatrace/models"
)

func init() {
	s3minio.RegisterMinioStorage()
	redisc.RegisterRedisCache()
}

type App struct {
	f      *fiber.App
	c      *config.Config
	svc    *service.Services
	tracer inatraceModels.CustomTraceProvider
}

func NewApp(c *config.Config) *App {
	os.Setenv("INALOG_SERVICE_ENV", c.Env)
	os.Setenv("INALOG_SERVICE_VERSION", c.AppVersion)

	cfg := inalog.Cfg{
		Source:     true,
		TextLog:    !c.JsonLog,
		MessageKey: c.MsgLog,
	}
	if !c.JsonLog {
		cfg.CustomFunc = logtint.CreateTintHandler()
	}
	inalog.Init(cfg)
	inalog.AddHook(inatrace.ExtractTraceSpanID)

	return &App{
		c: c,
	}
}

func (app *App) Init() {
	// Init Infra
	f := transport.InitFiber(app.c)
	app.f = f

	// Setup Database
	dbConn, err := db.Open(&app.c.DB)
	if err != nil {
		slog.Error("Cannot connect to DB", slog.String("DB", app.c.DB.DatabaseUri), slog.Any("error", err))
	} else if dbConn == nil {
		slog.Warn("Not using DB")
	}

	// Setup Cache
	cache, err := cache.NewCache(&app.c.Cache)
	if err != nil {
		slog.Error("Cannot connect to Redis", slog.Any("error", err))
	} else if cache == nil {
		slog.Warn("Not using Redis")
	}

	svc := &service.Services{
		DB:    dbConn,
		Cache: cache,
		Resty: httpclient.InitRestyClient(),
	}

	app.svc = svc
}

func (app *App) Routes() huma.API {
	// Init Routers
	api := handler.RegisterRoutes(app.f, app.svc)

	return api
}

func (app *App) setupTracingEnv() {
	if app.c.DDService == "" || app.c.OTELExporterOTLPEndpoint == "" {
		os.Setenv("DD_TRACE_ENABLED", "false")
		slog.Info("datadog tracing disabled (DD_SERVICE or OTEL_EXPORTER_OTLP_ENDPOINT empty)")
		return
	}
	// os.Setenv("DD_SERVICE", app.c.DDService)
	// os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", app.c.OTELExporterOTLPEndpoint)
	os.Setenv("DD_TRACE_ENABLED", "true")
	slog.Info("datadog tracing enabled",
		slog.String("service", app.c.DDService),
		slog.String("endpoint", app.c.OTELExporterOTLPEndpoint))

	if app.c.UptraceDSN != "" {
		os.Unsetenv("UPTRACE_DISABLED")
		slog.Info("enable uptrace", slog.String("dsn", app.c.UptraceDSN))
	} else {
		os.Setenv("UPTRACE_DISABLED", "yes")
	}
}

func (app *App) Start() {
	app.Init()

	app.setupTracingEnv()

	var tracer inatraceModels.CustomTraceProvider

	tracer = ddtrace.InitTracerDD()

	app.tracer = tracer

	app.Routes()
	if err := app.startServer(); err != nil {
		// humacli calls no Stop when the server fails to start, so the tracer is flushed here.
		app.shutdownTracer()
	}
}

func (app *App) Stop() {
	app.gracefulShutdown()
}

func (app *App) gracefulShutdown() {
	// Fiber first: requests still in flight enqueue activity events, so the workers have to outlive them.
	if err := app.f.ShutdownWithTimeout(5 * time.Second); err != nil {
		slog.Error("Error during graceful shutdown", slog.String("error", err.Error()))
	}

	// The tracer goes last, so in-flight requests and the worker drain still export their spans.
	app.shutdownTracer()
	slog.Info("Server shutdown completed")
}

func (app *App) shutdownTracer() {
	if app.tracer == nil {
		return
	}
	if err := app.tracer.Shutdown(context.Background()); err != nil {
		slog.Error("Error shutting down tracer provider", slog.Any("error", err))
	}
}

func (app *App) startServer() error {
	err := app.f.Listen(fmt.Sprintf("%s:%d", app.c.Host, app.c.Port),
		fiber.ListenConfig{DisableStartupMessage: !app.c.StartupMessage},
	)
	if err != nil {
		slog.Error("Error", slog.Any("error", err))
	}
	return err
}
