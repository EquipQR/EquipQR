package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterWebAuthnRoutes(app *fiber.App) {
	webAuthn := app.Group("/api/auth/webauthn")

	webAuthn.Post("/register/begin", beginRegistration)
	webAuthn.Post("/register/finish", finishRegistration)

	webAuthn.Post("/login/begin", beginLogin)
	webAuthn.Post("/login/finish", finishLogin)
}

func beginRegistration(c *fiber.Ctx) error {
	userID, err := utils.ValidateJWTFromCookie(c)
	if err != nil || userID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	print("hello")
	options, err := repositories.BeginWebAuthnRegistration(userID, c)
	print("hello")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to begin registration",
		})
	}
	print("hello")

	return c.JSON(options)
}

func finishRegistration(c *fiber.Ctx) error {
	if err := repositories.FinishWebAuthnRegistration(c); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "registration failed"})
	}
	return c.JSON(fiber.Map{"success": true})
}

func beginLogin(c *fiber.Ctx) error {
	type body struct {
		Email string `json:"email"`
	}
	var input body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	options, err := repositories.BeginWebAuthnLogin(input.Email, c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "login failed"})
	}

	return c.JSON(options)
}

func finishLogin(c *fiber.Ctx) error {
	token, err := repositories.FinishWebAuthnLogin(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "authentication failed"})
	}

	return c.JSON(fiber.Map{"token": token})
}
