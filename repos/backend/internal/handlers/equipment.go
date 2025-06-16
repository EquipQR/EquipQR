package handlers

import (
	"encoding/json"

	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

type CreateEquipmentRequest struct {
	BusinessID string `json:"business_id" validate:"required"`
	Status     string `json:"status" validate:"required,oneof='in service' 'not in service'"`
	Type       string `json:"type" validate:"required"`
	Location   string `json:"location"`
	MoreFields any    `json:"more_fields"`
}

func RegisterEquipmentRoutes(app *fiber.App) {
	app.Get("/equipment/:id/issues", func(c *fiber.Ctx) error {
		equipmentID := c.Params("id")
		issues, err := repositories.GetIssuesByEquipmentID(equipmentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch issues"})
		}
		return c.JSON(issues)
	})

	app.Get("/equipment/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		eq, err := repositories.GetEquipmentByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "equipment not found",
			})
		}

		return c.JSON(eq)
	})

	app.Post("/equipment", utils.ValidateBody[CreateEquipmentRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(CreateEquipmentRequest)

		_, err := repositories.GetBusinessByID(req.BusinessID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "business not found",
			})
		}

		moreFieldsJSON, err := json.Marshal(req.MoreFields)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid more_fields format",
			})
		}

		equipment := models.Equipment{
			BusinessID: req.BusinessID,
			Status:     req.Status,
			Type:       req.Type,
			Location:   req.Location,
			MoreFields: datatypes.JSON(moreFieldsJSON),
		}

		if err := repositories.CreateEquipment(&equipment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "could not create equipment",
			})
		}

		created, err := repositories.GetEquipmentByID(equipment.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "created but failed to load full record",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(created)
	})
}
