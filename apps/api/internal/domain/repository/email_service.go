package repository

// EmailService defines the interface for sending transactional emails.
type EmailService interface {
	// SendPasswordResetEmail sends a password reset email with the given reset URL.
	SendPasswordResetEmail(to, resetURL string) error
}
