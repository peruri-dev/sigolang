package handler

import "github.com/gofiber/fiber/v3"

// NotFound returns custom 404 page
func NotFound(c fiber.Ctx) error {
	return c.Status(404).JSON(fiber.Map{
		"error": "not found",
	})
}
