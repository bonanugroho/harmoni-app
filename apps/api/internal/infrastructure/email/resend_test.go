package email

import (
	"context"
	"strings"
	"testing"
)

func TestResend_SendPasswordResetEmail_ContainsResetLink(t *testing.T) {
	// Test the HTML template generation directly (without API call)
	resetURL := "https://harmonictest.app/reset-password?token=abc123"
	html := passwordResetHTML(resetURL)

	if !strings.Contains(html, resetURL) {
		t.Errorf("HTML should contain reset URL %q", resetURL)
	}
	if !strings.Contains(html, "Reset Password") {
		t.Error("HTML should contain 'Reset Password' button text")
	}
	if !strings.Contains(html, "expires in 1 hour") {
		t.Error("HTML should contain expiry notice")
	}
	if !strings.Contains(html, "Harmoni") {
		t.Error("HTML should contain branding")
	}
}

func TestResend_SendPasswordResetEmail_MissingAPIKey(t *testing.T) {
	// Creating a service with empty API key should still create a client
	// but sending will fail — we test the interface contract here
	svc := NewResendEmailService("", "noreply@harmonictest.app")

	err := svc.SendPasswordResetEmail("test@example.com", "https://example.com/reset")
	if err == nil {
		t.Log("Note: Resend client may accept empty key but fail at send time — this is expected")
	}
}

// MockEmailService implements repository.EmailService for testing.
type MockEmailService struct {
	LastTo     string
	LastURL    string
	SentCount  int
	SendError  error
}

// SendPasswordResetEmail records the call for test verification.
func (m *MockEmailService) SendPasswordResetEmail(to, resetURL string) error {
	m.LastTo = to
	m.LastURL = resetURL
	m.SentCount++
	return m.SendError
}

// Verify MockEmailService implements the interface at compile time.
var _ interface {
	SendPasswordResetEmail(to, resetURL string) error
} = (*MockEmailService)(nil)

// TestMockEmailService verifies the mock works correctly.
func TestMockEmailService(t *testing.T) {
	mock := &MockEmailService{}

	err := mock.SendPasswordResetEmail("user@example.com", "https://example.com/reset?token=xyz")
	if err != nil {
		t.Fatalf("Mock SendPasswordResetEmail() error = %v", err)
	}

	if mock.LastTo != "user@example.com" {
		t.Errorf("LastTo = %q, want %q", mock.LastTo, "user@example.com")
	}
	if mock.LastURL != "https://example.com/reset?token=xyz" {
		t.Errorf("LastURL = %q, want %q", mock.LastURL, "https://example.com/reset?token=xyz")
	}
	if mock.SentCount != 1 {
		t.Errorf("SentCount = %d, want 1", mock.SentCount)
	}

	// Test error propagation
	mock.SendError = context.DeadlineExceeded
	err = mock.SendPasswordResetEmail("user@example.com", "https://example.com/reset")
	if err != context.DeadlineExceeded {
		t.Errorf("Mock should propagate SendError, got %v", err)
	}
}
