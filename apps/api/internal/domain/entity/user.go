package entity

import "time"

// User represents a registered user in the Harmoni system.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash,omitempty"`
	Role         string    `json:"role"`
	TerritoryID  string    `json:"territory_id"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Sanitize returns a copy of the user with the password hash removed.
// Use this when returning user data to clients.
func (u *User) Sanitize() *User {
	return &User{
		ID:          u.ID,
		Email:       u.Email,
		Role:        u.Role,
		TerritoryID: u.TerritoryID,
		FullName:    u.FullName,
		Phone:       u.Phone,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
