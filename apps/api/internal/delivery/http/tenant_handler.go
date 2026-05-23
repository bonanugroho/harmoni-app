package http

import (
	"errors"
	"fmt"
	"time"

	"harmoni-api/internal/domain/service"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

// TenantHandler handles tenant and fee HTTP endpoints.
type TenantHandler struct {
	tenantService *service.TenantService
}

// NewTenantHandler creates a new tenant handler.
func NewTenantHandler(tenantService *service.TenantService) *TenantHandler {
	return &TenantHandler{tenantService: tenantService}
}

// --- Request Types ---

// CreateTenantRequest represents the tenant creation request body.
type CreateTenantRequest struct {
	Block         string                    `json:"block"`
	UnitNumber    string                    `json:"unit_number"`
	Occupancy     string                    `json:"occupancy"`
	MonthlyFee    float64                   `json:"monthly_fee"`
	TerritoryID   string                    `json:"territory_id,omitempty"`
	MandatoryFees []CreateMandatoryFeeRequest `json:"mandatory_fees"`
}

// CreateMandatoryFeeRequest represents a mandatory fee creation request in JSON.
type CreateMandatoryFeeRequest struct {
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	EffectiveDate string  `json:"effective_date"`
}

// CreateVoluntaryFeeRequest represents a voluntary fee creation request in JSON.
type CreateVoluntaryFeeRequest struct {
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	EffectiveDate string  `json:"effective_date"`
	PaidAt        *string `json:"paid_at"`
}

// UpdateTenantRequest represents the tenant update request body.
type UpdateTenantRequest struct {
	Block       string  `json:"block"`
	UnitNumber  string  `json:"unit_number"`
	Occupancy   string  `json:"occupancy"`
	MonthlyFee  float64 `json:"monthly_fee"`
	TerritoryID string  `json:"territory_id,omitempty"`
}

// CreateFeeRequest represents a fee creation request with a type discriminator.
type CreateFeeRequest struct {
	Type          string  `json:"type"` // "mandatory" or "voluntary"
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	EffectiveDate string  `json:"effective_date"`
	PaidAt        *string `json:"paid_at"`
}

// UpdateFeeRequest represents a fee update request body.
type UpdateFeeRequest struct {
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	EffectiveDate string  `json:"effective_date"`
	PaidAt        *string `json:"paid_at"`
}

// --- Route Registration ---

// RegisterRoutes registers tenant and fee endpoints on the Fiber API group.
func (h *TenantHandler) RegisterRoutes(api fiber.Router, pasetoSvc *auth.PasetoService) {
	// Tenant CRUD (plural per D-03: /api/tenants)
	api.Get("/tenants", h.ListTenants)
	api.Post("/tenants", h.CreateTenant)
	api.Get("/tenants/:id", h.GetTenant)
	api.Put("/tenants/:id", h.UpdateTenant)
	api.Delete("/tenants/:id", h.DeleteTenant)

	// Fee sub-resources (nested under tenant)
	api.Get("/tenants/:id/fees", h.ListFees)
	api.Post("/tenants/:id/fees", h.CreateFee)
	api.Put("/tenants/:id/fees/:feeId", h.UpdateFee)
	api.Delete("/tenants/:id/fees/:feeId", h.DeleteFee)
}

// --- Handler Methods ---

// ListTenants handles GET /api/tenants.
func (h *TenantHandler) ListTenants(c *fiber.Ctx) error {
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	tenants, err := h.tenantService.List(c.Context(), claims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list tenants",
			"code":  "LIST_FAILED",
		})
	}

	return c.JSON(fiber.Map{"tenants": tenants})
}

// CreateTenant handles POST /api/tenants.
func (h *TenantHandler) CreateTenant(c *fiber.Ctx) error {
	var req CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST",
		})
	}

	// Validate required fields
	if req.Block == "" || req.UnitNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "block and unit_number are required",
			"code":  "MISSING_FIELDS",
		})
	}

	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	// Convert handler request to service request
	serviceReq, err := toServiceCreateRequest(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "INVALID_DATE_FORMAT",
		})
	}

	tenant, err := h.tenantService.Create(c.Context(), serviceReq, claims)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateBlockUnit):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Block and unit number already exist in this territory",
				"code":  "DUPLICATE_BLOCK_UNIT",
			})
		case errors.Is(err, service.ErrCrossTerritoryAccess):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cross-territory access denied",
				"code":  "CROSS_TERRITORY_DENIED",
			})
		case errors.Is(err, service.ErrMandatoryFeeRequired):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one mandatory fee is required",
				"code":  "MANDATORY_FEE_REQUIRED",
			})
		case errors.Is(err, service.ErrFeeExceedsMonthlyCap):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fee amount exceeds monthly fee cap",
				"code":  "FEE_EXCEEDS_CAP",
			})
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"code":  "CREATION_FAILED",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(tenant)
}

