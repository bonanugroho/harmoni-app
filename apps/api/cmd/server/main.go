package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"harmoni-api/internal/config"
	"harmoni-api/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load and validate environment variables
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment configuration: %v", err)
	}

	// Establish database connection pool
	ctx := context.Background()
	db, err := database.NewConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run pending migrations on startup
	// Use working directory (where the binary is executed from) to find migrations
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	migrationsPath := filepath.Join(wd, "migrations")
	if err := database.RunMigrations(cfg.DatabaseURL, migrationsPath); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Health check endpoint with database status
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "connected"
		if err := db.HealthCheck(c.Context()); err != nil {
			dbStatus = "disconnected"
			log.Printf("Health check: database connection failed: %v", err)
		}
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": dbStatus,
		})
	})

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Harmoni API server on port %s (%s)", port, cfg.AppEnv)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
