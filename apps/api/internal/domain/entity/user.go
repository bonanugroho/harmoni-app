package entity

import "time"

// User represents a registered user in the Harmoni system.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	TerritoryID  string
	FullName     string
	Phone        string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
