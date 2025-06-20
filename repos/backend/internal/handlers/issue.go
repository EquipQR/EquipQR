package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
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

	app.Post("/issue", utils.ValidateBody[utils.CreateIssueRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(utils.CreateIssueRequest)

		// Check that Equipment and Assignee exists
		_, err := repositories.GetEquipmentByID(req.EquipmentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "equipment not found",
			})
		}
		_, err = repositories.GetUserByID(req.AssigneeID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "assignee not found",
			})
		}

		issue := models.Issue{
			Title:       req.Title,
			Description: req.Description,
			EquipmentID: req.EquipmentID,
			AssigneeID:  req.AssigneeID,
		}

		if err := repositories.CreateIssue(&issue); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "could not create issue",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(issue)
	})

}
