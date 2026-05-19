package middleware

import (
	"strings"

	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddlewareConfig holds configuration for the authentication middleware.
type AuthMiddlewareConfig struct {
	// PasetoService is the PASETO token service for validation.
	PasetoService *auth.PasetoService
	// PublicRoutes is a list of route prefixes that bypass authentication.
	PublicRoutes []string
}

// NewAuthMiddleware creates a Fiber middleware that validates PASETO tokens
// from httpOnly cookies and sets user claims in the request context.
func NewAuthMiddleware(cfg AuthMiddlewareConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Skip authentication for public routes
		for _, prefix := range cfg.PublicRoutes {
			if strings.HasPrefix(path, prefix) {
				return c.Next()
			}
		}

		// Extract token from httpOnly cookie
		token := c.Cookies("paseto_token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
				"code":  "MISSING_TOKEN",
			})
		}

		// Validate token
		claims, err := cfg.PasetoService.ValidateToken(token)
		if err != nil {
			errMsg := err.Error()
			code := "INVALID_TOKEN"

			// Distinguish between expired and invalid tokens
			if strings.Contains(errMsg, "expired") {
				code = "TOKEN_EXPIRED"
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
				"code":  code,
			})
		}

		// Set user claims in request context
		c.Locals("user", claims)

		return c.Next()
	}
}

// DefaultPublicRoutes returns the default list of routes that bypass authentication.
func DefaultPublicRoutes() []string {
	return []string{
		"/health",
		"/auth/register",
		"/auth/login",
		"/auth/reset",
	}
}
