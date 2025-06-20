package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
	"strconv"
)

type Config struct {
	App_Host           string
	App_Port           string
	SSL_CertPath       string
	SSL_KeyPath        string
	CORSAllowOrigins   string
	CORSAllowHeaders   string
	Host               string
	Port               string
	User               string
	Password           string
	Name               string
	SSLMode            string
	TimeZone           string
	JWT_Secret         string
	JWT_Expiry_Minutes int
	Cookie_Expiry_Days int
}

func LoadConfigFromEnv() Config {
	jwtSecret := getEnv("JWT_SECRET", generateRandomBase64(32))
	if jwtSecret == "" {
		log.Println("WARN: JWT_SECRET not set, generating one-time secret for this session")
	}

	jwtExpiryStr := getEnv("JWT_EXPIRY_MINUTES", "15")
	jwtExpiry, err := strconv.Atoi(jwtExpiryStr)
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRY_MINUTES: %v", err)
	}

	cookieExpiryStr := getEnv("JWT_EXPIRY_MINUTES", "15")
	cookieExpiry, err := strconv.Atoi(cookieExpiryStr)
	if err != nil {
		log.Fatalf("Invalid COOKIE_EXPIRY_DAYS: %v", err)
	}

	return Config{
		App_Host:           getEnv("APP_HOST", "0.0.0.0"),
		App_Port:           getEnv("APP_PORT", "8080"),
		SSL_CertPath:       getEnv("SSL_CERT", ""),
		SSL_KeyPath:        getEnv("SSL_KEY", ""),
		CORSAllowOrigins:   getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowHeaders:   getEnv("CORS_ALLOW_HEADERS", "Origin, Content-Type, Accept, Authorization"),
		Host:               getEnv("POSTGRES_HOST", "localhost"),
		Port:               getEnv("POSTGRES_PORT", "5432"),
		User:               getEnv("POSTGRES_USER", "postgres"),
		Password:           getEnv("POSTGRES_PASSWORD", "postgres"),
		Name:               getEnv("POSTGRES_DB", "appdb"),
		SSLMode:            getEnv("POSTGRES_SSLMODE", "disable"),
		TimeZone:           getEnv("POSTGRES_TIMEZONE", "UTC"),
		JWT_Secret:         jwtSecret,
		JWT_Expiry_Minutes: jwtExpiry,
		Cookie_Expiry_Days: cookieExpiry,
	}
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func generateRandomBase64(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("Failed to generate random JWT secret")
	}
	return base64.URLEncoding.EncodeToString(b)
}