// GetTenant handles GET /api/tenants/:id.
func (h *TenantHandler) GetTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	tenant, err := h.tenantService.GetByID(c.Context(), id, claims)
	if err != nil {
		if errors.Is(err, service.ErrCrossTerritoryAccess) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cross-territory access denied",
				"code":  "CROSS_TERRITORY_DENIED",
			})
		}
		if errors.Is(err, service.ErrTenantNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found",
				"code":  "TENANT_NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get tenant",
			"code":  "GET_FAILED",
		})
	}

	return c.JSON(tenant)
}

// UpdateTenant handles PUT /api/tenants/:id.
func (h *TenantHandler) UpdateTenant(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST",
		})
	}

	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	// Convert to service request
	serviceReq := &service.UpdateTenantRequest{
		Block:       req.Block,
		UnitNumber:  req.UnitNumber,
		Occupancy:   req.Occupancy,
		MonthlyFee:  decimal.NewFromFloat(req.MonthlyFee),
		TerritoryID: req.TerritoryID,
	}

	tenant, err := h.tenantService.Update(c.Context(), id, serviceReq, claims)
	if err != nil {
		if errors.Is(err, service.ErrCrossTerritoryAccess) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cross-territory access denied",
				"code":  "CROSS_TERRITORY_DENIED",
			})
		}
		if errors.Is(err, service.ErrTenantNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found",
				"code":  "TENANT_NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "UPDATE_FAILED",
		})
	}

	return c.JSON(tenant)
}

// DeleteTenant handles DELETE /api/tenants/:id.
func (h *TenantHandler) DeleteTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	err := h.tenantService.Delete(c.Context(), id, claims)
	if err != nil {
		if errors.Is(err, service.ErrCrossTerritoryAccess) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cross-territory access denied",
				"code":  "CROSS_TERRITORY_DENIED",
			})
		}
		if errors.Is(err, service.ErrTenantNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found",
				"code":  "TENANT_NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete tenant",
			"code":  "DELETE_FAILED",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListFees handles GET /api/tenants/:id/fees.
func (h *TenantHandler) ListFees(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	mandatory, voluntary, err := h.tenantService.ListFees(c.Context(), tenantID, claims)
	if err != nil {
		if errors.Is(err, service.ErrCrossTerritoryAccess) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cross-territory access denied",
				"code":  "CROSS_TERRITORY_DENIED",
			})
		}
		if errors.Is(err, service.ErrTenantNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found",
				"code":  "TENANT_NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list fees",
			"code":  "LIST_FEES_FAILED",
		})
	}

	return c.JSON(fiber.Map{
		"mandatory_fees": mandatory,
		"voluntary_fees": voluntary,
	})
}

