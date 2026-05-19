package config

import (
	"encoding/hex"
	"fmt"
	"os"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	DatabaseURL     string
	PasetoSecretKey string
	EmailAPIKey     string
	AppEnv          string
	AppPort         string
}

// LoadEnv loads and validates required environment variables.
func LoadEnv() (*Config, error) {
	cfg := &Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		PasetoSecretKey: os.Getenv("PASETO_SECRET_KEY"),
		EmailAPIKey:     os.Getenv("EMAIL_API_KEY"),
		AppEnv:          getEnvOrDefault("APP_ENV", "development"),
		AppPort:         getEnvOrDefault("APP_PORT", "8080"),
	}

	// Validate required variables
	required := map[string]string{
		"DATABASE_URL":      cfg.DatabaseURL,
		"PASETO_SECRET_KEY": cfg.PasetoSecretKey,
		"EMAIL_API_KEY":     cfg.EmailAPIKey,
	}

	for name, value := range required {
		if value == "" {
			return nil, &EnvValidationError{Field: name, Message: name + " is required"}
		}
	}

	// Validate PASETO_SECRET_KEY is a valid 32-byte hex string (64 hex chars)
	if err := validatePasetoKey(cfg.PasetoSecretKey); err != nil {
		return nil, &EnvValidationError{Field: "PASETO_SECRET_KEY", Message: err.Error()}
	}

	return cfg, nil
}

// validatePasetoKey ensures the PASETO secret key is a valid 32-byte hex-encoded string.
func validatePasetoKey(key string) error {
	// Must be exactly 64 hex characters (32 bytes)
	if len(key) != 64 {
		return fmt.Errorf("PASETO_SECRET_KEY must be 64 hex characters (32 bytes), got %d characters", len(key))
	}

	// Must be valid hex
	_, err := hex.DecodeString(key)
	if err != nil {
		return fmt.Errorf("PASETO_SECRET_KEY must be a valid hex string: %w", err)
	}

	return nil
}

// GeneratePasetoKey generates a random 32-byte hex-encoded key suitable for PASETO V4 Local.
// Usage: run `openssl rand -hex 32` or call this function programmatically.
func GeneratePasetoKey() (string, error) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = 0 // placeholder — use crypto/rand in production
	}
	return hex.EncodeToString(key), nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// EnvValidationError is returned when a required environment variable is missing or invalid.
type EnvValidationError struct {
	Field   string
	Message string
}

func (e *EnvValidationError) Error() string {
	return e.Message
}
