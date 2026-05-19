# Phase 1: Core Authentication & RBAC — Research

**Date:** 2026-05-19
**Phase:** 1
**Goal:** Implement secure user registration, login, password reset, and role-based access control for Residents, RT Officers, and RW Officers.

---

## 1. PASETO V4 Local — Go Implementation

### Recommended Library
- **`oasislmf/paseto`** or **`aidarkhanov/paseto`** — both support V4 Local (symmetric encryption)
- **Alternative:** `github.com/o1egl/paseto` — widely used, well-documented

### Key Management
- PASETO V4 Local uses **XChaCha20-Poly1305** for encryption
- **Key size:** 32 bytes (256-bit)
- **Storage:** Environment variable `PASETO_SECRET_KEY` (hex-encoded)
- **Rotation:** Generate new key, re-encrypt active sessions, retire old key after TTL

### httpOnly Cookie Handling
```go
http.SetCookie(w, &http.Cookie{
    Name:     "paseto_token",
    Value:    token,
    HttpOnly: true,        // Prevents XSS
    Secure:   true,        // HTTPS only
    SameSite: http.SameSiteStrictMode,
    Path:     "/",
    MaxAge:   3600,        // 1 hour
})
```

### Token Structure
- **Payload:** `{"user_id": "uuid", "role": "resident", "territory_id": "rt-01", "exp": "2026-05-19T12:00:00Z"}`
- **Validation:** Check expiration, verify signature, extract claims
- **Refresh:** Optional refresh token in separate httpOnly cookie with longer TTL

---

## 2. Casbin RBAC — Go Integration

### Recommended Library
- **`github.com/casbin/casbin/v2`** — official Go SDK
- **Adapter:** `github.com/casbin/gorm-adapter/v3` for PostgreSQL persistence

### Policy Structure (Hybrid Approach)
```csv
p, resident, tenant, read, {{territory_id}}
p, resident, income, read, {{territory_id}}

p, rt_officer, tenant, read, {{territory_id}}
p, rt_officer, tenant, write, {{territory_id}}
p, rt_officer, income, read, {{territory_id}}
p, rt_officer, income, write, {{territory_id}}

p, rw_officer, tenant, read, *
p, rw_officer, tenant, write, *
p, rw_officer, income, read, *
p, rw_officer, income, write, *
```

### Middleware Implementation
```go
func CasbinMiddleware(e *casbin.Enforcer) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*UserClaims)
        obj := c.Params("resource")
        act := c.Method() // map HTTP method to action
        
        // Substitute territory_id for non-RW officers
        domain := user.TerritoryID
        if user.Role == "rw_officer" {
            domain = "*"
        }
        
        ok, err := e.Enforce(user.Role, obj, act, domain)
        if err != nil || !ok {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
        }
        return c.Next()
    }
}
```

### Role Hierarchy
```csv
g, rw_officer, rt_officer  # RW inherits RT permissions
g, rt_officer, resident    # RT inherits Resident permissions (optional)
```

---

## 3. PostgreSQL Schema Design

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,  -- UUIDv7 (time-ordered, sortable)
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('resident', 'rt_officer', 'rw_officer')),
    territory_id VARCHAR(50) NOT NULL REFERENCES territories(id),
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_territory ON users(territory_id);
```

### Territories Table
```sql
CREATE TABLE territories (
    id VARCHAR(50) PRIMARY KEY,  -- e.g., 'rt-01', 'rt-02', 'rw-01'
    name VARCHAR(100) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('rt', 'rw')),
    parent_id VARCHAR(50) REFERENCES territories(id),  -- RW is parent of RTs
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Sessions Table (for token revocation)
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY,  -- UUIDv7
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- UUIDv7
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
```

### Password Reset Tokens Table
```sql
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY,  -- UUIDv7
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- UUIDv7
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reset_tokens_user ON password_reset_tokens(user_id);
```

### UUIDv7 Generation
- **PostgreSQL 18+:** Native `uuidv7()` function available — **RECOMMENDED**
- **Go Library (fallback for PG < 18):** `github.com/google/uuid` v1.6.0+ with `uuid.NewUUIDv7()`
- **Best Practice:** Use PostgreSQL 18's `uuidv7()` as DEFAULT for all UUID columns
- **Migration syntax:** `id UUID PRIMARY KEY DEFAULT uuidv7()`
- **Benefits:** Time-ordered, sortable, better index performance, chronological insertion order, no application-layer generation needed

---

## 4. golang-migrate CLI

### Installation
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Migration File Structure
```
migrations/
  001_create_territories_table.up.sql
  001_create_territories_table.down.sql
  002_create_users_table.up.sql
  002_create_users_table.down.sql
  003_create_sessions_table.up.sql
  003_create_sessions_table.down.sql
  004_create_password_reset_tokens_table.up.sql
  004_create_password_reset_tokens_table.down.sql
  005_seed_territories.up.sql
  005_seed_territories.down.sql
```

### Commands
```bash
# Apply all migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/harmoni?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgres://..." down 1

# Create new migration
migrate create -ext sql -dir migrations -seq create_users_table
```

---

## 5. Email Service Options

### Recommended: **Resend** (resend.com)
- **Free tier:** 3,000 emails/month
- **Go SDK:** `github.com/resend/resend-go/v2`
- **Features:** Transactional emails, webhooks, analytics
- **Setup:** API key, verified domain

### Alternative: **SendGrid**
- **Free tier:** 100 emails/day
- **Go SDK:** `github.com/sendgrid/sendgrid-go`
- **Features:** Templates, tracking, analytics

### Alternative: **AWS SES**
- **Free tier:** 62,000 emails/month (from EC2)
- **Go SDK:** `github.com/aws/aws-sdk-go-v2/service/ses`
- **Features:** High deliverability, cost-effective at scale

### Password Reset Email Template
```html
<h2>Password Reset Request</h2>
<p>Click the link below to reset your password:</p>
<a href="{{.ResetURL}}">Reset Password</a>
<p>This link expires in {{.ExpiryHours}} hours.</p>
```

---

## 6. Password Hashing — bcrypt vs argon2

### Recommended: **bcrypt** (simpler, widely adopted)
```go
import "golang.org/x/crypto/bcrypt"

