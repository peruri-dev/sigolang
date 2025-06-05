package transport

import (
	"sigolang/config"

	"github.com/gofiber/fiber/v2"
)

func InitFiber(c *config.Config) *fiber.App {
	f := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	return f
}
