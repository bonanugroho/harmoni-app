package auth

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrPasswordNoUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber    = errors.New("password must contain at least one number")
	ErrPasswordNoSymbol    = errors.New("password must contain at least one symbol")
)

// passwordComplexityRegex validates password complexity requirements.
// Using individual checks for clearer error messages.
var (
	hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
	hasLower   = regexp.MustCompile(`[a-z]`).MatchString
	hasDigit   = regexp.MustCompile(`[0-9]`).MatchString
	hasSpecial = regexp.MustCompile(`[^A-Za-z0-9]`).MatchString
)

// HashPassword hashes a password using bcrypt with DefaultCost (10).
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword compares a bcrypt hash with a plaintext password.
// Returns nil if they match, or an error if they don't.
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ValidatePassword checks password complexity requirements:
// - At least 8 characters
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one number
// - At least one symbol
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	hasUpperChar := false
	hasLowerChar := false
	hasDigitChar := false
	hasSpecialChar := false

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpperChar = true
		case unicode.IsLower(r):
			hasLowerChar = true
		case unicode.IsDigit(r):
			hasDigitChar = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecialChar = true
		}
	}

	if !hasUpperChar {
		return ErrPasswordNoUppercase
	}
	if !hasLowerChar {
		return ErrPasswordNoLowercase
	}
	if !hasDigitChar {
		return ErrPasswordNoNumber
	}
	if !hasSpecialChar {
		return ErrPasswordNoSymbol
	}

	return nil
}
