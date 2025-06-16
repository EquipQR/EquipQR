package database

import (
	"fmt"
	"log"

	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(config utils.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host, config.User, config.Password, config.Name, config.Port, config.SSLMode, config.TimeZone,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	DB = conn

	log.Printf("Connected to PostgreSQL at %s:%s as %s", config.Host, config.Port, config.User)
}

func Migrate(models ...any) {
	if err := DB.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
