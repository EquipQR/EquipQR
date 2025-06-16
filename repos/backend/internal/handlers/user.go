package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=256"`
}

func RegisterUserRoutes(app *fiber.App) {
	app.Get("/user/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		user, err := repositories.GetUserByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})

	app.Post("/user", utils.ValidateBody[CreateUserRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(CreateUserRequest)

		rawPassword := req.Password
		hashedPassword, err := utils.GeneratePasswordHash(rawPassword, utils.DefaultArgon2Config)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to hash password",
			})
		}

		user := models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
		}

		if err := repositories.CreateUser(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "could not create user",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	})
}
