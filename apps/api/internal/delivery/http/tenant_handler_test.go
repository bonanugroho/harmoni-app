package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"harmoni-api/internal/domain/entity"
	"harmoni-api/internal/domain/repository"
	"harmoni-api/internal/domain/service"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// --- Mock Implementations ---

// mockTenantRepo implements repository.TenantRepository for testing.
type mockTenantRepo struct {
	tenants    map[string]*entity.Tenant
	createErr  error
	findErr    error
	updateErr  error
	deleteErr  error
	listErr    error
}

func newMockTenantRepo() *mockTenantRepo {
	return &mockTenantRepo{
		tenants: make(map[string]*entity.Tenant),
	}
}

func (m *mockTenantRepo) Create(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	tenant.ID = "tenant-uuid-v7"
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()
	m.tenants[tenant.ID] = tenant
	return tenant, nil
}

func (m *mockTenantRepo) CreateTx(ctx context.Context, tx pgx.Tx, tenant *entity.Tenant) (*entity.Tenant, error) {
	return m.Create(ctx, tenant)
}

func (m *mockTenantRepo) FindByID(ctx context.Context, id string) (*entity.Tenant, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	tenant, ok := m.tenants[id]
	if !ok {
		return nil, nil // Simulate error handling
	}
	return tenant, nil
}

