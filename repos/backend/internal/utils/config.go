package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"
)

type Config struct {
	App_Host         string
	App_Port         string
	SSL_CertPath     string
	SSL_KeyPath      string
	CORSAllowOrigins string
	CORSAllowHeaders string
	Host             string
	Port             string
	User             string
	Password         string
	Name             string
	SSLMode          string
	TimeZone         string
	JWT_Secret       string
}

func LoadConfigFromEnv() Config {
	jwtSecret := getEnv("JWT_SECRET", generateRandomBase64(32))
	if jwtSecret == "" {
		log.Println("WARN: JWT_SECRET not set, generating one-time secret for this session")
	}

	return Config{
		App_Host:         getEnv("APP_HOST", "0.0.0.0"),
		App_Port:         getEnv("APP_PORT", "8080"),
		SSL_CertPath:     getEnv("SSL_CERT", ""),
		SSL_KeyPath:      getEnv("SSL_KEY", ""),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowHeaders: getEnv("CORS_ALLOW_HEADERS", "Origin, Content-Type, Accept, Authorization"),
		Host:             getEnv("POSTGRES_HOST", "localhost"),
		Port:             getEnv("POSTGRES_PORT", "5432"),
		User:             getEnv("POSTGRES_USER", "postgres"),
		Password:         getEnv("POSTGRES_PASSWORD", "postgres"),
		Name:             getEnv("POSTGRES_DB", "appdb"),
		SSLMode:          getEnv("POSTGRES_SSLMODE", "disable"),
		TimeZone:         getEnv("POSTGRES_TIMEZONE", "UTC"),
		JWT_Secret:       jwtSecret,
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
