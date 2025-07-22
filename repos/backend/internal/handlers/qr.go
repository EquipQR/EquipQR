package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/middleware"
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func RegisterQRCodeRoutes(app *fiber.App) {
	app.Post("/api/generate-qr-zip", middleware.RequireUser, generateQRCodeZip)
	app.Post("/api/generate-qr", middleware.RequireUser, generateSingleQRCode)
}

func generateQRCodeZip(c *fiber.Ctx) error {
	type request struct {
		EquipmentIDs []string `json:"equipment_ids"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(body.EquipmentIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No equipment IDs provided",
		})
	}

	data, err := repositories.GenerateQRCodeZipBytes(body.EquipmentIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", "attachment; filename=qr_codes.zip")
	return c.Send(data)
}

func generateSingleQRCode(c *fiber.Ctx) error {
	type request struct {
		EquipmentID string `json:"equipment_id"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.EquipmentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No equipment ID provided",
		})
	}

	qrCode, filename, err := repositories.GenerateSingleQRCodeBytes(body.EquipmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "image/png")
	c.Set("Content-Disposition", "attachment; filename="+filename+".png")
	return c.Send(qrCode)
}
