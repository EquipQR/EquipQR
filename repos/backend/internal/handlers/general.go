package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func RegisterGeneralRoutes(app *fiber.App) {
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		user, err := repositories.GetUserByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})

	app.Get("/business/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		business, err := repositories.GetBusinessByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "business not found"})
		}
		return c.JSON(business)
	})

	app.Get("/issue/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		issue, err := repositories.GetIssueByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "issue not found"})
		}
		return c.JSON(issue)
	})

	app.Get("/equipment/:id/issues", func(c *fiber.Ctx) error {
		equipmentID := c.Params("id")
		issues, err := repositories.GetIssuesByEquipmentID(equipmentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch issues"})
		}
		return c.JSON(issues)
	})
}
