package main

import (
	"log"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := database.LoadConfigFromEnv()
	database.Init(config)

	database.Migrate(
		&models.User{},
		&models.Business{},
		&models.UserBusiness{},
		&models.Issue{},
		&models.Equipment{},
	)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Status: Healthy")
	})

	// ✅ Attach equipment routes
	handlers.RegisterGeneralRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
