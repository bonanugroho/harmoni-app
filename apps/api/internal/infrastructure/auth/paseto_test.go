package auth

import (
	"testing"
	"time"
)

const testHexKey = "0000000000000000000000000000000000000000000000000000000000000000"

func newTestService(t *testing.T) *PasetoService {
	t.Helper()
	svc, err := NewPasetoService(testHexKey)
	if err != nil {
		t.Fatalf("failed to create paseto service: %v", err)
	}
	return svc
}

func TestPaseto_GenerateAndValidateToken(t *testing.T) {
	svc := newTestService(t)

	token, err := svc.GenerateToken("user-123", "rt_officer", "rt-01", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user-123")
	}
	if claims.Role != "rt_officer" {
		t.Errorf("Role = %q, want %q", claims.Role, "rt_officer")
	}
	if claims.TerritoryID != "rt-01" {
		t.Errorf("TerritoryID = %q, want %q", claims.TerritoryID, "rt-01")
	}
}

func TestPaseto_ExpiredToken(t *testing.T) {
	svc := newTestService(t)

	token, err := svc.GenerateToken("user-123", "resident", "rt-01", -time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = svc.ValidateToken(token)
	if err == nil {
		t.Fatal("ValidateToken() expected error for expired token, got nil")
	}
	if err.Error() != "token expired" {
		t.Errorf("error = %q, want %q", err.Error(), "token expired")
	}
}

func TestPaseto_InvalidToken(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.ValidateToken("this-is-not-a-valid-paseto-token")
	if err == nil {
		t.Fatal("ValidateToken() expected error for invalid token, got nil")
	}
	if err.Error() != "invalid token" {
		t.Errorf("error = %q, want %q", err.Error(), "invalid token")
	}
}

func TestPaseto_TamperedToken(t *testing.T) {
	svc := newTestService(t)

	token, err := svc.GenerateToken("user-123", "resident", "rt-01", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Tamper with the token
	tampered := token[:len(token)-1] + "X"

	_, err = svc.ValidateToken(tampered)
	if err == nil {
		t.Fatal("ValidateToken() expected error for tampered token, got nil")
	}
}

func TestPaseto_NewServiceInvalidKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"not hex", "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"},
		{"too short", "00000000000000000000000000000000"},
		{"too long", "000000000000000000000000000000000000000000000000000000000000000000"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPasetoService(tt.key)
			if err == nil {
				t.Errorf("NewPasetoService(%q) expected error, got nil", tt.key)
			}
		})
	}
}

func TestPaseto_ValidateTokenConstantTime(t *testing.T) {
	a := []byte("token-a")
	b := []byte("token-a")
	c := []byte("token-b")

	if !ValidateTokenConstantTime(a, b) {
		t.Error("ValidateTokenConstantTime(a, a) should be true")
	}
	if ValidateTokenConstantTime(a, c) {
		t.Error("ValidateTokenConstantTime(a, b) should be false for different tokens")
	}
}

func TestPaseto_TokenClaimsExpiryFuture(t *testing.T) {
	svc := newTestService(t)

	token, err := svc.GenerateToken("user-123", "resident", "rt-01", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}

	// Expiry should be approximately 1 hour from now
	if claims.Expiration.Before(time.Now()) {
		t.Error("Token expiration should be in the future")
	}
}
