package middleware

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func RequireUser(c *fiber.Ctx) error {
	cookie := c.Cookies("session")
	if cookie == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	userID, err := utils.ValidateJWTFromCookie(c)
	if err != nil || userID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	user, err := repositories.GetUserByID(userID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	c.Locals("user", user)
	return c.Next()
}
