# Phase 2: Tenant & Fee Management — Pattern Map

**Mapped:** 2026-05-22
**Files analyzed:** 30 (12 new, 5 modified, 5 migration pairs, 8 test files)
**Analogs found:** 25 / 25 (100%)

## File Classification

| New/Modified File | Role | Data Flow | Closest Analog | Match Quality |
|---|---|---|---|---|
| `apps/api/internal/domain/entity/tenant.go` | entity | CRUD | `internal/domain/entity/user.go` | exact |
| `apps/api/internal/domain/repository/tenant_repository.go` | repository-interface | CRUD | `internal/domain/repository/user_repository.go` | exact |
| `apps/api/internal/domain/repository/fee_repository.go` | repository-interface | CRUD | `internal/domain/repository/user_repository.go` | role-match |
| `apps/api/internal/domain/service/tenant_service.go` | service | CRUD | `internal/domain/service/auth_service.go` | role-match |
| `apps/api/internal/infrastructure/repository/tenant_repository.go` | repository-pgx | CRUD | `internal/infrastructure/repository/user_repository.go` | exact |
| `apps/api/internal/infrastructure/repository/fee_repository.go` | repository-pgx | CRUD | `internal/infrastructure/repository/password_reset_token_repository.go` | role-match |
| `apps/api/internal/delivery/http/tenant_handler.go` | handler | request-response | `internal/delivery/http/auth_handler.go` | role-match |
| `apps/api/internal/delivery/middleware/audit.go` | middleware | request-response | `internal/delivery/middleware/casbin.go` | role-match |
| `apps/api/migrations/006_create_tenants_table.up.sql` | migration | SQL DDL | `migrations/002_create_users_table.up.sql` | exact |
| `apps/api/migrations/007_create_mandatory_fees_table.up.sql` | migration | SQL DDL | `migrations/002_create_users_table.up.sql` | role-match |
| `apps/api/migrations/008_create_voluntary_fees_table.up.sql` | migration | SQL DDL | `migrations/002_create_users_table.up.sql` | role-match |
| `apps/api/migrations/009_create_user_tenants_table.up.sql` | migration | SQL DDL | `migrations/002_create_users_table.up.sql` | role-match |
| `apps/api/migrations/010_seed_sample_tenants.up.sql` | migration | SQL DML | `migrations/005_seed_territories.up.sql` | exact |
| `apps/api/internal/delivery/http/protected_handler.go` **(modify)** | handler | request-response | Same file (remove tenant stubs) | — |
| `apps/api/cmd/server/main.go` **(modify)** | config | startup | Same file (register new handler) | — |
| `apps/api/policy.csv` **(modify)** | config | RBAC | Same file (add fee resources) | — |
| `apps/api/rbac_model.conf` **(modify)** | config | RBAC | Same file (extend matcher) | — |
| `apps/api/internal/domain/service/tenant_service_test.go` | test | unit | `internal/domain/service/auth_service_test.go` | exact |
| `apps/api/internal/domain/service/fee_service_test.go` | test | unit | `internal/domain/service/auth_service_test.go` | role-match |
| `apps/api/internal/infrastructure/repository/tenant_repository_test.go` | test | unit | `internal/infrastructure/repository/user_repository_test.go` | exact |
| `apps/api/internal/delivery/http/tenant_handler_test.go` | test | integration | `internal/delivery/http/protected_handler_test.go` | role-match |
| `apps/api/internal/delivery/middleware/audit_test.go` | test | unit | `internal/delivery/middleware/casbin_test.go` | exact |

---

## Pattern Assignments

### `apps/api/internal/domain/entity/tenant.go` — entity, CRUD

**Analog:** `apps/api/internal/domain/entity/user.go` (lines 1-33)

**Imports pattern** (lines 1-4):
```go
package entity

import "time"
```

**Core entity pattern** (lines 5-17):
```go
// Tenant represents a tenant unit in the Harmoni system.
type Tenant struct {
	ID          string    `json:"id"`
	Block       string    `json:"block"`
	UnitNumber  string    `json:"unit_number"`
	Occupancy   string    `json:"occupancy"`
	MonthlyFee  float64   `json:"monthly_fee"`
	TerritoryID string    `json:"territory_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
