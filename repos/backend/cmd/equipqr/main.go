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

	// Register routes
	handlers.RegisterHealthRoutes(app)
	handlers.RegisterUserRoutes(app)
	handlers.RegisterEquipmentRoutes(app)
	handlers.RegisterBusinessRoutes(app)
	handlers.RegisterIssueRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
