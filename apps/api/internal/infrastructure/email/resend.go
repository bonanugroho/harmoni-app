package email

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

// ResendEmailService implements the EmailService interface using Resend API.
type ResendEmailService struct {
	client    *resend.Client
	fromEmail string
}

// NewResendEmailService creates a new email service using the Resend API.
func NewResendEmailService(apiKey, fromEmail string) *ResendEmailService {
	client := resend.NewClient(apiKey)
	return &ResendEmailService{
		client:    client,
		fromEmail: fromEmail,
	}
}

// SendPasswordResetEmail sends a password reset email with the given reset URL.
func (s *ResendEmailService) SendPasswordResetEmail(to, resetURL string) error {
	params := &resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{to},
		Subject: "Reset Your Harmoni Password",
		Html:    passwordResetHTML(resetURL),
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

// passwordResetHTML returns the HTML template for password reset emails.
func passwordResetHTML(resetURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
    .container { max-width: 480px; margin: 0 auto; background: #ffffff; border-radius: 8px; padding: 32px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
    .logo { font-size: 24px; font-weight: bold; color: #2563eb; margin-bottom: 24px; }
    h1 { font-size: 20px; color: #1f2937; margin: 0 0 16px; }
    p { color: #6b7280; line-height: 1.6; margin: 0 0 24px; }
    .button { display: inline-block; background: #2563eb; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 6px; font-weight: 500; }
    .expiry { font-size: 12px; color: #9ca3af; margin-top: 24px; padding-top: 16px; border-top: 1px solid #e5e7eb; }
    .footer { font-size: 12px; color: #9ca3af; margin-top: 16px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="logo">Harmoni</div>
    <h1>Reset Your Password</h1>
    <p>We received a request to reset your password. Click the button below to create a new password:</p>
    <p><a href="%s" class="button">Reset Password</a></p>
    <p>If you didn't request this, you can safely ignore this email.</p>
    <div class="expiry">
      <strong>⏰ This link expires in 1 hour.</strong> If it expires, please request a new password reset.
    </div>
    <div class="footer">
      Harmoni — Community Financial Management<br>
      This is an automated message. Please do not reply.
    </div>
  </div>
</body>
</html>
`, resetURL)
}
