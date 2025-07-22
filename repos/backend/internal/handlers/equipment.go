package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/middleware"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterEquipmentRoutes(app *fiber.App) {
	app.Get("/api/equipment/:id/issues", middleware.RequireUser, getEquipmentIssues)
	app.Get("/api/equipment/:id", middleware.RequireUser, getEquipmentByID)
	app.Post("/api/equipment", middleware.RequireUser, utils.ValidateBody[utils.CreateEquipmentRequest](), createEquipment)
}

func getEquipmentIssues(c *fiber.Ctx) error {
	equipmentID := c.Params("id")
	issues, err := repositories.GetIssuesByEquipmentID(equipmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch issues",
		})
	}
	return c.JSON(issues)
}

func getEquipmentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	eq, err := repositories.GetEquipmentByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "equipment not found",
		})
	}

	return c.JSON(eq)
}

func createEquipment(c *fiber.Ctx) error {
	req := c.Locals("body").(utils.CreateEquipmentRequest)

	moreFields, ok := req.MoreFields.(map[string]any)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "more_fields must be an object",
		})
	}

	equipment, err := repositories.CreateEquipmentEntry(
		req.BusinessID,
		req.Status,
		req.Type,
		req.Location,
		moreFields,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(equipment)
}