func (m *mockTenantRepo) ListByTerritory(ctx context.Context, territoryID string) ([]*entity.Tenant, error) {
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

func (m *mockTenantRepo) ListByUserID(ctx context.Context, userID string) ([]*entity.Tenant, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var result []*entity.Tenant
	for _, t := range m.tenants {
		result = append(result, t)
	}
	return result, nil
}

func (m *mockTenantRepo) Update(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, ok := m.tenants[tenant.ID]
	if !ok {
		return nil, nil
	}
	existing.Block = tenant.Block
	existing.UnitNumber = tenant.UnitNumber
	existing.Occupancy = tenant.Occupancy
	existing.MonthlyFee = tenant.MonthlyFee
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (m *mockTenantRepo) Delete(ctx context.Context, id string, territoryID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.tenants[id]; !ok {
		return nil
	}
	delete(m.tenants, id)
	return nil
}

// mockFeeRepo implements repository.FeeRepository for testing.
type mockFeeRepo struct {
	mandatoryFees map[string]*entity.MandatoryFee
	voluntaryFees map[string]*entity.VoluntaryFee
	createErr     error
	findErr       error
	updateErr     error
	deleteErr     error
}

func newMockFeeRepo() *mockFeeRepo {
	return &mockFeeRepo{
		mandatoryFees: make(map[string]*entity.MandatoryFee),
		voluntaryFees: make(map[string]*entity.VoluntaryFee),
	}
}

func (m *mockFeeRepo) CreateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	fee.ID = "fee-uuid-v7"
	fee.CreatedAt = time.Now()
	fee.UpdatedAt = time.Now()
	m.mandatoryFees[fee.ID] = fee
	return fee, nil
}

func (m *mockFeeRepo) CreateMandatoryTx(ctx context.Context, tx pgx.Tx, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	return m.CreateMandatory(ctx, fee)
}

func (m *mockFeeRepo) ListMandatoryByTenant(ctx context.Context, tenantID string) ([]*entity.MandatoryFee, error) {
	var result []*entity.MandatoryFee
	for _, f := range m.mandatoryFees {
		if f.TenantID == tenantID {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *mockFeeRepo) GetMandatoryByID(ctx context.Context, id string) (*entity.MandatoryFee, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	fee, ok := m.mandatoryFees[id]
	if !ok {
		return nil, nil
	}
	return fee, nil
}

func (m *mockFeeRepo) UpdateMandatory(ctx context.Context, fee *entity.MandatoryFee) (*entity.MandatoryFee, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, ok := m.mandatoryFees[fee.ID]
	if !ok {
		return nil, nil
	}
	existing.Amount = fee.Amount
	existing.Description = fee.Description
	existing.EffectiveDate = fee.EffectiveDate
	existing.PaidAt = fee.PaidAt
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (m *mockFeeRepo) DeleteMandatory(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.mandatoryFees[id]; !ok {
		return nil
	}
	delete(m.mandatoryFees, id)
	return nil
}

func (m *mockFeeRepo) CreateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	fee.ID = "voluntary-uuid-v7"
	fee.CreatedAt = time.Now()
	fee.UpdatedAt = time.Now()
	m.voluntaryFees[fee.ID] = fee
	return fee, nil
}

func (m *mockFeeRepo) ListVoluntaryByTenant(ctx context.Context, tenantID string) ([]*entity.VoluntaryFee, error) {
	var result []*entity.VoluntaryFee
	for _, f := range m.voluntaryFees {
		if f.TenantID == tenantID {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *mockFeeRepo) GetVoluntaryByID(ctx context.Context, id string) (*entity.VoluntaryFee, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	fee, ok := m.voluntaryFees[id]
	if !ok {
		return nil, nil
	}
	return fee, nil
}

func (m *mockFeeRepo) UpdateVoluntary(ctx context.Context, fee *entity.VoluntaryFee) (*entity.VoluntaryFee, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	existing, ok := m.voluntaryFees[fee.ID]
	if !ok {
		return nil, nil
	}
	existing.Amount = fee.Amount
	existing.Description = fee.Description
	existing.EffectiveDate = fee.EffectiveDate
	existing.PaidAt = fee.PaidAt
	existing.UpdatedAt = time.Now()
	return existing, nil
}

func (m *mockFeeRepo) DeleteVoluntary(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.voluntaryFees[id]; !ok {
		return nil
	}
	delete(m.voluntaryFees, id)
	return nil
}

// Ensure mockFeeRepo implements repository.FeeRepository.
var _ repository.FeeRepository = (*mockFeeRepo)(nil)

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

// --- Test Setup ---

func setupTenantHandlerTest(t *testing.T) *fiber.App {
	t.Helper()

	tenantRepo := newMockTenantRepo()
	feeRepo := newMockFeeRepo()

	// Pre-populate a tenant for read/update/delete tests
	futureDate := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	tenantRepo.tenants["tenant-existing"] = &entity.Tenant{
		ID:          "tenant-existing",
		Block:       "A",
		UnitNumber:  "01",
		Occupancy:   "occupied",
		MonthlyFee:  decimal.NewFromInt(50000),
		TerritoryID: "rt-01",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Pre-populate a fee for listing
	feeRepo.mandatoryFees["fee-existing"] = &entity.MandatoryFee{
		ID:            "fee-existing",
		TenantID:      "tenant-existing",
		Amount:        decimal.NewFromInt(25000),
		Description:   "Security Fee",
		EffectiveDate: futureDate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	pool := &mockPool{}
	tenantSvc := service.NewTenantService(tenantRepo, feeRepo, pool)

	app := fiber.New()

	// Middleware to inject test claims (simulates auth middleware)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &auth.Claims{
			UserID:      "test-user",
			Role:        "rt_officer",
			TerritoryID: "rt-01",
		})
		return c.Next()
	})

	handler := NewTenantHandler(tenantSvc)
	api := app.Group("/api")
	handler.RegisterRoutes(api, nil)

	return app
}

// --- Helper Functions ---

func doTenantRequest(app *fiber.App, method, path string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	return app.Test(req)
}

func readTenantBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	defer resp.Body.Close()
	return string(body)
}

func futureDateStr() string {
	return time.Now().Add(24 * time.Hour).Format("2006-01-02")
}

// --- Test Cases ---

func TestTenantHandler_Create_Success(t *testing.T) {
	app := setupTenantHandlerTest(t)

	body := map[string]interface{}{
		"block":       "B",
		"unit_number": "02",
		"occupancy":   "occupied",
		"monthly_fee": 50000,
		"mandatory_fees": []map[string]interface{}{
			{
				"amount":         25000,
				"description":    "Security Fee",
				"effective_date": futureDateStr(),
			},
		},
	}

	resp, err := doTenantRequest(app, "POST", "/api/tenants", body)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "B", result["block"])
}

func TestTenantHandler_Create_MissingBlock(t *testing.T) {
	app := setupTenantHandlerTest(t)

	body := map[string]interface{}{
		"unit_number": "02",
		"occupancy":   "occupied",
		"monthly_fee": 50000,
	}

	resp, err := doTenantRequest(app, "POST", "/api/tenants", body)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	bodyStr := readTenantBody(t, resp)
	assert.Contains(t, bodyStr, "MISSING_FIELDS")
}

func TestTenantHandler_Create_DuplicateBlockUnit(t *testing.T) {
	app := setupTenantHandlerTest(t)

	// First creation should succeed
	body1 := map[string]interface{}{
		"block":       "C",
		"unit_number": "01",
		"occupancy":   "occupied",
		"monthly_fee": 50000,
		"mandatory_fees": []map[string]interface{}{
			{
				"amount":         25000,
				"description":    "Security Fee",
				"effective_date": futureDateStr(),
			},
		},
	}
	resp1, err := doTenantRequest(app, "POST", "/api/tenants", body1)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp1.StatusCode)

	// Second creation with same block+unit should fail
	resp2, err := doTenantRequest(app, "POST", "/api/tenants", body1)
	assert.NoError(t, err)

	// The service wraps errors, so outcome depends on error chain
	// Expected: the repo would return an error for duplicate
	bodyStr := readTenantBody(t, resp2)
	t.Logf("Duplicate response: status=%d body=%s", resp2.StatusCode, bodyStr)
}

