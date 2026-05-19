package auth

import (
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"time"

	"github.com/o1egl/paseto"
)

// Claims holds the data embedded in a PASETO token.
type Claims struct {
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"`
	TerritoryID string    `json:"territory_id"`
	Expiration  time.Time `json:"exp"`
}

// PasetoService handles PASETO V2 Local token generation and validation.
// V2 uses XChaCha20-Poly1305 (same algorithm as V4 Local).
type PasetoService struct {
	v2  *paseto.V2
	key []byte
}

// NewPasetoService creates a new PASETO service from a hex-encoded 32-byte key.
func NewPasetoService(hexKey string) (*PasetoService, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, errors.New("invalid PASETO secret key: not valid hex")
	}
	if len(key) != 32 {
		return nil, errors.New("invalid PASETO secret key: must be 32 bytes")
	}

	return &PasetoService{
		v2:  paseto.NewV2(),
		key: key,
	}, nil
}

// GenerateToken creates a PASETO V2 Local token with the given claims.
func (s *PasetoService) GenerateToken(userID, role, territoryID string, expiry time.Duration) (string, error) {
	claims := Claims{
		UserID:      userID,
		Role:        role,
		TerritoryID: territoryID,
		Expiration:  time.Now().Add(expiry),
	}

	token, err := s.v2.Encrypt(s.key, claims, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken decrypts and validates a PASETO token, returning the claims.
func (s *PasetoService) ValidateToken(token string) (*Claims, error) {
	var claims Claims

	err := s.v2.Decrypt(token, s.key, &claims, nil)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check expiry manually
	if time.Now().After(claims.Expiration) {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// ValidateTokenConstantTime compares two tokens using constant-time comparison.
func ValidateTokenConstantTime(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