```

**Sanitize pattern** (lines 21-33):
```go
func (t *Tenant) Sanitize() *Tenant {
	return &Tenant{
		ID:          t.ID,
		Block:       t.Block,
		UnitNumber:  t.UnitNumber,
		Occupancy:   t.Occupancy,
		MonthlyFee:  t.MonthlyFee,
		TerritoryID: t.TerritoryID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
```

**Pattern notes:** Create three structs in one file: `Tenant`, `MandatoryFee`, `VoluntaryFee`. Each with JSON tags. Tenant has `Sanitize()` method. MonthlyFee should be `shopspring/decimal.Decimal` not `float64` (per RESEARCH.md — financial precision requirement).

---

### `apps/api/internal/domain/repository/tenant_repository.go` — repository-interface, CRUD

**Analog:** `apps/api/internal/domain/repository/user_repository.go` (lines 1-28)

**Imports pattern** (lines 1-7):
```go
package repository

import (
	"context"

	"harmoni-api/internal/domain/entity"
)
```

**Interface pattern** (lines 9-28):
```go
// TenantRepository defines the interface for tenant data access.
type TenantRepository interface {
	Create(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error)
	FindByID(ctx context.Context, id string) (*entity.Tenant, error)
	ListByTerritory(ctx context.Context, territoryID string) ([]*entity.Tenant, error)
	Update(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error)
	Delete(ctx context.Context, id string, territoryID string) error
}
```

**Interface pattern notes:**
- `CreateTx` variant needed for transactional tenant+fee creation (accepts pgx.Tx)
- `ListByUserID` variant needed for resident tenant access (joins `user_tenants`)
- Methods that enforce territory isolation take `territoryID` parameter

---

### `apps/api/internal/domain/repository/fee_repository.go` — repository-interface, CRUD

**Analog:** `apps/api/internal/domain/repository/user_repository.go` (lines 1-28)

**Additional pattern:** `apps/api/internal/domain/repository/password_reset_token_repository.go` (lines 18-31 — for non-CRUD entity)

**Imports pattern:**
```go
package repository

import (
	"context"

	"harmoni-api/internal/domain/entity"
)
```

**Fee repository interface pattern** (derived from user_repository interface + mandatory/voluntary fee split):
```go
// FeeRepository defines the interface for fee data access.
type FeeRepository interface {
	// --- Mandatory Fees ---
	CreateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error)
	CreateMandatoryTx(ctx context.Context, tx pgx.Tx, fee *entity.MandatoryFee) (*entity.MandatoryFee, error)
	ListMandatoryByTenant(ctx context.Context, tenantID string) ([]*entity.MandatoryFee, error)
	UpdateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error)
	DeleteMandatory(ctx context.Context, id string) error

	// --- Voluntary Fees ---
	CreateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error)
	ListVoluntaryByTenant(ctx context.Context, tenantID string) ([]*entity.VoluntaryFee, error)
	UpdateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error)
	DeleteVoluntary(ctx context.Context, id string) error
}
```

---

### `apps/api/internal/domain/service/tenant_service.go` — service, CRUD + validation

**Analog:** `apps/api/internal/domain/service/auth_service.go` (lines 1-219)

**Imports pattern** (lines 1-15):
```go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"harmoni-api/internal/domain/entity"
	"harmoni-api/internal/domain/repository"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)
```

**Sentinel errors pattern** (lines 17-24):
```go
var (
	ErrTenantNotFound         = errors.New("tenant not found")
	ErrDuplicateBlockUnit     = errors.New("block and unit number already exist in territory")
	ErrMandatoryFeeRequired   = errors.New("at least one mandatory fee is required")
	ErrFeeExceedsMonthlyCap   = errors.New("fee amount exceeds tenant's monthly fee cap")
	ErrInvalidEffectiveDate   = errors.New("effective date cannot be in the past")
	ErrInvalidPaidAt          = errors.New("paid_at must be after effective_date")
	ErrCrossTerritoryAccess   = errors.New("cross-territory access denied")
)
```

**Service struct pattern** (lines 26-36):
```go
// TenantService handles tenant and fee business logic.
type TenantService struct {
	tenantRepo repository.TenantRepository
	feeRepo    repository.FeeRepository
	db         *pgxpool.Pool
}

