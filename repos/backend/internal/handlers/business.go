package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func RegisterBusinessRoutes(app *fiber.App) {
	app.Get("/business/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		business, err := repositories.GetBusinessByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "business not found"})
		}
		return c.JSON(business)
	})
}
