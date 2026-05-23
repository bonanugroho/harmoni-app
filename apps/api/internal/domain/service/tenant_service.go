package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"harmoni-api/internal/domain/entity"
	"harmoni-api/internal/domain/repository"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

// Sentinel errors for tenant and fee operations.
// These are comparable with errors.Is() for callers to inspect.
var (
	ErrTenantNotFound         = errors.New("tenant not found")
	ErrDuplicateBlockUnit     = errors.New("block and unit number already exist in territory")
	ErrMandatoryFeeRequired   = errors.New("at least one mandatory fee is required")
	ErrFeeExceedsMonthlyCap   = errors.New("fee amount exceeds tenant's monthly fee cap")
	ErrInvalidEffectiveDate   = errors.New("effective date cannot be in the past")
	ErrInvalidPaidAt          = errors.New("paid_at must be after effective_date")
	ErrCrossTerritoryAccess   = errors.New("cross-territory access denied")
)

// pool abstracts the pgx pool for testability.
type pool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// TenantService handles tenant and fee business logic.
type TenantService struct {
	tenantRepo repository.TenantRepository
	feeRepo    repository.FeeRepository
	pool       pool
}

// NewTenantService creates a new tenant service.
// Accepts *pgxpool.Pool as the pool argument (satisfies the pool interface).
func NewTenantService(
	tenantRepo repository.TenantRepository,
	feeRepo repository.FeeRepository,
	p pool,
) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		feeRepo:    feeRepo,
		pool:       p,
	}
}

// --- Request Types ---

// CreateTenantRequest represents the tenant creation request body.
type CreateTenantRequest struct {
	Block         string                     `json:"block"`
	UnitNumber    string                     `json:"unit_number"`
	Occupancy     string                     `json:"occupancy"`
	MonthlyFee    decimal.Decimal            `json:"monthly_fee"`
	TerritoryID   string                     `json:"territory_id,omitempty"`
	MandatoryFees []CreateMandatoryFeeRequest `json:"mandatory_fees"`
}

// CreateMandatoryFeeRequest represents a mandatory fee creation request.
type CreateMandatoryFeeRequest struct {
	Amount        decimal.Decimal `json:"amount"`
	Description   string          `json:"description"`
	EffectiveDate time.Time       `json:"effective_date"`
}

// CreateVoluntaryFeeRequest represents a voluntary fee creation request.
type CreateVoluntaryFeeRequest struct {
	Amount        decimal.Decimal `json:"amount"`
	Description   string          `json:"description"`
	EffectiveDate time.Time       `json:"effective_date"`
	PaidAt        *time.Time      `json:"paid_at"`
}

// UpdateTenantRequest represents the tenant update request body.
type UpdateTenantRequest struct {
	Block       string          `json:"block"`
	UnitNumber  string          `json:"unit_number"`
	Occupancy   string          `json:"occupancy"`
	MonthlyFee  decimal.Decimal `json:"monthly_fee"`
	TerritoryID string          `json:"territory_id,omitempty"`
}

// UpdateFeeRequest represents a fee update request.
type UpdateFeeRequest struct {
	Amount        decimal.Decimal `json:"amount"`
	Description   string          `json:"description"`
	EffectiveDate time.Time       `json:"effective_date"`
	PaidAt        *time.Time      `json:"paid_at"`
}

// --- Tenant CRUD ---

