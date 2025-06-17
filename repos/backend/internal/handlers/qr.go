package handlers

import (
	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

func RegisterQRCodeRoutes(app *fiber.App) {
	app.Post("/generate-qr-zip", repositories.GenerateQRCodeZip)
	app.Post("/generate-qr", repositories.GenerateSingleQRCode)
}
