package middleware

import (
	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
)

// CasbinMiddlewareConfig holds configuration for the Casbin authorization middleware.
type CasbinMiddlewareConfig struct {
	// Enforcer is the Casbin enforcer for policy checks.
	Enforcer *auth.CasbinEnforcer
	// ResourceExtractor extracts the resource name from the request.
	// Defaults to extracting from the first URL path segment after /api/.
	ResourceExtractor func(c *fiber.Ctx) string
}

// NewCasbinMiddleware creates a Fiber middleware that enforces Casbin RBAC policies.
// It must be chained after the auth middleware (auth → casbin → handler).
func NewCasbinMiddleware(cfg CasbinMiddlewareConfig) fiber.Handler {
	if cfg.ResourceExtractor == nil {
		cfg.ResourceExtractor = defaultResourceExtractor
	}

	return func(c *fiber.Ctx) error {
		// Extract user claims set by auth middleware
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
				"code":  "NO_USER_CONTEXT",
			})
		}

		claims, ok := user.(*auth.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
				"code":  "INVALID_USER_CONTEXT",
			})
		}

		// Map HTTP method to action
		action := methodToAction(c.Method())

		// Extract resource from request
		resource := cfg.ResourceExtractor(c)

		// Determine territory domain
		// RW officers use "*" for all territories; others use their assigned territory
		domain := claims.TerritoryID
		if claims.Role == "rw_officer" {
			domain = "*"
		}

		// Check permission
		allowed, err := cfg.Enforcer.Enforce(claims.Role, resource, action, domain)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
				"code":  "ENFORCE_ERROR",
			})
		}

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
				"code":  "INSUFFICIENT_PERMISSIONS",
			})
		}

		return c.Next()
	}
}

// defaultResourceExtractor extracts the resource name from the URL path.
// For paths like /api/users, /api/tenants, /api/incomes, it returns the
// first segment after /api/. For other paths, returns the first path segment.
func defaultResourceExtractor(c *fiber.Ctx) string {
	path := c.Path()

	// Remove leading slash
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// Remove /api/ prefix if present
	if len(path) >= 4 && path[:4] == "api/" {
		path = path[4:]
	}

	// Extract first segment
	for i, ch := range path {
		if ch == '/' {
			return path[:i]
		}
	}

	return path
}

// TenantResourceExtractor extracts the resource name from tenant/fee URLs.
// For /api/tenants/:id/fees, returns "tenant" as the resource (singular, matching policy convention).
// For /api/fees/:feeId, returns "fee" as the resource.
func TenantResourceExtractor(c *fiber.Ctx) string {
	path := c.Path()
	// Remove leading slash
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	// Remove /api/ prefix
	if len(path) >= 4 && path[:4] == "api/" {
		path = path[4:]
	}
	// Get first segment
	var firstSeg string
	for i, ch := range path {
		if ch == '/' {
			firstSeg = path[:i]
			break
		}
	}
	if firstSeg == "" {
		firstSeg = path
	}
	// Normalize: "tenants" → "tenant", "fees" → "fee"
	switch firstSeg {
	case "tenants":
		return "tenant"
	case "fees":
		return "fee"
	default:
		return firstSeg
	}
}

// methodToAction maps HTTP methods to Casbin actions.
func methodToAction(method string) string {
	switch method {
	case "GET", "HEAD", "OPTIONS":
		return "read"
	case "POST", "PUT", "PATCH", "DELETE":
		return "write"
	default:
		return "read"
	}
}
