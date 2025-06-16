package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func RegisterIssueRoutes(app *fiber.App) {
	app.Get("/issue/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		issue, err := repositories.GetIssueByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "issue not found"})
		}
		return c.JSON(issue)
	})
}
