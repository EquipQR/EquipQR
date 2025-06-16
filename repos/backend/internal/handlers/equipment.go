package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func RegisterEquipmentRoutes(app *fiber.App) {
	app.Get("/equipment/:id/issues", func(c *fiber.Ctx) error {
		equipmentID := c.Params("id")
		issues, err := repositories.GetIssuesByEquipmentID(equipmentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch issues"})
		}
		return c.JSON(issues)
	})
}
