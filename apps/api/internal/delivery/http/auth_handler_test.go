package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"harmoni-api/internal/domain/entity"
	"harmoni-api/internal/domain/repository"
	"harmoni-api/internal/domain/service"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
)

const testHexKey = "0000000000000000000000000000000000000000000000000000000000000000"

// mockUserRepo implements repository.UserRepository for testing.
type mockUserRepo struct {
	users      map[string]*entity.User
	emailIndex map[string]*entity.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:      make(map[string]*entity.User),
		emailIndex: make(map[string]*entity.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	user.ID = "test-uuid-v7"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user
	return user, nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, ok := m.emailIndex[email]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, id, hash string) error {
	user, ok := m.users[id]
	if !ok {
		return sql.ErrNoRows
	}
	user.PasswordHash = hash
	return nil
}

func (m *mockUserRepo) ListByTerritory(ctx context.Context, territoryID string) ([]*entity.User, error) {
	return nil, nil
}

// mockResetRepo implements repository.PasswordResetTokenRepository.
type mockResetRepo struct {
	tokens map[string]*repository.PasswordResetToken
}

func newMockResetRepo() *mockResetRepo {
	return &mockResetRepo{tokens: make(map[string]*repository.PasswordResetToken)}
}

func (m *mockResetRepo) Create(ctx context.Context, token *repository.PasswordResetToken) error {
	token.ID = "reset-uuid-v7"
	token.CreatedAt = time.Now()
	m.tokens[token.TokenHash] = token
	return nil
}

func (m *mockResetRepo) FindByTokenHash(ctx context.Context, hash string) (*repository.PasswordResetToken, error) {
	token, ok := m.tokens[hash]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return token, nil
}

func (m *mockResetRepo) MarkUsed(ctx context.Context, id string) error {
	for _, t := range m.tokens {
		if t.ID == id {
			t.Used = true
			return nil
		}
	}
	return sql.ErrNoRows
}

func (m *mockResetRepo) DeleteByUserID(ctx context.Context, userID string) error {
	return nil
}

// mockEmail implements repository.EmailService.
type mockEmail struct {
	LastTo  string
	LastURL string
}

func (m *mockEmail) SendPasswordResetEmail(to, resetURL string) error {
	m.LastTo = to
	m.LastURL = resetURL
	return nil
}

func setupTestHandler(t *testing.T) (*fiber.App, *mockEmail) {
	t.Helper()

	userRepo := newMockUserRepo()
	resetRepo := newMockResetRepo()
	emailSvc := &mockEmail{}

	paseto, err := auth.NewPasetoService(testHexKey)
	if err != nil {
		t.Fatalf("failed to create paseto service: %v", err)
	}

	authSvc := service.NewAuthService(userRepo, resetRepo, paseto, emailSvc, "https://harmonictest.app")

	app := fiber.New()
	handler := NewAuthHandler(authSvc)
	handler.RegisterRoutes(app)

	return app, emailSvc
}

func TestAuthHandler_Register_Success(t *testing.T) {
	app, _ := setupTestHandler(t)

	body := map[string]string{
		"email":        "test@example.com",
		"password":     "SecurePass123!",
		"full_name":    "Test User",
		"territory_id": "rt-01",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 201 {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("StatusCode = %d, want 201, body: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["email"] != "test@example.com" {
		t.Errorf("email = %v, want test@example.com", result["email"])
	}
	if _, ok := result["password_hash"]; ok {
		t.Error("response should not contain password_hash")
	}
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	app, _ := setupTestHandler(t)

	body := map[string]string{
		"email":    "",
		"password": "SecurePass123!",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func TestAuthHandler_Register_WeakPassword(t *testing.T) {
	app, _ := setupTestHandler(t)

	body := map[string]string{
		"email":        "test@example.com",
		"password":     "weak",
		"full_name":    "Test User",
		"territory_id": "rt-01",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	app, _ := setupTestHandler(t)

	body := map[string]string{
		"email":        "dup@example.com",
		"password":     "SecurePass123!",
		"full_name":    "Test User",
		"territory_id": "rt-01",
	}
	jsonBody, _ := json.Marshal(body)

	// First registration
	req1 := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != 201 {
		t.Fatalf("first registration failed with status %d", resp1.StatusCode)
	}

	// Second registration with same email
	req2 := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(jsonBody))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)

	if resp2.StatusCode != 409 {
		t.Errorf("StatusCode = %d, want 409", resp2.StatusCode)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	app, _ := setupTestHandler(t)

	// Register first
	regBody := map[string]string{
		"email":        "login@example.com",
		"password":     "SecurePass123!",
		"full_name":    "Login User",
		"territory_id": "rt-01",
	}
	regJSON, _ := json.Marshal(regBody)
	regReq := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(regJSON))
	regReq.Header.Set("Content-Type", "application/json")
	regResp, _ := app.Test(regReq)
	if regResp.StatusCode != 201 {
		t.Fatalf("registration failed with status %d", regResp.StatusCode)
	}

	// Login
	loginBody := map[string]string{
		"email":    "login@example.com",
		"password": "SecurePass123!",
	}
	loginJSON, _ := json.Marshal(loginBody)
	loginReq := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(loginJSON))
	loginReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(loginReq)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("StatusCode = %d, want 200, body: %s", resp.StatusCode, string(respBody))
	}

	// Check for Set-Cookie header
	cookies := resp.Header.Values("Set-Cookie")
	if len(cookies) == 0 {
		t.Fatal("Expected Set-Cookie header, got none")
	}

	cookie := cookies[0]
	if !containsStr(cookie, "paseto_token=") {
		t.Errorf("Cookie should contain 'paseto_token=', got: %s", cookie)
	}
	if !containsStr(cookie, "httponly") {
		t.Errorf("Cookie should be httponly, got: %s", cookie)
	}
	if !containsStr(cookie, "secure") {
		t.Errorf("Cookie should be secure, got: %s", cookie)
	}
	if !containsStr(cookie, "samesite=strict") {
		t.Errorf("Cookie should be samesite=strict, got: %s", cookie)
	}
	if !containsStr(cookie, "path=/") {
		t.Errorf("Cookie should have path=/, got: %s", cookie)
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	app, _ := setupTestHandler(t)

	body := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "SecurePass123!",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", resp.StatusCode)
	}
}

func TestAuthHandler_ResetPasswordRequest_Success(t *testing.T) {
	app, emailMock := setupTestHandler(t)

	// Register first
	regBody := map[string]string{
		"email":        "reset@example.com",
		"password":     "SecurePass123!",
		"full_name":    "Reset User",
		"territory_id": "rt-01",
	}
	regJSON, _ := json.Marshal(regBody)
	regReq := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(regJSON))
	regReq.Header.Set("Content-Type", "application/json")
	regResp, _ := app.Test(regReq)
	if regResp.StatusCode != 201 {
		t.Fatalf("registration failed with status %d", regResp.StatusCode)
	}

	// Request reset
	resetBody := map[string]string{
		"email": "reset@example.com",
	}
	resetJSON, _ := json.Marshal(resetBody)
	resetReq := httptest.NewRequest("POST", "/auth/reset", bytes.NewReader(resetJSON))
	resetReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(resetReq)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}

	if emailMock.LastTo != "reset@example.com" {
		t.Errorf("email LastTo = %q, want %q", emailMock.LastTo, "reset@example.com")
	}
	if emailMock.LastURL == "" {
		t.Error("email LastURL should not be empty")
	}
}

func TestAuthHandler_ResetPasswordRequest_NonExistentEmail(t *testing.T) {
	app, _ := setupTestHandler(t)

	// Should return 200 even for non-existent email (prevent enumeration)
	body := map[string]string{
		"email": "nonexistent@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/reset", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200 (prevent enumeration)", resp.StatusCode)
	}
}

func TestAuthHandler_ResetPasswordConfirm_MissingFields(t *testing.T) {
	app, _ := setupTestHandler(t)

	body := map[string]string{
		"token": "some-token",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/reset/confirm", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func TestAuthHandler_ResetPasswordConfirm_InvalidToken(t *testing.T) {
	app, _ := setupTestHandler(t)

	body := map[string]string{
		"token":        "invalid-token",
		"new_password": "NewSecurePass456!",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/reset/confirm", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode)
	}
}

func containsStr(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			result[i] = byte(c + 32)
		} else {
			result[i] = byte(c)
		}
	}
	return string(result)
}
