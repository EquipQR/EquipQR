package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

func LoadConfigFromEnv() Config {
	return Config{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		Name:     getEnv("POSTGRES_DB", "appdb"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		TimeZone: getEnv("POSTGRES_TIMEZONE", "UTC"),
	}
}

func Init(config Config) {
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

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
