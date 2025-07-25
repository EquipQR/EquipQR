package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
)

type SMTP_Authentication struct {
}

type Config struct {
	// Server
	App_Host         string
	App_Port         string
	SSL_CertPath     string
	SSL_KeyPath      string
	CORSAllowOrigins string
	CORSAllowHeaders string

	// JWT
	JWT_Secret         string
	JWT_Expiry_Minutes int
	Inivte_Secret      string
	Cookie_Expiry_Days int
	BaseURL            string

	// Development
	Development_Mode               bool
	Verify_Frontend_Hash           bool
	Verify_Frontend_Hash_Frequency int

	// Frontend Dev Proxy
	Vite_Host         string
	Vite_Port         int
	Vite_SSL_CertPath string
	Vite_SSL_KeyPath  string
	Vite_Proxy_Target string

	// Database (Postgres)
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string

	// Email
	Email_Enabled                   bool
	Email_Display_Name              string
	Email_Reply_To                  string
	Email_SMTP_Enable               bool
	Email_SMTP_Address              string
	Email_SMPT_Port                 int
	Email_SMPT_Username             string
	Email_SMPT_Password             string
	Email_SMTP_Domain               string
	Email_SMTP_Authentication       string
	Email_SMTP_Enable_StartTLS_Auto bool
	Email_SMTP_TLS                  bool
	Email_Resend_API_Key            string

	// S3
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func LoadConfigFromEnv() Config {
	jwtSecret := getEnv("JWT_SECRET", generateRandomBase64(32))
	inviteSecret := getEnv("INVITE_SECRET", generateRandomBase64(32))

	if jwtSecret == "" {
		log.Println("WARN: JWT_SECRET not set, generating one-time secret for this session")
	}

	if inviteSecret == "" {
		log.Println("WARN: INVITE_SECRET not set, generating one-time secret for this session")
	}

	jwtExpiryStr := getEnv("JWT_EXPIRY_MINUTES", "15")
	jwtExpiry, err := strconv.Atoi(jwtExpiryStr)
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRY_MINUTES: %v", err)
	}

	cookieExpiryStr := getEnv("COOKIE_EXPIRY_DAYS", "15")
	cookieExpiry, err := strconv.Atoi(cookieExpiryStr)
	if err != nil {
		log.Fatalf("Invalid COOKIE_EXPIRY_DAYS: %v", err)
	}

	return Config{
		// Server
		App_Host:         getEnv("APP_HOST", "0.0.0.0"),
		App_Port:         getEnv("APP_PORT", "8080"),
		SSL_CertPath:     getEnv("SSL_CERT", ""),
		SSL_KeyPath:      getEnv("SSL_KEY", ""),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowHeaders: getEnv("CORS_ALLOW_HEADERS", "Origin, Content-Type, Accept, Authorization"),

		// Database (Postgres)
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		Name:     getEnv("POSTGRES_DB", "appdb"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		TimeZone: getEnv("POSTGRES_TIMEZONE", "UTC"),

		// JWT
		JWT_Secret:         jwtSecret,
		JWT_Expiry_Minutes: jwtExpiry,
		Inivte_Secret:      inviteSecret,
		Cookie_Expiry_Days: cookieExpiry,
		BaseURL:            getEnv("BASE_URL", "https://localhost:3000"),

		// Email
		Email_Enabled:                   getEnvBool("EMAIL_ENABLED", false),
		Email_Display_Name:              getEnv("EMAIL_DISPLAY_NAME", "App"),
		Email_Reply_To:                  getEnv("EMAIL_REPLY_TO", "noreply@example.com"),
		Email_SMTP_Enable:               getEnvBool("EMAIL_SMTP_ENABLE", false),
		Email_SMTP_Address:              getEnv("EMAIL_SMTP_ADDRESS", ""),
		Email_SMPT_Port:                 getEnvInt("EMAIL_SMPT_PORT", 587),
		Email_SMPT_Username:             getEnv("EMAIL_SMTP_USERNAME", ""),
		Email_SMPT_Password:             getEnv("EMAIL_SMTP_USERNAME", ""),
		Email_SMTP_Domain:               getEnv("EMAIL_SMTP_DOMAIN", ""),
		Email_SMTP_Authentication:       getEnv("EMAIL_SMTP_AUTHENTICATION", "plain"),
		Email_SMTP_Enable_StartTLS_Auto: getEnvBool("EMAIL_SMTP_ENABLE_STARTTLS_AUTO", true),
		Email_SMTP_TLS:                  getEnvBool("EMAIL_SMTP_TLS", true),
		Email_Resend_API_Key:            getEnv("RESEND_API_KEY", ""),

		Development_Mode:               getEnvBool("DEV_MODE", false),
		Verify_Frontend_Hash:           getEnvBool("VERIFY_FRONTEND_HASH", true),
		Verify_Frontend_Hash_Frequency: getEnvInt("VERIFY_FRONTEND_HASH_FREQUENCY", 5),

		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", ""),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", ""),
		MinioBucket:    getEnv("MINIO_BUCKET", ""),
		MinioEndpoint:  fmt.Sprintf("%s:%s", getEnv("MINIO_HOST", "localhost"), getEnv("MINIO_PORT", "9000")),
		MinioUseSSL:    getEnvBool("MINIO_USE_SSL", false),
	}

}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return parsed
}

func generateRandomBase64(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("Failed to generate random JWT secret")
	}
	return base64.URLEncoding.EncodeToString(b)
}
