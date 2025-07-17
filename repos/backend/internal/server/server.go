package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	"github.com/fatih/color"
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

	app := fiber.New(fiber.Config{
		AppName:               "EquipQR",
		ServerHeader:          "EquipQR-Server",
		DisableStartupMessage: true, // ← disables Fiber's built-in log
	})

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
		fmt.Println()
		fmt.Println("🔁  Frontend build hash mismatch detected!")
		fmt.Printf("📦  Stored Hash:    %s\n", color.New(color.FgHiRed).Sprint(storedHash))
		fmt.Printf("📁  Current Source: %s\n", color.New(color.FgHiGreen).Sprint(hash))
		fmt.Println("⚠️   Rebuild the frontend to match the current source.")
		fmt.Println()
	} else {
		fmt.Printf("✅  Frontend hash verified: %s\n", color.New(color.FgGreen).Sprint(hash))
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
	bold := color.New(color.FgWhite, color.Bold).SprintFunc()
	section := color.New(color.FgCyan).SprintFunc()
	key := color.New(color.FgHiBlack).SprintFunc()
	value := color.New(color.FgGreen).SprintFunc()
	dim := color.New(color.FgHiBlack).SprintFunc()

	log.Println(bold("🚀 EquipQR server is starting up..."))
	log.Println(dim("────────────────────────────────────"))

	fmt.Println(section("▸ Server"))
	fmt.Printf("   %s %s:%s\n", key("Host:         "), value(config.App_Host), value(config.App_Port))
	fmt.Printf("   %s %s\n", key("TLS Cert:     "), value(config.SSL_CertPath))
	fmt.Printf("   %s %s\n", key("TLS Key:      "), value(config.SSL_KeyPath))

	fmt.Println(section("▸ Database"))
	fmt.Printf("   %s %s\n", key("Host:         "), value(config.Host))
	fmt.Printf("   %s %s\n", key("Port:         "), value(config.Port))
	fmt.Printf("   %s %s\n", key("Name:         "), value(config.Name))
	fmt.Printf("   %s %s\n", key("User:         "), value(config.User))
	fmt.Printf("   %s %s\n", key("SSL Mode:     "), value(config.SSLMode))
	fmt.Printf("   %s %s\n", key("Time Zone:    "), value(config.TimeZone))

	fmt.Println(section("▸ Auth"))
	fmt.Printf("   %s %d min\n", key("JWT Expiry:   "), config.JWT_Expiry_Minutes)
	fmt.Printf("   %s %d days\n", key("Cookie Expiry:"), config.Cookie_Expiry_Days)

	fmt.Println(section("▸ CORS"))
	fmt.Printf("   %s %s\n", key("Allow Origins:"), value(config.CORSAllowOrigins))
	fmt.Printf("   %s %s\n", key("Allow Headers:"), value(config.CORSAllowHeaders))

	log.Println(dim("────────────────────────────────────"))
}
