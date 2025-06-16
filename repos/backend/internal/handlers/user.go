package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		user, err := repositories.GetUserByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})
}
