package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/middleware"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterIssueRoutes(app *fiber.App) {
	app.Get("/api/issue/:id", middleware.RequireUser, handleGetIssue)
	app.Post("/api/issue", middleware.RequireUser, utils.ValidateBody[utils.CreateIssueRequest](), handleCreateIssue)
}

func handleGetIssue(c *fiber.Ctx) error {
	id := c.Params("id")

	issue, err := repositories.GetIssueByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "issue not found"})
	}

	return c.JSON(issue)
}

func handleCreateIssue(c *fiber.Ctx) error {
	req := c.Locals("body").(utils.CreateIssueRequest)
	user := c.Locals("user").(*models.User)
	userID := user.ID

	issue, err := repositories.CreateIssueFromRequest(req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(issue)
}
