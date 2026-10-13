package transport

import (
	"sigolang/config"

	"github.com/gofiber/fiber/v3"
)

func InitFiber(c *config.Config) *fiber.App {

	f := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,

		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Check if the error is a panic recovered by the middleware
			//if recErr, ok := err.(*fiber.Error); ok && recErr.Code == fiber.StatusInternalServerError {}

			// Customize the JSON response for panics
			if wasPanic, ok := c.Locals("panic").(bool); ok {
				if wasPanic {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error":   "Internal Server Error",
						"message": "An unexpected error occurred. Please try again later.",
					})
				}
			}
			// For other errors, use Fiber's default error handling or a different custom logic
			return fiber.DefaultErrorHandler(c, err)
		},
	})

	return f
}
