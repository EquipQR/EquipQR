package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/middleware"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterPendingRoutes(app *fiber.App) {
	app.Get("/api/pending/:businessID", middleware.RequireUser, getPendingJoinRequests)
	app.Post("/api/pending/approve", middleware.RequireUser, approvePendingJoin)
	app.Post("/api/pending/deny", middleware.RequireUser, denyPendingJoin)
	app.Get("/api/pending/:businessID/invite", middleware.RequireUser, GenerateInviteLinkHandler)
	app.Get("/api/invite/accept", middleware.RequireUser, AcceptInviteHandler)
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

func denyPendingJoin(c *fiber.Ctx) error {
	type DenyJoinRequest struct {
		RequestID string `json:"request_id"`
	}

	var req DenyJoinRequest
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

	if err := repositories.DenyPendingJoin(requestID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to deny request",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GenerateInviteLinkHandler(c *fiber.Ctx) error {
	businessIDParam := c.Params("businessID")
	email := c.Query("email")

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is required",
		})
	}

	businessID, err := uuid.Parse(businessIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid business ID",
		})
	}

	link, err := repositories.GenerateInviteLinkWithEmail(
		businessID,
		email,
		utils.AppConfig.JWT_Secret,
		utils.AppConfig.BaseURL,
		60, // valid for 60 minutes
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate invite link",
		})
	}

	return c.JSON(fiber.Map{
		"invite_link": link,
	})
}

func AcceptInviteHandler(c *fiber.Ctx) error {
	userID, err := utils.ValidateJWTFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "please log in to accept invite",
		})
	}

	params := utils.InviteParams{
		BusinessID: c.Query("business"),
		Token:      c.Query("token"),
		Email:      c.Query("email"),
		Expiry:     c.Query("exp"),
		Signature:  c.Query("sig"),
	}

	if params.BusinessID == "" || params.Token == "" || params.Email == "" || params.Expiry == "" || params.Signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing invite parameters",
		})
	}

	err = repositories.ProcessInvite(params, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "joined business successfully",
	})
}
