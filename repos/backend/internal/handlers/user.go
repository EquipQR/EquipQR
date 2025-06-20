package handlers

import (
	"time"

	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=256"`
}

func RegisterUserRoutes(app *fiber.App) {
	app.Post("/auth/login", utils.ValidateBody[LoginRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(LoginRequest)

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

		signedToken, err := utils.GenerateJWT(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate token",
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "session",
			Value:    signedToken,
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		return c.JSON(fiber.Map{
			"message": "login successful",
			"user":    user,
		})
	})

	app.Post("/auth/logout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
			Expires:  time.Now().Add(-1 * time.Hour),
		})
		return c.JSON(fiber.Map{
			"message": "logout successful",
		})
	})

	app.Post("/user", utils.ValidateBody[CreateUserRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(CreateUserRequest)

		hashedPassword, err := utils.GeneratePasswordHash(req.Password, utils.DefaultArgon2Config)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to hash password",
			})
		}

		user := &models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
			IsActive: true,
		}

		if err := repositories.CreateUser(user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create user",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "user created",
			"user":    user,
		})
	})
}
