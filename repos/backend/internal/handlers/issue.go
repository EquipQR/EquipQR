package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterIssueRoutes(app *fiber.App) {
	app.Get("/api/issue/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		issue, err := repositories.GetIssueByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "issue not found"})
		}
		return c.JSON(issue)
	})

	app.Post("/api/issue", utils.ValidateBody[utils.CreateIssueRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(utils.CreateIssueRequest)

		equipmentID, err := uuid.Parse(req.EquipmentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid equipment_id",
			})
		}

		assigneeID, err := uuid.Parse(req.AssigneeID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid assignee_id",
			})
		}

		issue := models.Issue{
			Title:       req.Title,
			Description: req.Description,
			EquipmentID: equipmentID,
			AssigneeID:  assigneeID,
		}

		if err := repositories.CreateIssue(&issue); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "could not create issue",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(issue)
	})

}
