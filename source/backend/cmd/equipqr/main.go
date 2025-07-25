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

	repositories.InitWebAuthn()
	server.RunServer(config)
}
