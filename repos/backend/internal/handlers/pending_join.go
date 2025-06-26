package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
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
	type ApproveJoinRequest struct {
		RequestID string `json:"request_id"`
	}

	var req ApproveJoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	requestID, err := uuid.Parse(req.RequestID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request ID",
		})
	}

	if err := repositories.ApprovePendingJoin(requestID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to approve request",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
