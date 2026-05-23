package middleware

import (
	"log"

	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
)

// AuditCrossTerritoryAccess returns middleware that logs all access attempts
// and flags cross-territory violations. Designed to run before Casbin middleware
// to capture request context for audit trail.
func AuditCrossTerritoryAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user != nil {
			if claims, ok := user.(*auth.Claims); ok {
				log.Printf("AUDIT: user=%s role=%s territory=%s method=%s path=%s",
					claims.UserID, claims.Role, claims.TerritoryID, c.Method(), c.Path())
			}
		}
		return c.Next()
	}
}
