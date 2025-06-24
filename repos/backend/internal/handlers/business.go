package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterBusinessRoutes(app *fiber.App) {
	app.Get("/api/business/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		business, err := repositories.GetBusinessByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "business not found"})
		}
		return c.JSON(business)
	})

	app.Post("/api/business", utils.ValidateBody[utils.CreateBusinessRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(utils.CreateBusinessRequest)

		business := models.Business{
			ID:           uuid.New(),
			BusinessName: req.BusinessName,
		}

		if err := repositories.CreateBusiness(&business); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create business"})
		}

		return c.Status(fiber.StatusCreated).JSON(business)
	})
}