// CreateFee handles POST /api/tenants/:id/fees.
func (h *TenantHandler) CreateFee(c *fiber.Ctx) error {
	tenantID := c.Params("id")

	var req CreateFeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST",
		})
	}

	if req.Type != "mandatory" && req.Type != "voluntary" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type must be 'mandatory' or 'voluntary'",
			"code":  "INVALID_FEE_TYPE",
		})
	}

	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	if req.Type == "mandatory" {
		var paidAt *time.Time
		if req.PaidAt != nil {
			t := parseDate(*req.PaidAt)
			paidAt = &t
		}
		feeReq := &service.CreateMandatoryFeeRequest{
			Amount:        decimal.NewFromFloat(req.Amount),
			Description:   req.Description,
			EffectiveDate: parseDate(req.EffectiveDate),
			PaidAt:        paidAt,
		}
		fee, err := h.tenantService.CreateMandatoryFee(c.Context(), tenantID, feeReq, claims)
		if err != nil {
			if errors.Is(err, service.ErrCrossTerritoryAccess) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Cross-territory access denied",
					"code":  "CROSS_TERRITORY_DENIED",
				})
			}
			if errors.Is(err, service.ErrTenantNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Tenant not found",
					"code":  "TENANT_NOT_FOUND",
				})
			}
			if errors.Is(err, service.ErrFeeExceedsMonthlyCap) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Fee amount exceeds monthly fee cap",
					"code":  "FEE_EXCEEDS_CAP",
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"code":  "FEE_CREATION_FAILED",
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fee)
	}

	// Voluntary fee
	var paidAt *time.Time
	if req.PaidAt != nil {
		t := parseDate(*req.PaidAt)
		paidAt = &t
	}
	feeReq := &service.CreateVoluntaryFeeRequest{
		Amount:        decimal.NewFromFloat(req.Amount),
		Description:   req.Description,
		EffectiveDate: parseDate(req.EffectiveDate),
		PaidAt:        paidAt,
	}
	fee, err := h.tenantService.CreateVoluntaryFee(c.Context(), tenantID, feeReq, claims)
	if err != nil {
		if errors.Is(err, service.ErrCrossTerritoryAccess) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cross-territory access denied",
				"code":  "CROSS_TERRITORY_DENIED",
			})
		}
		if errors.Is(err, service.ErrTenantNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found",
				"code":  "TENANT_NOT_FOUND",
			})
		}
		if errors.Is(err, service.ErrFeeExceedsMonthlyCap) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Fee amount exceeds monthly fee cap",
				"code":  "FEE_EXCEEDS_CAP",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "FEE_CREATION_FAILED",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fee)
}

// UpdateFee handles PUT /api/tenants/:id/fees/:feeId.
func (h *TenantHandler) UpdateFee(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	feeID := c.Params("feeId")

	var req UpdateFeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST",
		})
	}

	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	serviceReq := &service.UpdateFeeRequest{
		Amount:        decimal.NewFromFloat(req.Amount),
		Description:   req.Description,
		EffectiveDate: parseDate(req.EffectiveDate),
	}
	if req.PaidAt != nil {
		t := parseDate(*req.PaidAt)
		serviceReq.PaidAt = &t
	}

	err := h.tenantService.UpdateFee(c.Context(), tenantID, feeID, serviceReq, claims)
	if err != nil {
		if errors.Is(err, service.ErrCrossTerritoryAccess) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cross-territory access denied",
				"code":  "CROSS_TERRITORY_DENIED",
			})
		}
		if errors.Is(err, service.ErrTenantNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found",
				"code":  "TENANT_NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "UPDATE_FEE_FAILED",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Fee updated successfully",
	})
}

// DeleteFee handles DELETE /api/tenants/:id/fees/:feeId.
func (h *TenantHandler) DeleteFee(c *fiber.Ctx) error {
	tenantID := c.Params("id")
	feeID := c.Params("feeId")
	claims := getUserClaims(c)
	if claims == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user context not found",
			"code":  "NO_USER_CONTEXT",
		})
	}

	err := h.tenantService.DeleteFee(c.Context(), tenantID, feeID, claims)
	if err != nil {
		if errors.Is(err, service.ErrCrossTerritoryAccess) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cross-territory access denied",
				"code":  "CROSS_TERRITORY_DENIED",
			})
		}
		if errors.Is(err, service.ErrTenantNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Tenant not found",
				"code":  "TENANT_NOT_FOUND",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  "DELETE_FEE_FAILED",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// --- Conversion Helpers ---

// toServiceCreateRequest converts a handler-level CreateTenantRequest to a service-level request.
func toServiceCreateRequest(req *CreateTenantRequest) (*service.CreateTenantRequest, error) {
	monthlyFee := decimal.NewFromFloat(req.MonthlyFee)
	serviceReq := &service.CreateTenantRequest{
		Block:       req.Block,
		UnitNumber:  req.UnitNumber,
		Occupancy:   req.Occupancy,
		MonthlyFee:  monthlyFee,
		TerritoryID: req.TerritoryID,
	}

	for _, f := range req.MandatoryFees {
		effectiveDate, err := time.Parse("2006-01-02", f.EffectiveDate)
		if err != nil {
			return nil, fmt.Errorf("invalid effective_date format: %w", err)
		}
		serviceReq.MandatoryFees = append(serviceReq.MandatoryFees, service.CreateMandatoryFeeRequest{
			Amount:        decimal.NewFromFloat(f.Amount),
			Description:   f.Description,
			EffectiveDate: effectiveDate,
		})
	}

	return serviceReq, nil
}

// parseDate parses a YYYY-MM-DD date string into time.Time.
// Returns zero time if the string is empty.
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}
