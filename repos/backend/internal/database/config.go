package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(config utils.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host, config.User, config.Password, config.Name, config.Port, config.SSLMode, config.TimeZone,
	)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn, // show only warnings and errors
			IgnoreRecordNotFoundError: true,        // suppress 'record not found'
			Colorful:                  true,
		},
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

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

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error retrieving sql.DB from gorm.DB: %v\n", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v\n", err)
	} else {
		log.Println("Database connection closed successfully.")
	}
}
