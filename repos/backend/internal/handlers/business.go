package handlers

import (
	"strconv"

	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/middleware"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RegisterBusinessRoutes(app *fiber.App) {
	app.Get("/api/business/:id", middleware.RequireUser, getBusinessByID)
	app.Get("/api/businesses", middleware.RequireUser, listBusinessesPaginated)
	app.Post("/api/business", middleware.RequireUser, utils.ValidateBody[utils.CreateBusinessRequest](), createBusiness)
}

func getBusinessByID(c *fiber.Ctx) error {
	id := c.Params("id")

	business, err := repositories.GetBusinessByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "business not found",
		})
	}

	memberCount, err := repositories.CountBusinessMembers(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to count members",
		})
	}

	return c.JSON(fiber.Map{
		"id":              business.ID,
		"businessName":    business.BusinessName,
		"businessType":    business.Type,
		"userCanRegister": business.UserCanRegister,
		"loginMethods":    business.LoginMethods,
		"memberCount":     memberCount,
	})
}

func listBusinessesPaginated(c *fiber.Ctx) error {
	pageParam := c.Query("page", "1")
	limitParam := c.Query("limit", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid page number",
		})
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid limit number",
		})
	}

	offset := (page - 1) * limit

	businesses, err := repositories.GetBusinessesPaginated(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch businesses",
		})
	}

	return c.JSON(businesses)
}

func createBusiness(c *fiber.Ctx) error {
	req := c.Locals("body").(utils.CreateBusinessRequest)

	business := models.Business{
		ID:           uuid.New(),
		BusinessName: req.BusinessName,
	}

	if err := repositories.CreateBusiness(&business); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not create business",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(business)
}
