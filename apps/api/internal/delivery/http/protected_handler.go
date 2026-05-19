package http

import (
	"harmoni-api/internal/delivery/middleware"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
)

// ProtectedHandler handles routes that require authentication and authorization.
type ProtectedHandler struct {
	enforcer *auth.CasbinEnforcer
}

// NewProtectedHandler creates a new protected handler.
func NewProtectedHandler(enforcer *auth.CasbinEnforcer) *ProtectedHandler {
	return &ProtectedHandler{enforcer: enforcer}
}

// RegisterRoutes registers protected routes with the Fiber app.
// Routes are protected by the auth → casbin middleware chain.
func (h *ProtectedHandler) RegisterRoutes(app *fiber.App, pasetoSvc *auth.PasetoService) {
	// Create middleware instances
	authMW := middleware.NewAuthMiddleware(middleware.AuthMiddlewareConfig{
		PasetoService: pasetoSvc,
		PublicRoutes:  middleware.DefaultPublicRoutes(),
	})

	casbinMW := middleware.NewCasbinMiddleware(middleware.CasbinMiddlewareConfig{
		Enforcer: h.enforcer,
	})

	// Protected API group with middleware chain
	api := app.Group("/api", authMW, casbinMW)

	// User routes (territory-aware)
	api.Get("/users", h.ListUsers)
	api.Get("/users/:id", h.GetUser)
	api.Post("/users", h.CreateUser)

	// Tenant routes (singular to match Casbin policy resource names)
	api.Get("/tenant", h.ListTenants)
	api.Get("/tenant/:id", h.GetTenant)
	api.Post("/tenant", h.CreateTenant)

	// Income routes
	api.Get("/income", h.ListIncomes)
	api.Post("/income", h.CreateIncome)

	// Expenditure routes
	api.Get("/expenditure", h.ListExpenditures)
	api.Post("/expenditure", h.CreateExpenditure)

	// Report routes
	api.Get("/report", h.ListReports)
	api.Post("/report", h.CreateReport)
}

// ListUsers returns users filtered by territory for RT officers,
// or all users for RW officers.
func (h *ProtectedHandler) ListUsers(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	// In a real implementation, this would query the database
	// filtered by claims.TerritoryID for RT officers
	return c.JSON(fiber.Map{
		"users":       []string{},
		"territory":   claims.TerritoryID,
		"role":        claims.Role,
		"filter_type": getFilterType(claims.Role),
	})
}

// GetUser returns user details if the user belongs to the same territory.
func (h *ProtectedHandler) GetUser(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	userID := c.Params("id")
	return c.JSON(fiber.Map{
		"id":          userID,
		"territory":   claims.TerritoryID,
		"role":        claims.Role,
		"can_access":  true, // Casbin already verified access
	})
}

// CreateUser creates a new user (RT/RW officer only).
func (h *ProtectedHandler) CreateUser(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        "new-user-id",
		"territory": claims.TerritoryID,
		"created_by": claims.UserID,
	})
}

// ListTenants returns tenants in the user's territory.
func (h *ProtectedHandler) ListTenants(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.JSON(fiber.Map{
		"tenants":     []string{},
		"territory":   claims.TerritoryID,
		"role":        claims.Role,
		"filter_type": getFilterType(claims.Role),
	})
}

// GetTenant returns tenant details.
func (h *ProtectedHandler) GetTenant(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	return c.JSON(fiber.Map{
		"id":   tenantID,
		"name": "Sample Tenant",
	})
}

// CreateTenant creates a new tenant.
func (h *ProtectedHandler) CreateTenant(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        "new-tenant-id",
		"territory": claims.TerritoryID,
	})
}

// ListIncomes returns income records.
func (h *ProtectedHandler) ListIncomes(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.JSON(fiber.Map{
		"incomes":   []string{},
		"territory": claims.TerritoryID,
	})
}

// CreateIncome creates a new income record.
func (h *ProtectedHandler) CreateIncome(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        "new-income-id",
		"territory": claims.TerritoryID,
	})
}

// ListExpenditures returns expenditure records.
func (h *ProtectedHandler) ListExpenditures(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.JSON(fiber.Map{
		"expenditures": []string{},
		"territory":    claims.TerritoryID,
	})
}

// CreateExpenditure creates a new expenditure record.
func (h *ProtectedHandler) CreateExpenditure(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        "new-expenditure-id",
		"territory": claims.TerritoryID,
	})
}

// ListReports returns report records.
func (h *ProtectedHandler) ListReports(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.JSON(fiber.Map{
		"reports":   []string{},
		"territory": claims.TerritoryID,
	})
}

// CreateReport creates a new report.
func (h *ProtectedHandler) CreateReport(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        "new-report-id",
		"territory": claims.TerritoryID,
	})
}

// getUserClaims extracts user claims from the Fiber context.
func getUserClaims(c *fiber.Ctx) *auth.Claims {
	user := c.Locals("user")
	if user == nil {
		return nil
	}
	claims, ok := user.(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

// getFilterType returns the filter type based on role.
func getFilterType(role string) string {
	switch role {
	case "rw_officer":
		return "all_territories"
	case "rt_officer":
		return "own_territory"
	default:
		return "own_data"
	}
}