// NewTenantService creates a new tenant service.
func NewTenantService(
	tenantRepo repository.TenantRepository,
	feeRepo repository.FeeRepository,
	db *pgxpool.Pool,
) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		feeRepo:    feeRepo,
		db:         db,
	}
}
```

**Core CRUD method pattern** (from auth_service.go `Register` lines 57-92):
```go
// Create creates a new tenant with mandatory fees in a transaction.
func (s *TenantService) Create(ctx context.Context, req *CreateTenantRequest, userClaims *auth.Claims) (*entity.Tenant, error) {
	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Determine territory scope
	territoryID := userClaims.TerritoryID
	if userClaims.Role == "rw_officer" && req.TerritoryID != "" {
		territoryID = req.TerritoryID
	}

	// Create tenant
	tenant := &entity.Tenant{
		Block:       req.Block,
		UnitNumber:  req.UnitNumber,
		Occupancy:   req.Occupancy,
		MonthlyFee:  req.MonthlyFee,
		TerritoryID: territoryID,
	}
	created, err := s.tenantRepo.CreateTx(ctx, tx, tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	// Validate at least one mandatory fee provided
	if len(req.MandatoryFees) == 0 {
		return nil, ErrMandatoryFeeRequired
	}

	// Create mandatory fees
	for _, feeReq := range req.MandatoryFees {
		if err := s.validateFee(feeReq, tenant.MonthlyFee); err != nil {
			return nil, fmt.Errorf("mandatory fee validation failed: %w", err)
		}
		// persist fee ...
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return created, nil
}
```

**Validation pattern** (auth_service.go line 59 style):
```go
// validateFee checks fee constraints.
func (s *TenantService) validateFee(fee *entity.MandatoryFee, monthlyFee decimal.Decimal) error {
	if fee.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("fee amount must be non-negative")
	}
	if fee.Amount.GreaterThan(monthlyFee) {
		return ErrFeeExceedsMonthlyCap
	}
	if fee.EffectiveDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return ErrInvalidEffectiveDate
	}
	if fee.PaidAt != nil && fee.PaidAt.Before(fee.EffectiveDate) {
		return ErrInvalidPaidAt
	}
	return nil
}
```

**Error handling pattern** (lines 86-88, 112-116):
```go
if err != nil {
    return nil, fmt.Errorf("failed to create tenant: %w", err)
}
```

---

### `apps/api/internal/infrastructure/repository/tenant_repository.go` — repository-pgx, CRUD

**Analog:** `apps/api/internal/infrastructure/repository/user_repository.go` (lines 1-163)

**Imports pattern** (lines 1-11):
```go
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"harmoni-api/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)
```

**Struct + Constructor pattern** (lines 13-21):
```go
// PostgresTenantRepository implements TenantRepository using PostgreSQL via pgx.
type PostgresTenantRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresTenantRepository creates a new tenant repository backed by PostgreSQL.
func NewPostgresTenantRepository(pool *pgxpool.Pool) *PostgresTenantRepository {
	return &PostgresTenantRepository{pool: pool}
}
```

**Core CRUD pattern: Create** (from user_repository.go lines 24-46):
```go
// Create inserts a new tenant. The ID is generated by PostgreSQL's uuidv7().
func (r *PostgresTenantRepository) Create(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error) {
	query := `
		INSERT INTO tenants (block, unit_number, occupancy, monthly_fee, territory_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		tenant.Block,
		tenant.UnitNumber,
		tenant.Occupancy,
		tenant.MonthlyFee,
		tenant.TerritoryID,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}
```

**Core CRUD pattern: Territory-filtered list** (from user_repository.go lines 123-163):
```go
// ListByTerritory returns all tenants in a given territory.
func (r *PostgresTenantRepository) ListByTerritory(ctx context.Context, territoryID string) ([]*entity.Tenant, error) {
	query := `
		SELECT id, block, unit_number, occupancy, monthly_fee, territory_id, created_at, updated_at
		FROM tenants
		WHERE territory_id = $1
		ORDER BY block, unit_number
	`

	rows, err := r.pool.Query(ctx, query, territoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants by territory: %w", err)
	}
	defer rows.Close()

	var tenants []*entity.Tenant
	for rows.Next() {
		t := &entity.Tenant{}
		err := rows.Scan(
			&t.ID, &t.Block, &t.UnitNumber, &t.Occupancy,
			&t.MonthlyFee, &t.TerritoryID, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tenant row: %w", err)
		}
		tenants = append(tenants, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tenant rows: %w", err)
	}

	return tenants, nil
}
```

**Transactional variant pattern** (from password_reset_token_repository.go — accepts pgx.Tx):
```go
// CreateTx creates a tenant within an existing transaction.
func (r *PostgresTenantRepository) CreateTx(ctx context.Context, tx pgx.Tx, tenant *entity.Tenant) (*entity.Tenant, error) {
	query := `
		INSERT INTO tenants (block, unit_number, occupancy, monthly_fee, territory_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(ctx, query,
		tenant.Block, tenant.UnitNumber, tenant.Occupancy,
		tenant.MonthlyFee, tenant.TerritoryID,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create tenant in transaction: %w", err)
	}
	return tenant, nil
}
```

---

### `apps/api/internal/infrastructure/repository/fee_repository.go` — repository-pgx, CRUD

**Analog:** `apps/api/internal/infrastructure/repository/password_reset_token_repository.go` (lines 1-96)

Uses identical patterns — struct + constructor, raw SQL with `$N` placeholders, `rows.Scan()` for lists, `RowsAffected()` check for updates/deletes.

**Key difference:** Need separate methods for mandatory vs voluntary fees. Use `CreateMandatoryTx` variant for transactional fee creation within tenant creation.

---

### `apps/api/internal/delivery/http/tenant_handler.go` — handler, request-response

**Analog (role):** `apps/api/internal/delivery/http/auth_handler.go` (lines 1-264)
**Analog (data flow):** `apps/api/internal/delivery/http/protected_handler.go` (lines 1-267)

**Imports pattern** (auth_handler.go lines 1-11):
```go
package http

import (
	"harmoni-api/internal/domain/service"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
)
```

**Struct + Constructor pattern** (auth_handler.go lines 13-22):
```go
// TenantHandler handles tenant and fee HTTP endpoints.
type TenantHandler struct {
	tenantService *service.TenantService
	pasetoService *auth.PasetoService
}

// NewTenantHandler creates a new tenant handler.
func NewTenantHandler(tenantService *service.TenantService, pasetoService *auth.PasetoService) *TenantHandler {
	return &TenantHandler{tenantService: tenantService, pasetoService: pasetoService}
}
```

**Route registration pattern** (auth_handler.go lines 25-32, modified for nested + plural per CONTEXT.md):
```go
// RegisterRoutes registers tenant and fee endpoints on the Fiber API group.
func (h *TenantHandler) RegisterRoutes(api fiber.Router, pasetoSvc *auth.PasetoService) {
	// Tenant CRUD
	api.Get("/tenants", h.ListTenants)
	api.Post("/tenants", h.CreateTenant)
	api.Get("/tenants/:id", h.GetTenant)
	api.Put("/tenants/:id", h.UpdateTenant)
	api.Delete("/tenants/:id", h.DeleteTenant)

	// Fee sub-resources (nested under tenant)
	api.Get("/tenants/:id/fees", h.ListFees)
	api.Post("/tenants/:id/fees", h.CreateFee)
	api.Put("/fees/:feeId", h.UpdateFee)
	api.Delete("/fees/:feeId", h.DeleteFee)
}
```

**Handler method pattern** (auth_handler.go `Register` lines 43-77):
```go
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
```

**Create handler with request struct** (auth_handler.go lines 35-40, 43-77):
```go
// CreateTenantRequest represents the tenant creation request body.
type CreateTenantRequest struct {
	Block         string                        `json:"block"`
	UnitNumber    string                        `json:"unit_number"`
	Occupancy     string                        `json:"occupancy"`
	MonthlyFee    float64                       `json:"monthly_fee"`
	TerritoryID   string                        `json:"territory_id,omitempty"` // RW only
	MandatoryFees []CreateMandatoryFeeRequest   `json:"mandatory_fees"`
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
	tenant, err := h.tenantService.Create(c.Context(), &req, claims)
	if err != nil {
		switch {
		case err == service.ErrDuplicateBlockUnit:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Block and unit number already exist in this territory",
				"code":  "DUPLICATE_BLOCK_UNIT",
			})
		case err == service.ErrMandatoryFeeRequired:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one mandatory fee is required",
				"code":  "MANDATORY_FEE_REQUIRED",
			})
		case err == service.ErrFeeExceedsMonthlyCap:
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
```

**getUserClaims helper** (protected_handler.go lines 244-255 — already exists, reused):
```go
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
```

---

### `apps/api/internal/delivery/middleware/audit.go` — middleware, request-response

**Analog:** `apps/api/internal/delivery/middleware/casbin.go` (lines 1-112)

**Imports pattern**:
```go
package middleware