// Hash
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Verify
err := bcrypt.CompareHashAndPassword(hash, []byte(password))
```

### Alternative: **argon2** (more secure, newer)
```go
import "golang.org/x/crypto/argon2"

// Hash
hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
```

### Decision Factors
- **bcrypt:** Battle-tested, built into many frameworks, sufficient for MVP
- **argon2:** Winner of Password Hashing Competition, memory-hard, better against GPU attacks
- **Recommendation:** Start with **bcrypt** (cost=12), migrate to argon2 later if needed

---

## 7. Clean Architecture Layout — Go + Fiber

### Project Structure
```
apps/api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/
│   │   ├── entity/              # User, Territory, Session
│   │   ├── repository/          # Interfaces
│   │   └── service/             # Business logic
│   ├── infrastructure/
│   │   ├── database/            # PostgreSQL connection, migrations
│   │   ├── repository/          # Implementations
│   │   ├── auth/                # PASETO, Casbin middleware
│   │   └── email/               # Email service client
│   ├── delivery/
│   │   ├── http/                # Fiber handlers
│   │   └── middleware/          # Auth, logging, error handling
│   └── config/                  # Environment variables
├── migrations/                  # golang-migrate files
├── pkg/                         # Shared utilities
├── go.mod
└── .env.example
```

### Dependency Injection
- Use **wire** (`github.com/google/wire`) for compile-time DI
- Or manual DI in `main.go` for simplicity

---

## 8. httpOnly Cookie Security Best Practices

### Cookie Configuration
```go
http.SetCookie(w, &http.Cookie{
    Name:     "paseto_token",
    Value:    token,
    HttpOnly: true,        // Prevents JavaScript access (XSS protection)
    Secure:   true,        // HTTPS only (set false in dev)
    SameSite: http.SameSiteStrictMode,  // CSRF protection
    Path:     "/",
    MaxAge:   3600,        // 1 hour
    Domain:   "",          // Current domain only
})
```

### CSRF Protection
- **SameSite=Strict** prevents cross-site requests
- **Double Submit Cookie** pattern for additional CSRF protection
- **Custom header** required for API requests (e.g., `X-CSRF-Token`)

### Token Refresh
- **Access token:** Short-lived (1 hour), httpOnly cookie
- **Refresh token:** Long-lived (7 days), separate httpOnly cookie
- **Rotation:** New refresh token issued on each refresh

---

## 9. Territory-Aware Data Isolation Patterns

### Application-Level Filtering
```go
// Repository method
func (r *UserRepository) FindByTerritory(ctx context.Context, territoryID string) ([]*User, error) {
    query := "SELECT * FROM users WHERE territory_id = $1"
    // Execute query
}
```

### Database-Level Enforcement (Optional Safety Net)
```sql
-- Row Level Security (RLS)
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

CREATE POLICY territory_isolation ON users
    USING (territory_id = current_setting('app.current_territory'));
```

### Casbin Policy Enforcement
- **RT Officer:** `{{territory_id}}` replaced with user's territory
- **RW Officer:** `*` wildcard grants access to all territories
- **Resident:** Read-only access to own territory data

---

## 10. Testing Strategies for Auth Flows

### Unit Tests
- **Password hashing:** Verify hash generation and comparison
- **Token generation:** Verify PASETO token creation and validation
- **Casbin policies:** Verify permission checks

### Integration Tests
- **Registration:** POST /auth/register → 201 + user created
- **Login:** POST /auth/login → 200 + httpOnly cookie set
- **Password Reset:** POST /auth/reset → 200 + email sent
- **Token Validation:** GET /protected → 401 without token, 200 with valid token
- **RBAC:** Verify role-based access control with different users

### Test Setup
```go
// Test database
func setupTestDB() (*sql.DB, error) {
    // Create test database, run migrations
}

// Test client with cookie jar
func newTestClient() *http.Client {
    return &http.Client{
        Jar: cookiejar.New(nil),
    }
}
```

### Test Commands
```bash
# Run all tests
go test ./... -v

# Run integration tests
go test ./internal/delivery/http -v -tags=integration

# Coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Summary

Phase 1 is well-defined with clear decisions captured in CONTEXT.md. The research confirms:

1. **PASETO V4 Local** is mature in Go with good library support
2. **Casbin** provides flexible RBAC with territory filtering via wildcards
3. **PostgreSQL schema** is straightforward with clear relationships
4. **golang-migrate CLI** is the standard for Go database migrations
5. **Email services** have multiple viable options (Resend recommended)
6. **bcrypt** is sufficient for password hashing at MVP stage
7. **Clean Architecture** layout is well-established for Go + Fiber
8. **httpOnly cookies** provide strong XSS protection with proper configuration
9. **Territory isolation** can be enforced at application and database levels
10. **Testing strategies** cover unit, integration, and E2E flows
11. **UUIDv7** is natively supported in PostgreSQL 18 via `uuidv7()` function — use `DEFAULT uuidv7()` for all UUID columns, no application-layer generation needed

**Next Step:** Proceed to planning with `/gsd-plan-phase 1` to create executable PLAN.md files.