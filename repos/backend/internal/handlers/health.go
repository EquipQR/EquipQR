package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterHealthRoutes(app *fiber.App) {
	app.Get("/api/status", healthStatus)
}

func healthStatus(c *fiber.Ctx) error {
	return c.SendString("Status: Healthy")
}
