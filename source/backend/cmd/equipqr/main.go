package main

import (
	"log"
	"path/filepath"

	"github.com/EquipQR/equipqr/backend/internal/repositories"
	"github.com/EquipQR/equipqr/backend/internal/server"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	envPath := filepath.Join("..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Warning: No .env file found at", envPath)
	}

	config := utils.LoadConfigFromEnv()
	to := "justinpitera@gmail.com"
	subject := "Test Email from EquipQR"
	body := "Hello! This is a test email from your service."

	if err := utils.SendEmail(to, subject, body); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully.")

	log.Println("Email sent successfully.")
	// utils.SendEmail("justinpitera@gmail.com", "test", "test")

	repositories.InitWebAuthn()
	server.RunServer(config)
}
