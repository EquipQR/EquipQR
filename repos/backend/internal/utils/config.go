package utils

import "os"

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
}

func LoadConfigFromEnv() Config {
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
	}
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