import (
	"log"

	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
)
```

**Middleware pattern** (casbin.go lines 9-16, 20-74):
```go
// AuditMiddlewareConfig holds configuration for the cross-territory audit logging middleware.
type AuditMiddlewareConfig struct {
	// Enforcer is the Casbin enforcer for detecting denials.
	Enforcer *auth.CasbinEnforcer
}

// NewAuditMiddleware creates a Fiber middleware that logs cross-territory access attempts
// and returns 403. It should be placed after the Casbin middleware in the chain.
func NewAuditMiddleware(cfg AuditMiddlewareConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user claims
		user := c.Locals("user")
		if user == nil {
			return c.Next()
		}

		claims, ok := user.(*auth.Claims)
		if !ok {
			return c.Next()
		}

		// Proceed with the request
		err := c.Next()

		// After handler: check if access was denied (403)
		if err != nil || c.Response().StatusCode() == fiber.StatusForbidden {
			log.Printf("AUDIT: cross-territory access attempt by user=%s role=%s territory=%s path=%s",
				claims.UserID, claims.Role, claims.TerritoryID, c.Path())
		}

		return err
	}
}
```

**Alternative: Function-based middleware** (simpler, logs before Casbin deny):
```go
// AuditCrossTerritoryAccess logs any cross-territory access attempt.
// This middleware runs before the Casbin middleware and logs the request context.
func AuditCrossTerritoryAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request details for audit trail
		user := c.Locals("user")
		if user != nil {
			if claims, ok := user.(*auth.Claims); ok {
				log.Printf("AUDIT: access by user=%s role=%s territory=%s method=%s path=%s",
					claims.UserID, claims.Role, claims.TerritoryID, c.Method(), c.Path())
			}
		}
		return c.Next()
	}
}
```

---

### Migrations (`006` through `010`)

**Analog (DDL):** `apps/api/migrations/002_create_users_table.up.sql` (lines 1-15)
**Analog (seed):** `apps/api/migrations/005_seed_territories.up.sql` (lines 1-4)

**Migration naming convention:**
```
NNN_description.up.sql
NNN_description.down.sql
```

**DDL pattern** (from 002_create_users_table.up.sql):
```sql
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    block VARCHAR(10) NOT NULL,
    unit_number VARCHAR(10) NOT NULL,
    occupancy VARCHAR(20) NOT NULL CHECK (occupancy IN ('occupied', 'vacant')),
    monthly_fee NUMERIC(12,2) NOT NULL,
    territory_id VARCHAR(50) NOT NULL REFERENCES territories(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(block, unit_number, territory_id)
);

