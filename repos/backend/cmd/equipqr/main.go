package main

import (
	"log"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/handlers"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	if err := utils.InitPikaGenerator(1); err != nil {
		log.Fatalf("Failed to initialize Pika ID generator: %v", err)
	}

	config := utils.LoadConfigFromEnv()
	database.Init(config)

	database.Migrate(
		&models.User{},
		&models.Business{},
		&models.UserBusiness{},
		&models.Issue{},
		&models.Equipment{},
	)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.CORSAllowOrigins,
		AllowHeaders: config.CORSAllowHeaders,
	}))

	handlers.RegisterHealthRoutes(app)
	handlers.RegisterUserRoutes(app)
	handlers.RegisterEquipmentRoutes(app)
	handlers.RegisterBusinessRoutes(app)
	handlers.RegisterIssueRoutes(app)

	if config.SSL_CertPath == "" || config.SSL_KeyPath == "" {
		log.Fatal("SSL_CERT or SSL_KEY environment variables are not set")
	}

	log.Fatal(app.ListenTLS("0.0.0.0:8080", config.SSL_CertPath, config.SSL_KeyPath))
}
