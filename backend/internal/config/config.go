package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
)

var (
	ErrJWTSecretRequired      = errors.New("JWT_SECRET environment variable is required")
	ErrJWTSecretTooShort      = errors.New("JWT_SECRET must be at least 32 characters long")
	ErrJWTSecretBlocked       = errors.New("JWT_SECRET is a known weak secret and must not be used")
	ErrJWTSecretLowComplexity = errors.New("JWT_SECRET must contain at least 3 character types (uppercase, lowercase, digits, symbols)")
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
		JWTSecret: os.Getenv("JWT_SECRET"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}

	if err := validateJWTSecret(cfg.JWTSecret); err != nil {
		return nil, err
	}

	return cfg, nil
}

var blockedSecrets = []string{
	"dev-secret-change-in-production!",
	"your-256-bit-secret-key-here!!!!", // common JWT tutorial defaults
}

func validateJWTSecret(secret string) error {
	if secret == "" {
		return ErrJWTSecretRequired
	}
	if len(secret) < 32 {
		return fmt.Errorf("%w, got %d", ErrJWTSecretTooShort, len(secret))
	}

	if slices.Contains(blockedSecrets, strings.ToLower(secret)) {
		return ErrJWTSecretBlocked
	}

	var hasUpper, hasLower, hasDigit, hasOther bool
	for _, r := range secret {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasOther = true
		}
	}
	categories := 0
	for _, has := range []bool{hasUpper, hasLower, hasDigit, hasOther} {
		if has {
			categories++
		}
	}
	if categories < 3 {
		return ErrJWTSecretLowComplexity
	}

	return nil
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
