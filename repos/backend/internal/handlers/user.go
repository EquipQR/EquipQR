package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	app.Post("/api/auth/login", utils.ValidateBody[utils.LoginRequest](), handleLogin)
	app.Post("/api/auth/logout", handleLogout)
	app.Post("/api/auth/register", utils.ValidateBody[utils.CreateUserRequest](), handleRegister)
	app.Get("/api/user", handleGetUser)
}

func handleLogin(c *fiber.Ctx) error {
	req := c.Locals("body").(utils.LoginRequest)

	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	match, err := utils.ComparePasswordHash(user.Password, req.Password, utils.DefaultArgon2Config)
	if err != nil || !match {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	signedToken, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	utils.SetOrRemoveSessionCookie(c, signedToken)

	return c.JSON(fiber.Map{
		"message": "login successful",
		"user":    user,
	})
}

func handleLogout(c *fiber.Ctx) error {
	utils.SetOrRemoveSessionCookie(c, "")
	return c.JSON(fiber.Map{
		"message": "logout successful",
	})
}

func handleRegister(c *fiber.Ctx) error {
	req := c.Locals("body").(utils.CreateUserRequest)

	user, business, token, err := repositories.RegisterNewUser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.SetOrRemoveSessionCookie(c, token)

	status := fiber.StatusCreated
	message := "user registered and business created"
	if business == nil {
		status = fiber.StatusAccepted
		message = "registration submitted and pending business approval"
	}

	return c.Status(status).JSON(fiber.Map{
		"message":  message,
		"user":     user,
		"business": business,
	})
}

func handleGetUser(c *fiber.Ctx) error {
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

	return c.JSON(fiber.Map{
		"user": user,
	})
}
