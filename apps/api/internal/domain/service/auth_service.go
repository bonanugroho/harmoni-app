package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"harmoni-api/internal/domain/entity"
	"harmoni-api/internal/domain/repository"
	"harmoni-api/internal/infrastructure/auth"
)

var (
	ErrDuplicateEmail     = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidResetToken  = errors.New("invalid or expired reset token")
	ErrResetTokenUsed     = errors.New("reset token has already been used")
	ErrUserNotFound       = errors.New("user not found")
	ErrInactiveUser       = errors.New("user account is inactive")
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo    repository.UserRepository
	resetRepo   repository.PasswordResetTokenRepository
	paseto      *auth.PasetoService
	email       repository.EmailService
	tokenExpiry time.Duration
	resetExpiry time.Duration
	frontendURL string
}

// NewAuthService creates a new auth service.
func NewAuthService(
	userRepo repository.UserRepository,
	resetRepo repository.PasswordResetTokenRepository,
	paseto *auth.PasetoService,
	email repository.EmailService,
	frontendURL string,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		resetRepo:   resetRepo,
		paseto:      paseto,
		email:       email,
		tokenExpiry: time.Hour,
		resetExpiry: time.Hour,
		frontendURL: frontendURL,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(email, password, fullName, territoryID string) (*entity.User, error) {
	// Validate password complexity
	if err := auth.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Check if email already exists
	_, err := s.userRepo.FindByEmail(nil, email)
	if err == nil {
		return nil, ErrDuplicateEmail
	}

	// Hash password
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entity.User{
		Email:        email,
		PasswordHash: hash,
		Role:         "resident", // Default role
		TerritoryID:  territoryID,
		FullName:     fullName,
		IsActive:     true,
	}

	created, err := s.userRepo.Create(nil, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return created.Sanitize(), nil
}

// Login authenticates a user and returns a PASETO token.
func (s *AuthService) Login(email, password string) (*entity.User, string, error) {
	user, err := s.userRepo.FindByEmail(nil, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", ErrInactiveUser
	}

	if err := auth.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate token
	token, err := s.paseto.GenerateToken(user.ID, user.Role, user.TerritoryID, s.tokenExpiry)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user.Sanitize(), token, nil
}

// hashToken creates a SHA-256 hash of a token for storage/lookup.
// SHA-256 is used (not bcrypt) because we need deterministic lookup.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// ResetPasswordRequest initiates a password reset flow.
func (s *AuthService) ResetPasswordRequest(email string) error {
	// Find user — but always return success to prevent email enumeration
	user, err := s.userRepo.FindByEmail(nil, email)
	if err != nil {
		// User not found — still return nil to prevent enumeration
		return nil
	}

	// Generate reset token
	rawToken, err := generateRandomToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Hash the token for storage (SHA-256 for deterministic lookup)
	tokenHash := hashToken(rawToken)

	// Store reset token
	resetToken := &repository.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.resetExpiry),
		Used:      false,
	}

	if err := s.resetRepo.Create(nil, resetToken); err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Send email
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, rawToken)
	if err := s.email.SendPasswordResetEmail(user.Email, resetURL); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword completes a password reset using a token.
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Validate new password
	if err := auth.ValidatePassword(newPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Hash the incoming token to find the record (SHA-256 for deterministic lookup)
	tokenHash := hashToken(token)

	// Find the reset token record
	record, err := s.resetRepo.FindByTokenHash(nil, tokenHash)
	if err != nil {
		return ErrInvalidResetToken
	}

	if record.Used {
		return ErrResetTokenUsed
	}

	// Check if token has expired
	if time.Now().After(record.ExpiresAt) {
		return ErrInvalidResetToken
	}

	// Hash new password
	newHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(nil, record.UserID, newHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.resetRepo.MarkUsed(nil, record.ID); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}

// generateRandomToken creates a 32-byte random hex token.
func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
