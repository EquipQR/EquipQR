package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
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

		// Branch A: Joining existing business
		if req.BusinessID != "" {
			business, err := repositories.GetBusinessByID(req.BusinessID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid business ID"})
			}

			if err := repositories.CreatePendingJoinRequest(user.ID, business.ID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to request business approval"})
			}

			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
				"message": "registration submitted and pending business approval",
				"user":    user,
			})
		}

		// Branch B: Creating new business
		if req.BusinessName == "" || req.BusinessType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "business name and type required for new business creation",
			})
		}

		business := &models.Business{
			BusinessName:    req.BusinessName,
			BusinessEmail:   req.BusinessEmail,
			Phone:           req.Phone,
			CountryCode:     req.CountryCode,
			Type:            req.BusinessType,
			CompanySize:     req.CompanySize,
			Country:         req.Country,
			UserCanRegister: false,
			LoginMethods:    pq.StringArray{"password"},
		}

		if err := repositories.CreateBusiness(business); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create business",
			})
		}

		if err := repositories.AddUserToBusiness(user.ID.String(), business.ID.String(), true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to assign user to business",
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
			"message":  "user registered and business created",
			"user":     user,
			"business": business,
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