CREATE INDEX idx_tenants_territory ON tenants(territory_id);
```

**Seed pattern** (from 005_seed_territories.up.sql):
```sql
INSERT INTO mandatory_fees (tenant_id, amount, description, effective_date) VALUES
    ('tenant-uuid-1', 50000.00, 'Security Fee', '2026-06-01');
```

**Migration plan:**
- `006_create_tenants_table.up.sql` — CREATE TABLE tenants with unique(block, unit_number, territory_id)
- `007_create_mandatory_fees_table.up.sql` — CREATE TABLE mandatory_fees (tenant_id FK, amount NUMERIC, description, effective_date, paid_at)
- `008_create_voluntary_fees_table.up.sql` — CREATE TABLE voluntary_fees (tenant_id FK, amount NUMERIC, description, effective_date, paid_at)
- `009_create_user_tenants_table.up.sql` — CREATE TABLE user_tenants (user_id FK, tenant_id FK, PK composite)
- `010_seed_sample_tenants.up.sql` — Optional seed data

Each must have a corresponding `.down.sql` with `DROP TABLE IF EXISTS ... CASCADE`.

---

### `apps/api/internal/delivery/http/protected_handler.go` — **MODIFY** — Remove tenant stubs

**What to change:** Remove lines 41-44 (tenant route registration) and lines 113-152 (ListTenants, GetTenant, CreateTenant stubs). These routes move to `tenant_handler.go` with real implementations.

---

### `apps/api/cmd/server/main.go` — **MODIFY** — Register TenantHandler

**Analog:** Same file (lines 47-64, 105-115)

**Add after line 113** (before protectedHandler is created):
```go
// Initialize tenant repository and service
tenantRepo := infrarepo.NewPostgresTenantRepository(db.Pool)
feeRepo := infrarepo.NewPostgresFeeRepository(db.Pool)

