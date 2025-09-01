package app

import (
	"fmt"
	"log/slog"
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
	//"github.com/peruri-dev/inatrace/integrations/estrace"
	//"github.com/peruri-dev/inatrace/integrations/ddtrace"
)

type App struct {
	f *fiber.App
	c *config.Config
}

func NewApp(c *config.Config) *App {
	inalog.Init(inalog.Cfg{
		Source: true,
		Tinted: !c.JsonLog,
	})
	//inalog.AddHook(estrace.ExtractTraceSpanID)
	//inalog.AddHook(ddtrace.ExtractTraceSpanID)

	return &App{
		c: c,
	}
}

func (app *App) Start() {
	f := transport.InitFiber(app.c)
	app.f = f

	svc := &service.Services{}

	dbConn, err := db.Open(app.c)
	if err != nil {
		inalog.Log().Error("Error", slog.Any("error", err))
	}
	svc.DB = dbConn

	cache, err := cache.NewCache(app.c)
	if err != nil {
		inalog.Log().Error("Error", slog.Any("error", err))
	}
	svc.Cache = cache

	svc.Resty = httpclient.InitRestyClient()

	handler.RegisterRoutes(app.f, svc)

	//tp := ddtrace.InitTracerDD()
	// OR:
	//tp := estrace.InitTracerES()

	// defer func() {
	// 	if err := tp.Shutdown(context.Background()); err != nil {
	// 		log.Printf("Error shutting down tracer provider: %v", err)
	// 	}
	// }()

	// Start your server here
	err = f.Listen(fmt.Sprintf("%s:%d", app.c.Host, app.c.Port))
	if err != nil {
		inalog.Log().Error("Error", slog.Any("error", err))
	}
}

func (app *App) Stop() {
	app.f.ShutdownWithTimeout(5 * time.Second)
}
