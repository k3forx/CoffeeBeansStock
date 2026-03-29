package config

import (
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
	JWTSecret string
	LogLevel  string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "coffee_user"),
		DBPass:    getEnv("DB_PASSWORD", "coffee_pass"),
		DBName:    getEnv("DB_NAME", "coffee_stock"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

// DatabaseURL returns the PostgreSQL connection string.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		c.DBUser, c.DBPass, net.JoinHostPort(c.DBHost, c.DBPort), c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