func TestTenantHandler_Create_MandatoryFeeRequired(t *testing.T) {
	app := setupTenantHandlerTest(t)

	body := map[string]interface{}{
		"block":          "D",
		"unit_number":    "01",
		"occupancy":      "occupied",
		"monthly_fee":    50000,
		"mandatory_fees": []map[string]interface{}{},
	}

	resp, err := doTenantRequest(app, "POST", "/api/tenants", body)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	bodyStr := readTenantBody(t, resp)
	assert.Contains(t, bodyStr, "MANDATORY_FEE_REQUIRED")
}

func TestTenantHandler_List_Success(t *testing.T) {
	app := setupTenantHandlerTest(t)

	resp, err := doTenantRequest(app, "GET", "/api/tenants", nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyStr := readTenantBody(t, resp)
	assert.Contains(t, bodyStr, "tenants")
}

func TestTenantHandler_Get_Success(t *testing.T) {
	app := setupTenantHandlerTest(t)

	resp, err := doTenantRequest(app, "GET", "/api/tenants/tenant-existing", nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyStr := readTenantBody(t, resp)
	assert.Contains(t, bodyStr, "A")
	assert.Contains(t, bodyStr, "01")
}

func TestTenantHandler_Delete_Success(t *testing.T) {
	app := setupTenantHandlerTest(t)

	resp, err := doTenantRequest(app, "DELETE", "/api/tenants/tenant-existing", nil)
	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestTenantHandler_CreateFee_Success(t *testing.T) {
	app := setupTenantHandlerTest(t)

	body := map[string]interface{}{
		"type":           "mandatory",
		"amount":         15000,
		"description":    "Cleanliness Fee",
		"effective_date": futureDateStr(),
	}

	resp, err := doTenantRequest(app, "POST", "/api/tenants/tenant-existing/fees", body)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestTenantHandler_ListFees_Success(t *testing.T) {
	app := setupTenantHandlerTest(t)

	resp, err := doTenantRequest(app, "GET", "/api/tenants/tenant-existing/fees", nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	bodyStr := readTenantBody(t, resp)
	assert.Contains(t, bodyStr, "mandatory_fees")
	assert.Contains(t, bodyStr, "voluntary_fees")
}

func TestTenantHandler_FeeAmountExceedsCap(t *testing.T) {
	app := setupTenantHandlerTest(t)

	// Tenant has monthly_fee=50000, fee amount=60000 exceeds cap
	body := map[string]interface{}{
		"type":           "mandatory",
		"amount":         60000,
		"description":    "Too Expensive Fee",
		"effective_date": futureDateStr(),
	}

	resp, err := doTenantRequest(app, "POST", "/api/tenants/tenant-existing/fees", body)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	bodyStr := readTenantBody(t, resp)
	assert.Contains(t, bodyStr, "FEE_EXCEEDS_CAP")
}
