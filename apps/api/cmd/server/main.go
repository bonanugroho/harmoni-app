package main

import (
	"log"
	"os"

	"harmoni-api/internal/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load and validate environment variables
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment configuration: %v", err)
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

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
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
