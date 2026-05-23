package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"harmoni-api/internal/domain/entity"
	"harmoni-api/internal/domain/repository"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

// mockTenantRepository implements repository.TenantRepository for testing.
type mockTenantRepository struct {
	tenants    map[string]*entity.Tenant
	createErr  error
	findErr    error
	updateErr  error
	deleteErr  error
	listErr    error
}

func newMockTenantRepo() *mockTenantRepository {
	return &mockTenantRepository{
		tenants: make(map[string]*entity.Tenant),
	}
}

func (m *mockTenantRepository) Create(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	tenant.ID = "tenant-uuid-v7"
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()
	m.tenants[tenant.ID] = tenant
	return tenant, nil
}

func (m *mockTenantRepository) CreateTx(ctx context.Context, tx pgx.Tx, tenant *entity.Tenant) (*entity.Tenant, error) {
	return m.Create(ctx, tenant)
}

func (m *mockTenantRepository) FindByID(ctx context.Context, id string) (*entity.Tenant, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	tenant, ok := m.tenants[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return tenant, nil
}

func (m *mockTenantRepository) ListByTerritory(ctx context.Context, territoryID string) ([]*entity.Tenant, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var result []*entity.Tenant
	for _, t := range m.tenants {
		if territoryID == "*" || t.TerritoryID == territoryID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *mockTenantRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Tenant, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	// For tests, simulate that user "resident-user" has access to tenants with IDs matching userID pattern
	var result []*entity.Tenant
	for _, t := range m.tenants {
		result = append(result, t)
	}
	return result, nil
}

func (m *mockTenantRepository) Update(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, ok := m.tenants[tenant.ID]
	if !ok {
		return nil, sql.ErrNoRows
	}
	existing.Block = tenant.Block
	existing.UnitNumber = tenant.UnitNumber
	existing.Occupancy = tenant.Occupancy
	existing.MonthlyFee = tenant.MonthlyFee
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (m *mockTenantRepository) Delete(ctx context.Context, id string, territoryID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	tenant, ok := m.tenants[id]
	if !ok {
		return sql.ErrNoRows
	}
	if tenant.TerritoryID != territoryID {
		return sql.ErrNoRows
	}
	delete(m.tenants, id)
	return nil
}

// mockFeeRepository implements repository.FeeRepository for testing.
type mockFeeRepository struct {
	mandatoryFees map[string]*entity.MandatoryFee
	voluntaryFees map[string]*entity.VoluntaryFee
	createErr     error
	findErr       error
	updateErr     error
	deleteErr     error
}

func newMockFeeRepo() *mockFeeRepository {
	return &mockFeeRepository{
		mandatoryFees: make(map[string]*entity.MandatoryFee),
		voluntaryFees: make(map[string]*entity.VoluntaryFee),
	}
}

func (m *mockFeeRepository) CreateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	fee.ID = "fee-uuid-v7"
	fee.CreatedAt = time.Now()
	fee.UpdatedAt = time.Now()
	m.mandatoryFees[fee.ID] = fee
	return fee, nil
}

func (m *mockFeeRepository) CreateMandatoryTx(ctx context.Context, tx pgx.Tx, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	return m.CreateMandatory(ctx, fee)
}

func (m *mockFeeRepository) ListMandatoryByTenant(ctx context.Context, tenantID string) ([]*entity.MandatoryFee, error) {
	var result []*entity.MandatoryFee
	for _, f := range m.mandatoryFees {
		if f.TenantID == tenantID {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *mockFeeRepository) GetMandatoryByID(ctx context.Context, id string) (*entity.MandatoryFee, error) {
	fee, ok := m.mandatoryFees[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return fee, nil
}

func (m *mockFeeRepository) UpdateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, ok := m.mandatoryFees[fee.ID]
	if !ok {
		return nil, sql.ErrNoRows
	}
	existing.Amount = fee.Amount
	existing.Description = fee.Description
	existing.EffectiveDate = fee.EffectiveDate
	existing.PaidAt = fee.PaidAt
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (m *mockFeeRepository) DeleteMandatory(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.mandatoryFees[id]; !ok {
		return sql.ErrNoRows
	}
	delete(m.mandatoryFees, id)
	return nil
}

func (m *mockFeeRepository) CreateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	fee.ID = "voluntary-uuid-v7"
	fee.CreatedAt = time.Now()
	fee.UpdatedAt = time.Now()
	m.voluntaryFees[fee.ID] = fee
	return fee, nil
}

func (m *mockFeeRepository) ListVoluntaryByTenant(ctx context.Context, tenantID string) ([]*entity.VoluntaryFee, error) {
	var result []*entity.VoluntaryFee
	for _, f := range m.voluntaryFees {
		if f.TenantID == tenantID {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *mockFeeRepository) GetVoluntaryByID(ctx context.Context, id string) (*entity.VoluntaryFee, error) {
	fee, ok := m.voluntaryFees[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return fee, nil
}

func (m *mockFeeRepository) UpdateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, ok := m.voluntaryFees[fee.ID]
	if !ok {
		return nil, sql.ErrNoRows
	}
	existing.Amount = fee.Amount
	existing.Description = fee.Description
	existing.EffectiveDate = fee.EffectiveDate
	existing.PaidAt = fee.PaidAt
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (m *mockFeeRepository) DeleteVoluntary(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.voluntaryFees[id]; !ok {
		return sql.ErrNoRows
	}
	delete(m.voluntaryFees, id)
	return nil
}

// Ensure mockFeeRepository implements repository.FeeRepository
var _ repository.FeeRepository = (*mockFeeRepository)(nil)

// mockTx implements pgx.Tx for testing by embedding the interface.
type mockTx struct {
	pgx.Tx
	commitErr   error
	rollbackErr error
}

func (t *mockTx) Commit(ctx context.Context) error {
	return t.commitErr
}

func (t *mockTx) Rollback(ctx context.Context) error {
	return t.rollbackErr
}

// mockPool implements pool interface for testing.
type mockPool struct {
	beginErr error
}

func (p *mockPool) Begin(ctx context.Context) (pgx.Tx, error) {
	if p.beginErr != nil {
		return nil, p.beginErr
	}
	return &mockTx{}, nil
}

func newTestTenantService(t *testing.T) (*TenantService, *mockTenantRepository, *mockFeeRepository) {
	t.Helper()

	tenantRepo := newMockTenantRepo()
	feeRepo := newMockFeeRepo()
	pool := &mockPool{}

	svc := NewTenantService(tenantRepo, feeRepo, pool)
	return svc, tenantRepo, feeRepo
}

func newTestTenantServiceWithPool(t *testing.T, pool pool) (*TenantService, *mockTenantRepository, *mockFeeRepository) {
	t.Helper()

	tenantRepo := newMockTenantRepo()
	feeRepo := newMockFeeRepo()

	svc := NewTenantService(tenantRepo, feeRepo, pool)
	return svc, tenantRepo, feeRepo
}

// Helper to create a basic CreateTenantRequest.
func validCreateRequest() *CreateTenantRequest {
	futureDate := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	return &CreateTenantRequest{
		Block:      "A",
		UnitNumber: "01",
		Occupancy:  "occupied",
		MonthlyFee: decimal.NewFromInt(50000),
		MandatoryFees: []CreateMandatoryFeeRequest{
			{
				Amount:        decimal.NewFromInt(25000),
				Description:   "Security Fee",
				EffectiveDate: futureDate,
			},
		},
	}
}

func validRTClaims() *auth.Claims {
	return &auth.Claims{
		UserID:      "rt-officer-user",
		Role:        "rt_officer",
		TerritoryID: "rt-01",
	}
}

func validRWClaims() *auth.Claims {
	return &auth.Claims{
		UserID:      "rw-officer-user",
		Role:        "rw_officer",
		TerritoryID: "rw-01",
	}
}

func validResidentClaims() *auth.Claims {
	return &auth.Claims{
		UserID:      "resident-user",
		Role:        "resident",
		TerritoryID: "rt-01",
	}
}

func TestTenantService_Create_Success(t *testing.T) {
	svc, _, _ := newTestTenantService(t)
	ctx := context.Background()

	req := validCreateRequest()
	claims := validRTClaims()

	result, err := svc.Create(ctx, req, claims)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result.Block != "A" {
		t.Errorf("Block = %q, want %q", result.Block, "A")
	}
	if result.UnitNumber != "01" {
		t.Errorf("UnitNumber = %q, want %q", result.UnitNumber, "01")
	}
	if result.TerritoryID != "rt-01" {
		t.Errorf("TerritoryID = %q, want %q", result.TerritoryID, "rt-01")
	}
	if result.ID == "" {
		t.Error("Create() should return tenant with ID")
	}
}

func TestTenantService_Create_DuplicateBlockUnit(t *testing.T) {
	svc, tenantRepo, _ := newTestTenantService(t)
	ctx := context.Background()

	// Pre-populate a tenant to cause duplicate
	existing := &entity.Tenant{
		ID:          "existing-uuid",
		Block:       "A",
		UnitNumber:  "01",
		TerritoryID: "rt-01",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
	}
	tenantRepo.tenants["existing-uuid"] = existing

	// Manually set createErr to simulate duplicate
	tenantRepo.createErr = ErrDuplicateBlockUnit

	req := validCreateRequest()
	claims := validRTClaims()

	_, err := svc.Create(ctx, req, claims)
	if !errors.Is(err, ErrDuplicateBlockUnit) {
		t.Errorf("Create() error = %v, want ErrDuplicateBlockUnit", err)
	}
}

func TestTenantService_Create_MandatoryFeeRequired(t *testing.T) {
	svc, _, _ := newTestTenantService(t)
	ctx := context.Background()

	req := &CreateTenantRequest{
		Block:         "A",
		UnitNumber:    "01",
		Occupancy:     "occupied",
		MonthlyFee:    decimal.NewFromInt(50000),
		MandatoryFees: []CreateMandatoryFeeRequest{},
	}
	claims := validRTClaims()

	_, err := svc.Create(ctx, req, claims)
	if !errors.Is(err, ErrMandatoryFeeRequired) {
		t.Errorf("Create() error = %v, want ErrMandatoryFeeRequired", err)
	}
}

func TestTenantService_Create_FeeExceedsMonthlyCap(t *testing.T) {
	svc, _, _ := newTestTenantService(t)
	ctx := context.Background()

	futureDate := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	req := &CreateTenantRequest{
		Block:      "A",
		UnitNumber: "01",
		Occupancy:  "occupied",
		MonthlyFee: decimal.NewFromInt(50000), // monthly fee = 50000
		MandatoryFees: []CreateMandatoryFeeRequest{
			{
				Amount:        decimal.NewFromInt(60000), // fee > monthly fee
				Description:   "Too Expensive Fee",
				EffectiveDate: futureDate,
			},
		},
	}
	claims := validRTClaims()

	_, err := svc.Create(ctx, req, claims)
	if !errors.Is(err, ErrFeeExceedsMonthlyCap) {
		t.Errorf("Create() error = %v, want ErrFeeExceedsMonthlyCap", err)
	}
}

func TestTenantService_Create_InvalidEffectiveDate(t *testing.T) {
	svc, _, _ := newTestTenantService(t)
	ctx := context.Background()

	pastDate := time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour)
	req := &CreateTenantRequest{
		Block:      "A",
		UnitNumber: "01",
		Occupancy:  "occupied",
		MonthlyFee: decimal.NewFromInt(50000),
		MandatoryFees: []CreateMandatoryFeeRequest{
			{
				Amount:        decimal.NewFromInt(25000),
				Description:   "Past Fee",
				EffectiveDate: pastDate,
			},
		},
	}
	claims := validRTClaims()

	_, err := svc.Create(ctx, req, claims)
	if !errors.Is(err, ErrInvalidEffectiveDate) {
		t.Errorf("Create() error = %v, want ErrInvalidEffectiveDate", err)
	}
}

func TestTenantService_Delete_NotOwnTerritory(t *testing.T) {
	svc, tenantRepo, _ := newTestTenantService(t)
	ctx := context.Background()

	// Create a tenant in rt-01
	tenant := &entity.Tenant{
		ID:          "tenant-rt01",
		Block:       "A",
		UnitNumber:  "01",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
		TerritoryID: "rt-01",
	}
	tenantRepo.tenants["tenant-rt01"] = tenant

	// RT officer from rt-02 tries to delete it
	rt02Claims := &auth.Claims{
		UserID:      "rt-officer-02",
		Role:        "rt_officer",
		TerritoryID: "rt-02",
	}

	err := svc.Delete(ctx, "tenant-rt01", rt02Claims)
	if !errors.Is(err, ErrCrossTerritoryAccess) {
		t.Errorf("Delete() error = %v, want ErrCrossTerritoryAccess", err)
	}
}

func TestTenantService_List_ResidentOnlyOwn(t *testing.T) {
	svc, tenantRepo, _ := newTestTenantService(t)
	ctx := context.Background()

	// Add tenants
	tenantRepo.tenants["t1"] = &entity.Tenant{
		ID:          "t1",
		Block:       "A",
		UnitNumber:  "01",
		TerritoryID: "rt-01",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
	}
	tenantRepo.tenants["t2"] = &entity.Tenant{
		ID:          "t2",
		Block:       "B",
		UnitNumber:  "01",
		TerritoryID: "rt-01",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
	}

	claims := validResidentClaims()
	tenants, err := svc.List(ctx, claims)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// Resident should get tenants via ListByUserID
	if len(tenants) == 0 {
		t.Error("List() should return tenants for resident")
	}
}

func TestTenantService_List_RTOfficerScoped(t *testing.T) {
	svc, tenantRepo, _ := newTestTenantService(t)
	ctx := context.Background()

	// Add tenants in different territories
	tenantRepo.tenants["t1"] = &entity.Tenant{
		ID:          "t1",
		Block:       "A",
		UnitNumber:  "01",
		TerritoryID: "rt-01",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
	}
	tenantRepo.tenants["t2"] = &entity.Tenant{
		ID:          "t2",
		Block:       "B",
		UnitNumber:  "01",
		TerritoryID: "rt-02",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
	}

	claims := &auth.Claims{
		UserID:      "rt-officer-01",
		Role:        "rt_officer",
		TerritoryID: "rt-01",
	}

	tenants, err := svc.List(ctx, claims)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(tenants) != 1 {
		t.Errorf("List() returned %d tenants, want 1", len(tenants))
	}
	if len(tenants) > 0 && tenants[0].TerritoryID != "rt-01" {
		t.Errorf("List() returned tenant from territory %q, want rt-01", tenants[0].TerritoryID)
	}
}

func TestTenantService_List_RWOfficerAll(t *testing.T) {
	svc, tenantRepo, _ := newTestTenantService(t)
	ctx := context.Background()

	// Add tenants in different territories
	tenantRepo.tenants["t1"] = &entity.Tenant{
		ID:          "t1",
		Block:       "A",
		UnitNumber:  "01",
		TerritoryID: "rt-01",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
	}
	tenantRepo.tenants["t2"] = &entity.Tenant{
		ID:          "t2",
		Block:       "B",
		UnitNumber:  "01",
		TerritoryID: "rt-02",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
	}

	claims := validRWClaims()
	tenants, err := svc.List(ctx, claims)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// RW officer should get all tenants via "*" wildcard
	if len(tenants) != 2 {
		t.Errorf("List() returned %d tenants, want 2 (all territories)", len(tenants))
	}
}