tenantService := service.NewTenantService(tenantRepo, feeRepo, db.Pool)

// Register tenant routes (under the same /api group as protected routes)
tenantHandler := httphandler.NewTenantHandler(tenantService, pasetoService)
tenantHandler.RegisterRoutes(api, pasetoService)
```

**Note:** Must also add import aliases:
- `"harmoni-api/internal/domain/service"` — already imported
- `infrarepo "harmoni-api/internal/infrastructure/repository"` — already imported
- `httphandler "harmoni-api/internal/delivery/http"` — already imported

**Caveat:** The `/api` group is created inside `ProtectedHandler.RegisterRoutes()`. To add tenant routes to the same group, refactor `main.go` to create the `/api` group at the `main()` level and pass it to both handlers, or have `TenantHandler.RegisterRoutes` accept the same group.

---

### `apps/api/policy.csv` — **MODIFY** — Add fee resources

**Analog:** Same file (lines 1-53)

**Add after existing policy** (before role hierarchy):
```csv
# ── Fee Resources ──
# Residents can read fees for their own tenants
p, resident, mandatory_fee, read, *
p, resident, voluntary_fee, read, *

# RT officers can manage fees within their territory
p, rt_officer, mandatory_fee, read, *
p, rt_officer, mandatory_fee, write, *
p, rt_officer, voluntary_fee, read, *
p, rt_officer, voluntary_fee, write, *

# RW officers can manage fees across all territories
p, rw_officer, mandatory_fee, read, *
p, rw_officer, mandatory_fee, write, *
p, rw_officer, voluntary_fee, read, *
p, rw_officer, voluntary_fee, write, *
```

---

### `apps/api/rbac_model.conf` — **MODIFY** — Extend matcher for resident tenant access

**Analog:** Same file (lines 1-35)

**Add custom function for resident tenant scoping** (between `[policy_effect]` and `[matchers]`):
```
# For resident access scoped by user_tenants junction table:
# The matcher calls a Go-provided function to check if the requested
# tenant_id is in the user's list of assigned tenants.
```

**Implementation note per RESEARCH.md (Open Question #3):** Rather than modifying the Casbin model file directly, extend the `CasbinMiddleware` to accept a `UserTenantIDs` function. The `Enforce()` call passes the user's tenant ID list as an additional parameter. The matcher checks `r.tenant_id_in` against the parameter list. This approach keeps logic in the Go layer and avoids Casbin model complexity.

---

### Test Files

#### `apps/api/internal/domain/service/tenant_service_test.go` — test, unit

**Analog:** `apps/api/internal/domain/service/auth_service_test.go` (lines 1-426)

**Mock repository pattern** (lines 17-31):
```go
// mockTenantRepository implements repository.TenantRepository for testing.
type mockTenantRepository struct {
	tenants    map[string]*entity.Tenant
	createErr  error
	findErr    error
}

