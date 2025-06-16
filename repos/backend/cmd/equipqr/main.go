package main

import (
	"log"
	"path/filepath"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/handlers"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	envPath := filepath.Join("..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Warning: No .env file found at", envPath)
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
	handlers.RegisterQRCodeRoutes(app)

	if config.SSL_CertPath == "" || config.SSL_KeyPath == "" {
		log.Fatal("SSL_CERT or SSL_KEY environment variables are not set")
	}

	address := config.App_Host + ":" + config.App_Port
	log.Fatal(app.ListenTLS(address, config.SSL_CertPath, config.SSL_KeyPath))
}
