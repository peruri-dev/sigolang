package handler

import (
	"log/slog"
	"os"
	"sigolang/config"
	"sigolang/frontend"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
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
		//f.Static("/", "./dist")
	} else {
		viteHandler, err = vite.NewHandler(vite.Config{
			FS:      os.DirFS("."),
			IsDev:   true,
			ViteURL: "http://localhost:5173",
		})

		if err != nil {
			slog.Error("error setup vite in dev", slog.Any("error", err))
			return false
		}

		slog.Info("Serving frontend vite-dev at /")

		f.Static("/src/assets", "./src/assets")
	}

	f.Use(adaptor.HTTPHandler(viteHandler))
	return true
}
