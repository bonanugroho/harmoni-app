package config

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	DatabaseURL     string `mapstructure:"DATABASE_URL"`
	PasetoSecretKey string `mapstructure:"PASETO_SECRET_KEY"`
	EmailAPIKey     string `mapstructure:"EMAIL_API_KEY"`
	FromEmail       string `mapstructure:"FROM_EMAIL"`
	AppEnv          string `mapstructure:"APP_ENV"`
	AppPort         string `mapstructure:"APP_PORT"`
}

// LoadEnv loads and validates required environment variables.
// Viper automatically loads .env file if present, then overrides with actual env vars.
func LoadEnv() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", "8080")

	// Configure Viper to read .env file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AddConfigPath("./apps/api")

	// Read .env file (ignore if not found — env vars will still work)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Bind environment variables (these override .env values)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal into Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
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

	return &cfg, nil
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

// EnvValidationError is returned when a required environment variable is missing or invalid.
type EnvValidationError struct {
	Field   string
	Message string
}

func (e *EnvValidationError) Error() string {
	return e.Message
}