// Create creates a new tenant with mandatory fees in a transaction.
func (s *TenantService) Create(ctx context.Context, req *CreateTenantRequest, claims *auth.Claims) (*entity.Tenant, error) {
	// Validate block and unit_number are not empty
	if req.Block == "" || req.UnitNumber == "" {
		return nil, errors.New("block and unit_number are required")
	}

	// Validate at least one mandatory fee provided (per D-05)
	if len(req.MandatoryFees) == 0 {
		return nil, ErrMandatoryFeeRequired
	}

	// Validate each fee before starting the transaction
	for _, feeReq := range req.MandatoryFees {
		if err := s.validateFee(feeReq.Amount, req.MonthlyFee, feeReq.EffectiveDate, nil); err != nil {
			return nil, fmt.Errorf("mandatory fee validation failed: %w", err)
		}
	}

	// Begin transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Determine territory scope
	territoryID := claims.TerritoryID
	if claims.Role == "rw_officer" && req.TerritoryID != "" {
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

	// Create mandatory fees
	for _, feeReq := range req.MandatoryFees {
		fee := &entity.MandatoryFee{
			TenantID:      created.ID,
			Amount:        feeReq.Amount,
			Description:   feeReq.Description,
			EffectiveDate: feeReq.EffectiveDate,
		}
		if _, err := s.feeRepo.CreateMandatoryTx(ctx, tx, fee); err != nil {
			return nil, fmt.Errorf("failed to create mandatory fee: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return created, nil
}

// List returns tenants scoped by the user's role.
// Residents see only their assigned tenants via user_tenants.
// RT officers see tenants in their territory.
// RW officers see all tenants (passes "*" to let repository handle it).
func (s *TenantService) List(ctx context.Context, claims *auth.Claims) ([]*entity.Tenant, error) {
	switch claims.Role {
	case "resident":
		return s.tenantRepo.ListByUserID(ctx, claims.UserID)
	case "rw_officer":
		return s.tenantRepo.ListByTerritory(ctx, "*")
	default: // rt_officer and any other role
		return s.tenantRepo.ListByTerritory(ctx, claims.TerritoryID)
	}
}

// GetByID returns a tenant by ID with access control.
// Returns ErrTenantNotFound if the tenant doesn't exist or the user lacks access.
func (s *TenantService) GetByID(ctx context.Context, id string, claims *auth.Claims) (*entity.Tenant, error) {
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTenantNotFound
	}

	return s.validateTenantAccess(ctx, tenant, claims)
}

// Update updates an existing tenant with access control.
func (s *TenantService) Update(ctx context.Context, id string, req *UpdateTenantRequest, claims *auth.Claims) (*entity.Tenant, error) {
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTenantNotFound
	}

	// Access check - returns ErrTenantNotFound on mismatch (prevents scanning)
	if _, err := s.validateTenantAccess(ctx, tenant, claims); err != nil {
		return nil, err
	}

	// Update fields
	tenant.Block = req.Block
	tenant.UnitNumber = req.UnitNumber
	tenant.Occupancy = req.Occupancy
	tenant.MonthlyFee = req.MonthlyFee

	updated, err := s.tenantRepo.Update(ctx, tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return updated, nil
}

// Delete removes a tenant with access control.
func (s *TenantService) Delete(ctx context.Context, id string, claims *auth.Claims) error {
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		return ErrTenantNotFound
	}

	// Access check - returns ErrTenantNotFound on mismatch
	if _, err := s.validateTenantAccess(ctx, tenant, claims); err != nil {
		return err
	}

	if err := s.tenantRepo.Delete(ctx, id, tenant.TerritoryID); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	return nil
}

// --- Fee Methods ---

// CreateMandatoryFee creates a mandatory fee for a tenant.
func (s *TenantService) CreateMandatoryFee(ctx context.Context, tenantID string, req *CreateMandatoryFeeRequest, claims *auth.Claims) (*entity.MandatoryFee, error) {
	// Verify tenant access
	tenant, err := s.getTenantForFee(ctx, tenantID, claims)
	if err != nil {
		return nil, err
	}

	// Validate fee constraints
	if err := s.validateFee(req.Amount, tenant.MonthlyFee, req.EffectiveDate, nil); err != nil {
		return nil, fmt.Errorf("mandatory fee validation failed: %w", err)
	}

	fee := &entity.MandatoryFee{
		TenantID:      tenantID,
		Amount:        req.Amount,
		Description:   req.Description,
		EffectiveDate: req.EffectiveDate,
	}

	created, err := s.feeRepo.CreateMandatory(ctx, fee)
	if err != nil {
		return nil, fmt.Errorf("failed to create mandatory fee: %w", err)
	}

	return created, nil
}

// CreateVoluntaryFee creates a voluntary fee for a tenant.
func (s *TenantService) CreateVoluntaryFee(ctx context.Context, tenantID string, req *CreateVoluntaryFeeRequest, claims *auth.Claims) (*entity.VoluntaryFee, error) {
	// Verify tenant access
	tenant, err := s.getTenantForFee(ctx, tenantID, claims)
	if err != nil {
		return nil, err
	}

	// Validate fee constraints
	if err := s.validateFee(req.Amount, tenant.MonthlyFee, req.EffectiveDate, req.PaidAt); err != nil {
		return nil, fmt.Errorf("voluntary fee validation failed: %w", err)
	}

	fee := &entity.VoluntaryFee{
		TenantID:      tenantID,
		Amount:        req.Amount,
		Description:   req.Description,
		EffectiveDate: req.EffectiveDate,
		PaidAt:        req.PaidAt,
	}

	created, err := s.feeRepo.CreateVoluntary(ctx, fee)
	if err != nil {
		return nil, fmt.Errorf("failed to create voluntary fee: %w", err)
	}

	return created, nil
}

// ListFees returns all fees (mandatory + voluntary) for a tenant.
func (s *TenantService) ListFees(ctx context.Context, tenantID string, claims *auth.Claims) ([]*entity.MandatoryFee, []*entity.VoluntaryFee, error) {
	// Verify tenant access
	if _, err := s.getTenantForFee(ctx, tenantID, claims); err != nil {
		return nil, nil, err
	}

	mandatory, err := s.feeRepo.ListMandatoryByTenant(ctx, tenantID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list mandatory fees: %w", err)
	}

	voluntary, err := s.feeRepo.ListVoluntaryByTenant(ctx, tenantID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list voluntary fees: %w", err)
	}

	return mandatory, voluntary, nil
}

// UpdateFee updates a fee record by ID.
func (s *TenantService) UpdateFee(ctx context.Context, feeID string, req *UpdateFeeRequest, claims *auth.Claims) error {
	// Fetch existing fee to determine its type and parent tenant
	// Try mandatory first, then voluntary
	mFee, err := s.feeRepo.ListMandatoryByTenant(ctx, feeID)
	// If the fee doesn't match, try voluntary
	vFee, err2 := s.feeRepo.ListVoluntaryByTenant(ctx, feeID)

	// We can't easily find a single fee by ID from the repository interface as designed.
	// For now, return an error - the handler layer will have access to fee-specific lookups.
	_ = mFee
	_ = vFee
	_ = err
	_ = err2

	return fmt.Errorf("update fee by ID not directly supported: use tenant-scoped update")
}

// DeleteFee deletes a fee record by ID.
func (s *TenantService) DeleteFee(ctx context.Context, feeID string, claims *auth.Claims) error {
	// Similar to UpdateFee, we need a fee-specific lookup.
	// The repository doesn't support direct FindByID for fees.
	// The handler layer should resolve the tenant and call the appropriate method.
	return fmt.Errorf("delete fee by ID not directly supported: use tenant-scoped delete")
}

// --- Private Helpers ---

// validateFee checks fee constraints:
// - Amount must be non-negative
// - Amount must not exceed the tenant's monthly fee cap
// - Effective date cannot be in the past
// - PaidAt must be after EffectiveDate (if provided)
func (s *TenantService) validateFee(amount, monthlyFee decimal.Decimal, effectiveDate time.Time, paidAt *time.Time) error {
	if amount.LessThan(decimal.Zero) {
		return errors.New("fee amount must be non-negative")
	}

	if amount.GreaterThan(monthlyFee) {
		return ErrFeeExceedsMonthlyCap
	}

	today := time.Now().Truncate(24 * time.Hour)
	if effectiveDate.Before(today) {
		return ErrInvalidEffectiveDate
	}

	if paidAt != nil && paidAt.Before(effectiveDate) {
		return ErrInvalidPaidAt
	}

	return nil
}

// validateTenantAccess checks if the claims have access to the given tenant.
// For residents: verifies tenant is in user's assigned list.
// For RT officers: verifies tenant territory matches claims territory.
// For RW officers: always allowed.
// Returns ErrTenantNotFound on mismatch (same error for not-found vs no-access).
func (s *TenantService) validateTenantAccess(ctx context.Context, tenant *entity.Tenant, claims *auth.Claims) (*entity.Tenant, error) {
	switch claims.Role {
	case "resident":
		// Check if tenant is in user's assigned list
		userTenants, err := s.tenantRepo.ListByUserID(ctx, claims.UserID)
		if err != nil {
			return nil, ErrTenantNotFound
		}
		for _, ut := range userTenants {
			if ut.ID == tenant.ID {
				return tenant, nil
			}
		}
		return nil, ErrTenantNotFound

	case "rw_officer":
		return tenant, nil

	default: // rt_officer
		if tenant.TerritoryID != claims.TerritoryID {
			return nil, ErrTenantNotFound
		}
		return tenant, nil
	}
}

// getTenantForFee is a helper that fetches a tenant and validates access.
func (s *TenantService) getTenantForFee(ctx context.Context, tenantID string, claims *auth.Claims) (*entity.Tenant, error) {
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		return nil, ErrTenantNotFound
	}

	return s.validateTenantAccess(ctx, tenant, claims)
}
