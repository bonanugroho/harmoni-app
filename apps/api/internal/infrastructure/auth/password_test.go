package auth

import (
	"testing"
)

func TestPassword_HashAndCompare(t *testing.T) {
	password := "SecurePass123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}

	// Compare correct password
	err = ComparePassword(hash, password)
	if err != nil {
		t.Errorf("ComparePassword() error = %v, want nil", err)
	}

	// Compare wrong password
	err = ComparePassword(hash, "WrongPassword1!")
	if err == nil {
		t.Fatal("ComparePassword() expected error for wrong password, got nil")
	}
}

func TestPassword_DifferentHashesForSamePassword(t *testing.T) {
	password := "SecurePass123!"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword() returned same hash for same password — salt should be random")
	}

	// Both should still validate
	if err := ComparePassword(hash1, password); err != nil {
		t.Errorf("ComparePassword(hash1) error = %v", err)
	}
	if err := ComparePassword(hash2, password); err != nil {
		t.Errorf("ComparePassword(hash2) error = %v", err)
	}
}

func TestPassword_ValidatePassword_Accepts(t *testing.T) {
	validPasswords := []string{
		"SecurePass123!",
		"MyP@ssw0rd",
		"Abcdef1!",
		"Test1234$%^&",
	}

	for _, pw := range validPasswords {
		t.Run(pw, func(t *testing.T) {
			err := ValidatePassword(pw)
			if err != nil {
				t.Errorf("ValidatePassword(%q) error = %v, want nil", pw, err)
			}
		})
	}
}

func TestPassword_ValidatePassword_Rejects(t *testing.T) {
	tests := []struct {
		password string
		wantErr  error
	}{
		{"password", ErrPasswordNoUppercase},
		{"short1!", ErrPasswordTooShort},
		{"NoSpecial1", ErrPasswordNoSymbol},
		{"nouppercase1!", ErrPasswordNoUppercase},
		{"", ErrPasswordTooShort},
		{"abc", ErrPasswordTooShort},
		{"NOLOWER1!", ErrPasswordNoLowercase},
		{"NoDigits!", ErrPasswordNoNumber},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if err == nil {
				t.Errorf("ValidatePassword(%q) expected error, got nil", tt.password)
				return
			}
			if err != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, want %v", tt.password, err, tt.wantErr)
			}
		})
	}
}
