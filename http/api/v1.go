package api

import (
	"github.com/gofiber/fiber/v2"
)

func V1(fiberBackend *fiber.App) {
	fiberBackend.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
