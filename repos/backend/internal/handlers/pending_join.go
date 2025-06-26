package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterPendingRoutes(app *fiber.App) {
	app.Get("/api/pending/:businessID", getPendingJoinRequests)
	app.Post("/api/pending/approve", approvePendingJoin)
}

func getPendingJoinRequests(c *fiber.Ctx) error {
	businessIDStr := c.Params("businessID")
	businessID, err := uuid.Parse(businessIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid business ID",
		})
	}

	requests, err := repositories.GetAllPendingJoinsForBusiness(businessID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch pending requests",
		})
	}

	return c.JSON(requests)
}

func approvePendingJoin(c *fiber.Ctx) error {
	var req utils.ApproveJoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	businessID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid business ID",
		})
	}

	if err := repositories.ApprovePendingJoin(userID, businessID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to approve request",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
