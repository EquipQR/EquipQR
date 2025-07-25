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
	"github.com/EquipQR/equipqr/backend/internal/s3"
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
		&models.IssueAttachment{},
		&models.Equipment{},
	)

	log.Println("🔐 Initializing MinIO...")
	s3.Init()
	log.Println("✅ MinIO ready")

	app := fiber.New(fiber.Config{
		AppName:               "EquipQR",
		ServerHeader:          "EquipQR-Server",
		DisableStartupMessage: true,
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
	handlers.RegisterMediaRoutes(app)

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

	if config.Development_Mode && config.Verify_Frontend_Hash {
		log.Println("[DEBUG] Continuously verifying frontend hash...")
		go startFrontendHashChecker(config)
	}

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

var (
	lastFrontendHash      string
	lastFrontendError     string
	mismatchWarnedAlready bool
	mismatchReminderTimer time.Time
)

func checkFrontendHash() {
	srcPath := "../frontend/src"
	hashPath := "./web/.frontend_build_hash"

	hash, err := calculateDirectoryHash(srcPath)
	if err != nil {
		msg := fmt.Sprintf("⚠️  Could not hash frontend source: %v", err)
		if msg != lastFrontendError {
			log.Println(msg)
			lastFrontendError = msg
		}
		return
	}

	stored, err := os.ReadFile(hashPath)
	if err != nil {
		msg := fmt.Sprintf("⚠️  Could not read frontend build hash from %s: %v", hashPath, err)
		if msg != lastFrontendError {
			log.Println(msg)
			lastFrontendError = msg
		}
		return
	}

	storedHash := string(bytesTrimSpace(stored))

	if storedHash != hash {
		if !mismatchWarnedAlready {
			fmt.Println()
			fmt.Println("🔁  Frontend build hash mismatch detected!")
			fmt.Printf("📦  Stored Hash:    %s\n", color.New(color.FgHiRed).Sprint(storedHash))
			fmt.Printf("📁  Current Source: %s\n", color.New(color.FgHiGreen).Sprint(hash))
			fmt.Println("⚠️   Rebuild the frontend to match the current source.")
			fmt.Println()
			mismatchWarnedAlready = true
			mismatchReminderTimer = time.Now()
		} else if time.Since(mismatchReminderTimer) > 30*time.Second {
			log.Println("🔁  Frontend still mismatched — waiting on rebuild.")
			mismatchReminderTimer = time.Now()
		}
		lastFrontendHash = "" // invalidate match cache
	} else if lastFrontendHash != hash {
		fmt.Printf("✅  Frontend hash verified: %s\n", color.New(color.FgGreen).Sprint(hash))
		lastFrontendHash = hash
		lastFrontendError = ""
		mismatchWarnedAlready = false
	}
}

func startFrontendHashChecker(config utils.Config) {
	interval := time.Duration(config.Verify_Frontend_Hash_Frequency) * time.Second

	for {
		time.Sleep(interval)
		checkFrontendHash()
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
	dim := color.New(color.FgHiBlack).SprintFunc()

	log.Println(bold("🚀 EquipQR server is starting up..."))
	log.Println(dim("────────────────────────────────────"))

	printConfigSection("▸ Server", map[string]string{
		"Host":     fmt.Sprintf("%s:%s", config.App_Host, config.App_Port),
		"TLS Cert": config.SSL_CertPath,
		"TLS Key":  config.SSL_KeyPath,
	}, nil)

	printConfigSection("▸ Database", map[string]string{
		"Host":      config.Host,
		"Port":      config.Port,
		"Name":      config.Name,
		"User":      config.User,
		"SSL Mode":  config.SSLMode,
		"Time Zone": config.TimeZone,
	}, nil)

	printConfigSection("▸ Auth", map[string]string{
		"JWT Expiry":    fmt.Sprintf("%d min", config.JWT_Expiry_Minutes),
		"Cookie Expiry": fmt.Sprintf("%d days", config.Cookie_Expiry_Days),
	}, nil)

	printConfigSection("▸ CORS", map[string]string{
		"Allow Origins": config.CORSAllowOrigins,
		"Allow Headers": config.CORSAllowHeaders,
	}, nil)

	printConfigSection("▸ Email", map[string]string{
		"Enabled":       fmt.Sprintf("%t", config.Email_Enabled),
		"Display Name":  config.Email_Display_Name,
		"Reply-To":      config.Email_Reply_To,
		"SMTP Enabled":  fmt.Sprintf("%t", config.Email_SMTP_Enable),
		"SMTP Address":  config.Email_SMTP_Address,
		"SMTP Port":     fmt.Sprintf("%d", config.Email_SMPT_Port),
		"SMTP Username": config.Email_SMPT_Username,
		"SMTP Password": config.Email_SMPT_Password,
		"SMTP Domain":   config.Email_SMTP_Domain,
		"Auth Method":   config.Email_SMTP_Authentication,
		"StartTLS Auto": fmt.Sprintf("%t", config.Email_SMTP_Enable_StartTLS_Auto),
		"TLS Required":  fmt.Sprintf("%t", config.Email_SMTP_TLS),
	}, map[string]bool{
		"SMTP Username": true,
		"SMTP Password": true,
	})

	printConfigSection("▸ Development", map[string]string{
		"Development Mode":           fmt.Sprintf("%t", config.Development_Mode),
		"Verify Frontend Hash":       fmt.Sprintf("%t", config.Verify_Frontend_Hash),
		"Hash Check Frequency (sec)": fmt.Sprintf("%d", config.Verify_Frontend_Hash_Frequency),
	}, nil)

	printConfigSection("▸ Minio", map[string]string{
		"Endpoint": config.MinioEndpoint,
		"Bucket":   config.MinioBucket,
	}, nil)

	log.Println(dim("────────────────────────────────────"))
}

func printConfigSection(title string, settings map[string]string, redactedKeys map[string]bool) {
	section := color.New(color.FgCyan).SprintFunc()
	key := color.New(color.FgHiBlack).SprintFunc()
	value := color.New(color.FgGreen).SprintFunc()

	fmt.Println(section(title))
	for k, v := range settings {
		if redactedKeys != nil && redactedKeys[k] {
			v = redactValue(k, v)
		}
		fmt.Printf("   %s %s\n", key(k+":"), value(v))
	}
}

func redactValue(key string, val string) string {
	if val == "" {
		return ""
	}

	switch key {
	case "SMTP Username":
		parts := strings.Split(val, "@")
		if len(parts) != 2 {
			return "[REDACTED]"
		}
		username := parts[0]
		if len(username) > 3 {
			username = username[:3] + "..."
		}
		return fmt.Sprintf("%s@%s", username, parts[1])

	default:
		return "[REDACTED]"
	}
}
