package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	app.Post("/api/auth/login", utils.ValidateBody[utils.LoginRequest](), func(c *fiber.Ctx) error {
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
	})

	app.Post("/api/auth/logout", func(c *fiber.Ctx) error {
		utils.SetOrRemoveSessionCookie(c, "")
		return c.JSON(fiber.Map{
			"message": "logout successful",
		})
	})

	app.Post("/api/user", utils.ValidateBody[utils.CreateUserRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(utils.CreateUserRequest)

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

	app.Post("/api/auth/register", utils.ValidateBody[utils.CreateUserRequest](), func(c *fiber.Ctx) error {
		req := c.Locals("body").(utils.CreateUserRequest)

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

		signedToken, err := utils.GenerateJWT(user.ID.String())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate token",
			})
		}

		utils.SetOrRemoveSessionCookie(c, signedToken)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "user registered and logged in",
			"user":    user,
		})
	})

	app.Get("/api/user", func(c *fiber.Ctx) error {
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
	})
}