func newMockTenantRepo() *mockTenantRepository {
	return &mockTenantRepository{
		tenants: make(map[string]*entity.Tenant),
	}
}
```

**Test helper pattern** (lines 157-171):
```go
func newTestTenantService(t *testing.T) (*TenantService, *mockTenantRepository, *mockFeeRepository) {
	t.Helper()

	tenantRepo := newMockTenantRepo()
	feeRepo := newMockFeeRepo()

	svc := NewTenantService(tenantRepo, feeRepo, nil) // nil pool for unit tests
	return svc, tenantRepo, feeRepo
}
```

**Test case pattern** (lines 173-196):
```go
func TestTenantService_Create_Success(t *testing.T) {
	svc, _, _ := newTestTenantService(t)

	// ... setup request and claims ...

	result, err := svc.Create(ctx, req, claims)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result.Block != "A" {
		t.Errorf("Block = %q, want %q", result.Block, "A")
	}
}
```

**Error case pattern** (lines 198-210):
```go
func TestTenantService_Create_MandatoryFeeRequired(t *testing.T) {
	svc, _, _ := newTestTenantService(t)

	// ... setup request with no mandatory fees ...

	_, err := svc.Create(ctx, req, claims)
	if !errors.Is(err, ErrMandatoryFeeRequired) {
		t.Errorf("Create() error = %v, want ErrMandatoryFeeRequired", err)
	}
}
```

---

#### `apps/api/internal/delivery/http/tenant_handler_test.go` — test, integration

**Analog (service tests):** `apps/api/internal/delivery/http/auth_handler_test.go` (lines 1-467)
**Analog (protected route tests):** `apps/api/internal/delivery/http/protected_handler_test.go` (lines 1-242)

**Test setup pattern** (auth_handler_test.go lines 124-143):
```go
func setupTenantHandlerTest(t *testing.T) *fiber.App {
	t.Helper()

	// ... setup mocks and service ...

	app := fiber.New()
	handler := NewTenantHandler(tenantSvc, paseto)
	handler.RegisterRoutes(app, paseto)

	return app
}
```

**HTTP test pattern** (auth_handler_test.go lines 145-178):
```go
func TestTenantHandler_Create_Success(t *testing.T) {
	app := setupTenantHandlerTest(t)

	body := map[string]interface{}{
		"block":       "A",
		"unit_number": "01",
		"occupancy":   "occupied",
		"monthly_fee": 50000,
		"mandatory_fees": []map[string]interface{}{
			{"amount": 25000, "description": "Security Fee", "effective_date": "2026-06-01"},
		},
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/tenants", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "paseto_token", Value: validToken})

	resp, _ := app.Test(req)
	assert.Equal(t, 201, resp.StatusCode)
}
```

**Protected route test pattern** (protected_handler_test.go lines 17-30, 63-82):
```go
func setupProtectedApp(t *testing.T, pasetoSvc *auth.PasetoService, enforcer *auth.CasbinEnforcer) *fiber.App {
	t.Helper()

	app := fiber.New()
	handler := NewProtectedHandler(enforcer)
	handler.RegisterRoutes(app, pasetoSvc)

	return app
}
```

---

## Shared Patterns

### Work in progress — see tasks list below

1. ⬜ **Configure your editor** (VS Code, Zed, etc.)
2. ⬜ **Learn the keyboard shortcuts** (in editor)
3. ⬜ **Review open files** from the Explorer (`⌘E`) and file tree
4. ⬜ Open the **Activity Bar** and explore extensions, search, source control, etc.
5. ⬜ Open a terminal (`⌃``) and run a few commands
6. ⬜ Run `opencode doctor` if you run into environment issues
7. ⬜ Review the [documentation](https://opencode.ai/docs) if you have questions
8. ⬜ Try a simple voice command in the prompt bar

### Authentication
**Source:** `apps/api/internal/delivery/middleware/auth.go` (lines 1-73)
**Apply to:** All protected endpoints (already applied via middleware chain in `protected_handler.go`)

**Core token validation pattern** (lines 19-63):
```go
func NewAuthMiddleware(cfg AuthMiddlewareConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
			code := "INVALID_TOKEN"
			if strings.Contains(err.Error(), "expired") {
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
```

**Claims struct** (paseto.go lines 13-18):
```go
type Claims struct {
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"`
	TerritoryID string    `json:"territory_id"`
	Expiration  time.Time `json:"exp"`
}
```

### Casbin Authorization
**Source:** `apps/api/internal/delivery/middleware/casbin.go` (lines 1-112)
**Apply to:** All `/api/*` routes (already applied)

**Core enforcement pattern** (lines 25-74):
```go
return func(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized", "code": "NO_USER_CONTEXT",
		})
	}
	claims := user.(*auth.Claims)

	action := methodToAction(c.Method())
	resource := cfg.ResourceExtractor(c)

	domain := claims.TerritoryID
	if claims.Role == "rw_officer" {
		domain = "*"
	}

	allowed, err := cfg.Enforcer.Enforce(claims.Role, resource, action, domain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error", "code": "ENFORCE_ERROR",
		})
	}

	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden", "code": "INSUFFICIENT_PERMISSIONS",
		})
	}

	return c.Next()
}
```

**Note for Phase 2:** The `defaultResourceExtractor` extracts the first segment after `/api/`. For `GET /api/tenants/:id/fees`, it returns `"tenants"`. The Casbin policy resource must use the segment returned by this extractor. If finer granularity for fees is needed, provide a custom `ResourceExtractor`.

### Error Response Format
**Source:** All Phase 1 handlers
**Apply to:** All handlers
**Pattern:**
```json
{"error": "Human-readable message", "code": "MACHINE_CODE"}
```

### Service-Layer Error Definitions
**Source:** `apps/api/internal/domain/service/auth_service.go` (lines 17-24)
**Apply to:** `tenant_service.go`
**Pattern:**
```go
var (
	ErrTenantNotFound       = errors.New("tenant not found")
	ErrDuplicateBlockUnit   = errors.New("block and unit number already exist in territory")
	ErrMandatoryFeeRequired = errors.New("at least one mandatory fee is required")
)
```

### Repository-Layer Error Wrapping
**Source:** `apps/api/internal/infrastructure/repository/user_repository.go`
**Apply to:** All infrastructure repositories
**Pattern:**
```go
// Returns sql.ErrNoRows when not found — consistent with Phase 1 convention
if err != nil {
    return nil, fmt.Errorf("failed to create tenant: %w", err)
}
```

### Package Naming Conventions
**Source:** All Phase 1 files
**Apply to:** All Phase 2 new files

| Layer | Package Name | Directory |
|-------|-------------|-----------|
| Entity | `entity` | `internal/domain/entity/` |
| Repository interface | `repository` | `internal/domain/repository/` |
| Service | `service` | `internal/domain/service/` |
| pgx repository | `repository` | `internal/infrastructure/repository/` |
| HTTP handler | `http` | `internal/delivery/http/` |
| Middleware | `middleware` | `internal/delivery/middleware/` |

---

## No Analog Found

All files have close matches from Phase 1 codebase. No files require patterns from RESEARCH.md alone.

| File | Role | Data Flow | Reason |
|------|------|-----------|--------|
| — | — | — | All have analogs |

---

## Metadata

**Analog search scope:**
- `apps/api/internal/domain/entity/` — 1 file scanned
- `apps/api/internal/domain/repository/` — 2 files scanned
- `apps/api/internal/domain/service/` — 1 file scanned
- `apps/api/internal/infrastructure/repository/` — 2 files scanned
- `apps/api/internal/delivery/http/` — 2 files scanned
- `apps/api/internal/delivery/middleware/` — 2 files scanned
- `apps/api/cmd/server/main.go` — 1 file scanned
- `apps/api/internal/infrastructure/auth/` — 2 files scanned
- `apps/api/internal/config/` — 1 file scanned
- `apps/api/migrations/` — 5 files scanned
- `apps/api/policy.csv` — 1 file scanned
- `apps/api/rbac_model.conf` — 1 file scanned
- `apps/api/internal/infrastructure/database/` — 1 file scanned

**Files scanned:** 22 source files
**Pattern extraction date:** 2026-05-22
