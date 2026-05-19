package http

import (
	"log"
	"time"

	"harmoni-api/internal/domain/service"
	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	authService  *service.AuthService
	pasetoService *auth.PasetoService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService *service.AuthService, pasetoService *auth.PasetoService) *AuthHandler {
	return &AuthHandler{authService: authService, pasetoService: pasetoService}
}

// RegisterRoutes registers auth endpoints on the Fiber app.
func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/auth")
	api.Post("/register", h.Register)
	api.Post("/login", h.Login)
	api.Post("/reset", h.ResetPasswordRequest)
	api.Post("/reset/confirm", h.ResetPasswordConfirm)
	api.Get("/me", h.Me)
}

// RegisterRequest represents the registration request body.
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	TerritoryID string `json:"territory_id"`
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST",
		})
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FullName == "" || req.TerritoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email, password, full_name, and territory_id are required",
			"code":  "MISSING_FIELDS",
		})
	}

	user, err := h.authService.Register(req.Email, req.Password, req.FullName, req.TerritoryID)
	if err != nil {
		switch {
		case err == service.ErrDuplicateEmail:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already registered",
				"code":  "DUPLICATE_EMAIL",
			})
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"code":  "REGISTRATION_FAILED",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email and password are required",
			"code":  "MISSING_FIELDS",
		})
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case err == service.ErrInvalidCredentials, err == service.ErrUserNotFound:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
				"code":  "INVALID_CREDENTIALS",
			})
		case err == service.ErrInactiveUser:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Account is inactive",
				"code":  "INACTIVE_ACCOUNT",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
				"code":  "INTERNAL_ERROR",
			})
		}
	}

	// Set httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "paseto_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   3600, // 1 hour
	})

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// ResetRequest represents the password reset request body.
type ResetRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest handles POST /auth/reset.
func (h *AuthHandler) ResetPasswordRequest(c *fiber.Ctx) error {
	var req ResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is required",
			"code":  "MISSING_FIELDS",
		})
	}

	// Always return 200 to prevent email enumeration
	err := h.authService.ResetPasswordRequest(req.Email)
	if err != nil {
		// Log the actual error for debugging
		log.Printf("Password reset request failed for %s: %v", req.Email, err)
		// Still return 200 to prevent email enumeration
		return c.JSON(fiber.Map{
			"message": "If an account exists with that email, a password reset link has been sent",
		})
	}

	return c.JSON(fiber.Map{
		"message": "If an account exists with that email, a password reset link has been sent",
	})
}

// ResetConfirmRequest represents the password reset confirmation body.
type ResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// ResetPasswordConfirm handles POST /auth/reset/confirm.
func (h *AuthHandler) ResetPasswordConfirm(c *fiber.Ctx) error {
	var req ResetConfirmRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "INVALID_REQUEST",
		})
	}

	if req.Token == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "token and new_password are required",
			"code":  "MISSING_FIELDS",
		})
	}

	err := h.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		switch {
		case err == service.ErrInvalidResetToken:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid or expired reset token",
				"code":  "INVALID_TOKEN",
			})
		case err == service.ErrResetTokenUsed:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Reset token has already been used",
				"code":  "TOKEN_USED",
			})
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
				"code":  "RESET_FAILED",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

// Me handles GET /auth/me — returns the authenticated user from the PASETO cookie.
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	token := c.Cookies("paseto_token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No authentication token",
			"code":  "MISSING_TOKEN",
		})
	}

	claims, err := h.pasetoService.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
			"code":  "INVALID_TOKEN",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":           claims.UserID,
			"role":         claims.Role,
			"territory_id": claims.TerritoryID,
		},
	})
}

// CookieConfig returns the standard cookie configuration for auth tokens.
func CookieConfig(token string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "paseto_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   int((1 * time.Hour).Seconds()),
	}
}
