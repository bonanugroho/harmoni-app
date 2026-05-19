package config

import "os"

// Config holds the application configuration loaded from environment variables.
type Config struct {
	DatabaseURL    string
	PasetoSecretKey string
	EmailAPIKey    string
	AppEnv         string
	AppPort        string
}

// LoadEnv loads and validates required environment variables.
func LoadEnv() (*Config, error) {
	cfg := &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		PasetoSecretKey: os.Getenv("PASETO_SECRET_KEY"),
		EmailAPIKey:    os.Getenv("EMAIL_API_KEY"),
		AppEnv:         getEnvOrDefault("APP_ENV", "development"),
		AppPort:        getEnvOrDefault("APP_PORT", "8080"),
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

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// EnvValidationError is returned when a required environment variable is missing.
type EnvValidationError struct {
	Field   string
	Message string
}

func (e *EnvValidationError) Error() string {
	return e.Message
}
