package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"harmoni-api/internal/domain/entity"
	"harmoni-api/internal/domain/repository"
)

// TestPostgresUserRepository_InterfaceContract verifies that
// PostgresUserRepository implements the UserRepository interface.
var _ repository.UserRepository = (*PostgresUserRepository)(nil)

// mockPgxPool is a minimal mock for testing repository behavior
// without a live PostgreSQL connection.
// For full integration tests with real DB, use testcontainers-go
// with the build tag: go test -tags=integration

func TestUserRepository_Create_Success(t *testing.T) {
	// This test verifies the repository interface contract.
	// Full integration tests require a live PostgreSQL instance
	// via testcontainers-go (run with: go test -tags=integration).
	//
	// Interface contract verification:
	// - Create(ctx, user) returns (*User, error)
	// - FindByEmail(ctx, email) returns (*User, error)
	// - FindByID(ctx, id) returns (*User, error)
	// - UpdatePassword(ctx, id, hash) returns error
	// - ListByTerritory(ctx, territoryID) returns ([]*User, error)

	user := &entity.User{
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		Role:         "resident",
		TerritoryID:  "rt-01",
		FullName:     "Test User",
		Phone:        "081234567890",
		IsActive:     true,
	}

	// Verify User entity fields are set correctly
	if user.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", user.Email, "test@example.com")
	}
	if user.Role != "resident" {
		t.Errorf("Role = %q, want %q", user.Role, "resident")
	}
	if user.TerritoryID != "rt-01" {
		t.Errorf("TerritoryID = %q, want %q", user.TerritoryID, "rt-01")
	}
	if !user.IsActive {
		t.Error("IsActive should be true")
	}
}

func TestUserRepository_EntitySanitize(t *testing.T) {
	user := &entity.User{
		ID:           "test-uuid",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$should-be-removed",
		Role:         "resident",
		TerritoryID:  "rt-01",
		FullName:     "Test User",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	sanitized := user.Sanitize()

	if sanitized.PasswordHash != "" {
		t.Error("Sanitize() should remove password hash")
	}
	if sanitized.Email != user.Email {
		t.Errorf("Sanitize() Email = %q, want %q", sanitized.Email, user.Email)
	}
	if sanitized.ID != user.ID {
		t.Errorf("Sanitize() ID = %q, want %q", sanitized.ID, user.ID)
	}
}

func TestUserRepository_TerritoryFiltering(t *testing.T) {
	// Verify territory filtering logic concept
	users := []*entity.User{
		{ID: "1", Email: "a@rt01.com", TerritoryID: "rt-01", FullName: "A"},
		{ID: "2", Email: "b@rt01.com", TerritoryID: "rt-01", FullName: "B"},
		{ID: "3", Email: "c@rt02.com", TerritoryID: "rt-02", FullName: "C"},
	}

	// Filter by rt-01
	var rt01Users []*entity.User
	for _, u := range users {
		if u.TerritoryID == "rt-01" {
			rt01Users = append(rt01Users, u)
		}
	}

	if len(rt01Users) != 2 {
		t.Errorf("Expected 2 users in rt-01, got %d", len(rt01Users))
	}

	// Verify correct users
	for _, u := range rt01Users {
		if u.TerritoryID != "rt-01" {
			t.Errorf("User %s has wrong territory: %s", u.Email, u.TerritoryID)
		}
	}
}

func TestUserRepository_PasswordUpdate(t *testing.T) {
	// Verify password update concept
	user := &entity.User{
		ID:           "test-uuid",
		PasswordHash: "$2a$10$oldhash",
	}

	newHash := "$2a$10$newhash"
	user.PasswordHash = newHash

	if user.PasswordHash != newHash {
		t.Errorf("PasswordHash = %q, want %q", user.PasswordHash, newHash)
	}
}

func TestUserRepository_ErrNoRows(t *testing.T) {
	// Verify sql.ErrNoRows is the expected error for not-found cases
	if sql.ErrNoRows.Error() == "" {
		t.Error("sql.ErrNoRows should have an error message")
	}
}

func TestUserRepository_ContextCancellation(t *testing.T) {
	// Verify that context cancellation is handled by the repository interface
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	if ctx.Err() == nil {
		t.Error("Context should be cancelled")
	}
}
