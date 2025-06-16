package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterHealthRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Status: Healthy")
	})
}
