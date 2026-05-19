package middleware

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"harmoni-api/internal/infrastructure/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp(t *testing.T, pasetoSvc *auth.PasetoService, publicRoutes []string) *fiber.App {
	t.Helper()

	app := fiber.New()

	if publicRoutes == nil {
		publicRoutes = DefaultPublicRoutes()
	}

	mw := NewAuthMiddleware(AuthMiddlewareConfig{
		PasetoService: pasetoSvc,
		PublicRoutes:  publicRoutes,
	})

	app.Use(mw)

	// Protected route
	app.Get("/api/protected", func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not set"})
		}
		claims := user.(*auth.Claims)
		return c.JSON(fiber.Map{
			"user_id":      claims.UserID,
			"role":         claims.Role,
			"territory_id": claims.TerritoryID,
		})
	})

	// Public route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Auth routes
	app.Post("/auth/register", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "registered"})
	})
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "logged in"})
	})
	app.Post("/auth/reset", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "reset requested"})
	})
	app.Post("/auth/reset/confirm", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "password reset"})
	})

	return app
}

func newTestPasetoService(t *testing.T) *auth.PasetoService {
	t.Helper()
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 1)
	}
	hexKey := hex.EncodeToString(key)

	svc, err := auth.NewPasetoService(hexKey)
	if err != nil {
		t.Fatalf("failed to create paseto service: %v", err)
	}
	return svc
}

func doRequest(app *fiber.App, method, path, token string) (*http.Response, error) {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.AddCookie(&http.Cookie{
			Name:  "paseto_token",
			Value: token,
		})
	}
	return app.Test(req)
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	defer resp.Body.Close()
	return string(body)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	svc := newTestPasetoService(t)
	app := setupTestApp(t, svc, nil)

	resp, err := doRequest(app, "GET", "/api/protected", "")
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body := readBody(t, resp)
	assert.Contains(t, body, "MISSING_TOKEN")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	svc := newTestPasetoService(t)
	app := setupTestApp(t, svc, nil)

	resp, err := doRequest(app, "GET", "/api/protected", "invalid-token-value")
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body := readBody(t, resp)
	assert.Contains(t, body, "INVALID_TOKEN")
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	svc := newTestPasetoService(t)
	app := setupTestApp(t, svc, nil)

	token, err := svc.GenerateToken("user-123", "rt_officer", "rt-01", -time.Hour)
	assert.NoError(t, err)

	resp, err := doRequest(app, "GET", "/api/protected", token)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body := readBody(t, resp)
	assert.Contains(t, body, "TOKEN_EXPIRED")
}

func TestAuthMiddleware_PublicRouteBypass(t *testing.T) {
	svc := newTestPasetoService(t)
	app := setupTestApp(t, svc, nil)

	resp, err := doRequest(app, "GET", "/health", "")
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := readBody(t, resp)
	assert.Contains(t, body, "ok")
}

func TestAuthMiddleware_AuthRouteBypass(t *testing.T) {
	svc := newTestPasetoService(t)
	app := setupTestApp(t, svc, nil)

	authRoutes := []string{"/auth/register", "/auth/login", "/auth/reset", "/auth/reset/confirm"}
	for _, route := range authRoutes {
		resp, err := doRequest(app, "POST", route, "")
		assert.NoError(t, err)
		assert.NotEqual(t, fiber.StatusUnauthorized, resp.StatusCode, "route %s should bypass auth", route)
	}
}

func TestAuthMiddleware_ValidTokenSetsUserClaims(t *testing.T) {
	svc := newTestPasetoService(t)
	app := setupTestApp(t, svc, nil)

	token, err := svc.GenerateToken("user-123", "rt_officer", "rt-01", time.Hour)
	assert.NoError(t, err)

	resp, err := doRequest(app, "GET", "/api/protected", token)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body := readBody(t, resp)
	assert.Contains(t, body, "user-123")
	assert.Contains(t, body, "rt_officer")
	assert.Contains(t, body, "rt-01")
}

func TestAuthMiddleware_CustomPublicRoutes(t *testing.T) {
	svc := newTestPasetoService(t)

	app := fiber.New()
	mw := NewAuthMiddleware(AuthMiddlewareConfig{
		PasetoService: svc,
		PublicRoutes:  []string{"/public", "/api/open"},
	})
	app.Use(mw)
	app.Get("/public", func(c *fiber.Ctx) error {
		return c.SendString("public")
	})
	app.Get("/api/open/data", func(c *fiber.Ctx) error {
		return c.SendString("open data")
	})

	resp, err := doRequest(app, "GET", "/public", "")
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	resp, err = doRequest(app, "GET", "/api/open/data", "")
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestDefaultPublicRoutes(t *testing.T) {
	routes := DefaultPublicRoutes()
	assert.Contains(t, routes, "/health")
	assert.Contains(t, routes, "/auth/register")
	assert.Contains(t, routes, "/auth/login")
	assert.Contains(t, routes, "/auth/reset")
	assert.Len(t, routes, 4)
}
