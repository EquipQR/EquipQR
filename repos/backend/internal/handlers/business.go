package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateBusinessRequest struct {
	BusinessName string `json:"business_name" validate:"required,min=2,max=64"`
}

func RegisterBusinessRoutes(app *fiber.App) {
	app.Get("/business/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		business, err := repositories.GetBusinessByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "business not found"})
		}
		return c.JSON(business)
	})

	app.Post("/business", utils.ValidateBody[CreateBusinessRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(CreateBusinessRequest)

		business := models.Business{
			ID:           uuid.NewString(),
			BusinessName: req.BusinessName,
		}

		if err := repositories.CreateBusiness(&business); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create business"})
		}

		return c.Status(fiber.StatusCreated).JSON(business)
	})
}
