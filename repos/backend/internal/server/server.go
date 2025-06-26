package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/handlers"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RunServer(config utils.Config) {
	checkFrontendHash()

	printStartupBanner(config)

	database.Init(config)

	database.Migrate(
		&models.User{},
		&models.Business{},
		&models.UserBusiness{},
		&models.PendingJoinRequest{},
		&models.Credential{},
		&models.Issue{},
		&models.Equipment{},
	)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.CORSAllowOrigins,
		AllowHeaders:     config.CORSAllowHeaders,
		AllowCredentials: true,
	}))

	handlers.RegisterHealthRoutes(app)
	handlers.RegisterUserRoutes(app)
	handlers.RegisterEquipmentRoutes(app)
	handlers.RegisterPendingRoutes(app)
	handlers.RegisterBusinessRoutes(app)
	handlers.RegisterWebAuthnRoutes(app)
	handlers.RegisterIssueRoutes(app)
	handlers.RegisterQRCodeRoutes(app)

	app.Static("/", "./web")

	app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./web/index.html")
	})

	if config.SSL_CertPath == "" || config.SSL_KeyPath == "" {
		log.Fatal("❌ SSL_CERT or SSL_KEY environment variables are not set")
	}

	address := config.App_Host + ":" + config.App_Port

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.ListenTLS(address, config.SSL_CertPath, config.SSL_KeyPath); err != nil {
			log.Printf("❌ Fiber ListenTLS error: %v\n", err)
		}
	}()

	log.Printf("✅ Server is running on https://%s\n", address)

	<-quit
	log.Println("🛑 Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("❌ Error shutting down server: %v\n", err)
	}

	database.Close()
	log.Println("✅ Server gracefully stopped")
}

func checkFrontendHash() {
	srcPath := "../frontend/src"
	hashPath := "./web/.frontend_build_hash"

	hash, err := calculateDirectoryHash(srcPath)
	if err != nil {
		log.Printf("⚠️  Could not hash frontend source: %v\n", err)
		return
	}

	stored, err := os.ReadFile(hashPath)
	if err != nil {
		log.Printf("⚠️  Could not read frontend build hash from %s: %v\n", hashPath, err)
		return
	}

	storedHash := string(bytesTrimSpace(stored))
	if storedHash != hash {
		log.Println("⚠️  Frontend source has changed since last build. Rebuild the frontend!")
	}
}

func calculateDirectoryHash(root string) (string, error) {
	hasher := sha256.New()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		hasher.Write(content)
		return nil
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func bytesTrimSpace(b []byte) string {
	return strings.TrimSpace(string(b))
}

func printStartupBanner(config utils.Config) {
	log.Println("🚀 Starting EquipQR server with the following configuration:")
	printConfig := struct {
		AppHost          string `json:"AppHost"`
		AppPort          string `json:"AppPort"`
		SSL_CertPath     string `json:"SSL_CertPath"`
		SSL_KeyPath      string `json:"SSL_KeyPath"`
		PostgresHost     string `json:"PostgresHost"`
		PostgresPort     string `json:"PostgresPort"`
		PostgresUser     string `json:"PostgresUser"`
		PostgresDB       string `json:"PostgresDB"`
		PostgresSSLMode  string `json:"PostgresSSLMode"`
		PostgresTimeZone string `json:"PostgresTimeZone"`
		CORSAllowOrigins string `json:"CORSAllowOrigins"`
		CORSAllowHeaders string `json:"CORSAllowHeaders"`
		JWTExpiryMinutes int    `json:"JWTExpiryMinutes"`
		CookieExpiryDays int    `json:"CookieExpiryDays"`
	}{
		AppHost:          config.App_Host,
		AppPort:          config.App_Port,
		SSL_CertPath:     config.SSL_CertPath,
		SSL_KeyPath:      config.SSL_KeyPath,
		PostgresHost:     config.Host,
		PostgresPort:     config.Port,
		PostgresUser:     config.User,
		PostgresDB:       config.Name,
		PostgresSSLMode:  config.SSLMode,
		PostgresTimeZone: config.TimeZone,
		CORSAllowOrigins: config.CORSAllowOrigins,
		CORSAllowHeaders: config.CORSAllowHeaders,
		JWTExpiryMinutes: config.JWT_Expiry_Minutes,
		CookieExpiryDays: config.Cookie_Expiry_Days,
	}

	configJSON, err := json.MarshalIndent(printConfig, "", "  ")
	if err != nil {
		log.Println("⚠️  Failed to print config:", err)
		return
	}

	log.Println(string(configJSON))
	log.Println("────────────────────────────────────────────────────")
}
