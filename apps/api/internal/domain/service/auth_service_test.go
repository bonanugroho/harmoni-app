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
)

const testHexKey = "0000000000000000000000000000000000000000000000000000000000000000"

// mockUserRepository implements repository.UserRepository for testing.
type mockUserRepository struct {
	users      map[string]*entity.User
	emailIndex map[string]*entity.User
	createErr  error
	findErr    error
	updateErr  error
}

func newMockUserRepo() *mockUserRepository {
	return &mockUserRepository{
		users:      make(map[string]*entity.User),
		emailIndex: make(map[string]*entity.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	user.ID = "test-uuid-v7"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user
	return user, nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	user, ok := m.emailIndex[email]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *mockUserRepository) UpdatePassword(ctx context.Context, id, hash string) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	user, ok := m.users[id]
	if !ok {
		return sql.ErrNoRows
	}
	user.PasswordHash = hash
	user.UpdatedAt = time.Now()
	return nil
}

func (m *mockUserRepository) ListByTerritory(ctx context.Context, territoryID string) ([]*entity.User, error) {
	var result []*entity.User
	for _, u := range m.users {
		if u.TerritoryID == territoryID {
			result = append(result, u)
		}
	}
	return result, nil
}

// mockResetTokenRepository implements repository.PasswordResetTokenRepository.
type mockResetTokenRepository struct {
	tokens    map[string]*repository.PasswordResetToken
	createErr error
	findErr   error
	markErr   error
}

func newMockResetRepo() *mockResetTokenRepository {
	return &mockResetTokenRepository{
		tokens: make(map[string]*repository.PasswordResetToken),
	}
}

func (m *mockResetTokenRepository) Create(ctx context.Context, token *repository.PasswordResetToken) error {
	if m.createErr != nil {
		return m.createErr
	}
	token.ID = "reset-uuid-v7"
	token.CreatedAt = time.Now()
	m.tokens[token.TokenHash] = token
	return nil
}

func (m *mockResetTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*repository.PasswordResetToken, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	token, ok := m.tokens[hash]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return token, nil
}

func (m *mockResetTokenRepository) MarkUsed(ctx context.Context, id string) error {
	if m.markErr != nil {
		return m.markErr
	}
	for _, t := range m.tokens {
		if t.ID == id {
			t.Used = true
			return nil
		}
	}
	return sql.ErrNoRows
}

func (m *mockResetTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	for k, t := range m.tokens {
		if t.UserID == userID {
			delete(m.tokens, k)
		}
	}
	return nil
}

// mockEmailService implements repository.EmailService.
type mockEmailService struct {
	lastTo     string
	lastURL    string
	sendErr    error
}

func (m *mockEmailService) SendPasswordResetEmail(to, resetURL string) error {
	m.lastTo = to
	m.lastURL = resetURL
	return m.sendErr
}

func newTestAuthService(t *testing.T) (*AuthService, *mockUserRepository, *mockResetTokenRepository, *mockEmailService) {
	t.Helper()

	userRepo := newMockUserRepo()
	resetRepo := newMockResetRepo()
	emailSvc := &mockEmailService{}

	paseto, err := auth.NewPasetoService(testHexKey)
	if err != nil {
		t.Fatalf("failed to create paseto service: %v", err)
	}

	svc := NewAuthService(userRepo, resetRepo, paseto, emailSvc, "https://harmonictest.app")
	return svc, userRepo, resetRepo, emailSvc
}

func TestAuthService_Register_Success(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	user, err := svc.Register("test@example.com", "SecurePass123!", "Test User", "rt-01")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", user.Email, "test@example.com")
	}
	if user.FullName != "Test User" {
		t.Errorf("FullName = %q, want %q", user.FullName, "Test User")
	}
	if user.TerritoryID != "rt-01" {
		t.Errorf("TerritoryID = %q, want %q", user.TerritoryID, "rt-01")
	}
	if user.Role != "resident" {
		t.Errorf("Role = %q, want %q", user.Role, "resident")
	}
	if user.PasswordHash != "" {
		t.Error("Register() should not return password hash")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	_, err := svc.Register("test@example.com", "SecurePass123!", "Test User", "rt-01")
	if err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	_, err = svc.Register("test@example.com", "SecurePass123!", "Another User", "rt-01")
	if !errors.Is(err, ErrDuplicateEmail) {
		t.Errorf("Register() error = %v, want ErrDuplicateEmail", err)
	}
}

func TestAuthService_Register_WeakPassword(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	_, err := svc.Register("test@example.com", "weak", "Test User", "rt-01")
	if err == nil {
		t.Fatal("Register() expected error for weak password, got nil")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	// Register first
	_, err := svc.Register("test@example.com", "SecurePass123!", "Test User", "rt-01")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Login
	user, token, err := svc.Login("test@example.com", "SecurePass123!")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", user.Email, "test@example.com")
	}
	if token == "" {
		t.Fatal("Login() returned empty token")
	}
	if user.PasswordHash != "" {
		t.Error("Login() should not return password hash")
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	_, err := svc.Register("test@example.com", "SecurePass123!", "Test User", "rt-01")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	_, _, err = svc.Login("test@example.com", "WrongPassword1!")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	_, _, err := svc.Login("nonexistent@example.com", "SecurePass123!")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestAuthService_ResetPasswordRequest_Success(t *testing.T) {
	svc, _, _, emailSvc := newTestAuthService(t)

	_, err := svc.Register("test@example.com", "SecurePass123!", "Test User", "rt-01")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	err = svc.ResetPasswordRequest("test@example.com")
	if err != nil {
		t.Fatalf("ResetPasswordRequest() error = %v", err)
	}

	if emailSvc.lastTo != "test@example.com" {
		t.Errorf("email lastTo = %q, want %q", emailSvc.lastTo, "test@example.com")
	}
	if emailSvc.lastURL == "" {
		t.Error("email lastURL should not be empty")
	}
}

func TestAuthService_ResetPasswordRequest_NonExistentEmail(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	// Should not error — prevents email enumeration
	err := svc.ResetPasswordRequest("nonexistent@example.com")
	if err != nil {
		t.Errorf("ResetPasswordRequest() error = %v, want nil (prevent enumeration)", err)
	}
}

func TestAuthService_ResetPassword_Success(t *testing.T) {
	svc, _, resetRepo, _ := newTestAuthService(t)

	// Register
	_, err := svc.Register("test@example.com", "SecurePass123!", "Test User", "rt-01")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	// Manually create a reset token (bypassing the email flow)
	tokenHash := hashToken("raw-reset-token")
	resetToken := &repository.PasswordResetToken{
		ID:        "reset-uuid-v7",
		UserID:    "test-uuid-v7",
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(time.Hour),
		Used:      false,
	}
	resetRepo.tokens[tokenHash] = resetToken

	// Reset password using the raw token
	err = svc.ResetPassword("raw-reset-token", "NewSecurePass456!")
	if err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}

	// Verify password was updated — login with new password
	user, _, err := svc.Login("test@example.com", "NewSecurePass456!")
	if err != nil {
		t.Fatalf("Login() with new password error = %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", user.Email, "test@example.com")
	}
}

func TestAuthService_ResetPassword_InvalidToken(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	err := svc.ResetPassword("invalid-token", "NewSecurePass456!")
	if !errors.Is(err, ErrInvalidResetToken) {
		t.Errorf("ResetPassword() error = %v, want ErrInvalidResetToken", err)
	}
}

func TestAuthService_ResetPassword_WeakNewPassword(t *testing.T) {
	svc, _, _, _ := newTestAuthService(t)

	err := svc.ResetPassword("some-token", "weak")
	if err == nil {
		t.Fatal("ResetPassword() expected error for weak new password, got nil")
	}
}
