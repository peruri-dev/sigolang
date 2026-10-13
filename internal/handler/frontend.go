package handler

import (
	"fmt"
	"log/slog"
	"os"
	"sigolang/config"
	"sigolang/frontend"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/olivere/vite"
)

func RegisterFrontend(f *fiber.App) bool {
	c := config.Get()

	// Handle the Vite server.
	var viteHandler *vite.Handler
	var err error

	if c.ViteDist {
		viteHandler, err = vite.NewHandler(vite.Config{
			FS:    frontend.GetDist(),
			IsDev: false,
		})

		if err != nil {
			slog.Error("error setup vite in dist", slog.Any("error", err))
			return false
		}

		slog.Info("Serving frontend vite-dist")
		//f.Get("/", mStatic.New("./dist"))

		f.Get("/assets/*", adaptor.HTTPHandler(viteHandler))
		f.Get("/*", func(c fiber.Ctx) error {
			// slog.Info("route path: " + c.Route().Path + ", path: " + c.Path())
			c.RequestCtx().Request.SetRequestURI("/index.html")
			return adaptor.HTTPHandler(viteHandler)(c)
		})
	} else {
		viteHandler, err = vite.NewHandler(vite.Config{
			FS:      os.DirFS("."),
			IsDev:   true,
			ViteURL: fmt.Sprintf("http://localhost:%d", c.VitePort),
		})

		if err != nil {
			slog.Error("error setup vite in dev", slog.Any("error", err))
			return false
		}

		slog.Info(fmt.Sprintf("Serving frontend vite-dev at %d", c.VitePort))

		f.Get("/src/assets/*", static.New("./src/assets"))

		// In dev mode, use All to handle all HTTP methods (needed for HMR websockets)
		// and properly proxy all requests to the Vite dev server
		f.All("/*", func(c fiber.Ctx) error {
			//slog.Info("route path: " + c.Route().Path + ", path: " + c.Path())
			c.RequestCtx().Request.SetRequestURI("/index.html")
			return adaptor.HTTPHandler(viteHandler)(c)
		})
	}

	return true
}
