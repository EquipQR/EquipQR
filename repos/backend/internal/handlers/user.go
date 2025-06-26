package handlers

import (
	"fmt"

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
	fmt.Println("Fetching user")
	if cookie == "" {
		fmt.Println("No 'token' cookie found in request")
	} else {
		fmt.Println("Received cookie 'token':", cookie)
	}

	userID, err := utils.ValidateJWTFromCookie(c)
	if err != nil || userID == "" {
		fmt.Println("JWT validation failed:", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	fmt.Println("JWT validated, userID:", userID)

	user, err := repositories.GetUserByID(userID)
	if err != nil {
		fmt.Println("Error fetching user from DB:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	if user == nil {
		fmt.Println("No user returned for ID:", userID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	fmt.Println("Returning user:", user.ID, user.Username)
	return c.JSON(fiber.Map{
		"user": user,
	})
}
