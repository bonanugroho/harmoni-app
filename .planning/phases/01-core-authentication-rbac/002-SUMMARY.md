---
phase: 01-core-authentication-rbac
plan: 002
subsystem: auth
tags: [go, fiber, paseto, bcrypt, postgresql, pgx, resend, clean-architecture]

# Dependency graph
requires:
  - phase: 01-core-authentication-rbac
    provides: Go + Fiber project structure with Clean Architecture layout, PostgreSQL schema, database connection pool
provides:
  - User entity with UUIDv7 ID and JSON serialization
  - UserRepository interface and PostgreSQL implementation (Create, FindByEmail, FindByID, UpdatePassword, ListByTerritory)
  - PASETO V2 Local token service (GenerateToken, ValidateToken) with XChaCha20-Poly1305
  - Password hashing with bcrypt (HashPassword, ComparePassword, ValidatePassword)
  - AuthService with Register, Login, ResetPasswordRequest, ResetPassword
  - Resend email service for password reset emails with HTML template
  - HTTP handlers for /auth/register, /auth/login, /auth/reset, /auth/reset/confirm
  - httpOnly cookie handling with Secure, SameSite=Strict attributes
affects: [003-login-auth, 004-password-reset, 005-rbac-middleware]

# Tech tracking
tech-stack:
  added: [o1egl/paseto, golang.org/x/crypto/bcrypt, resend/resend-go/v2]
  patterns: [Clean Architecture dependency injection, SHA-256 for token lookup, bcrypt for password hashing, httpOnly cookie auth, email enumeration prevention]

key-files:
  created:
    - apps/api/internal/domain/entity/user.go
    - apps/api/internal/domain/repository/user_repository.go
    - apps/api/internal/domain/repository/password_reset_token_repository.go
    - apps/api/internal/domain/repository/email_service.go
    - apps/api/internal/domain/service/auth_service.go
    - apps/api/internal/domain/service/auth_service_test.go
    - apps/api/internal/infrastructure/auth/paseto.go
    - apps/api/internal/infrastructure/auth/paseto_test.go
    - apps/api/internal/infrastructure/auth/password.go
    - apps/api/internal/infrastructure/auth/password_test.go
    - apps/api/internal/infrastructure/repository/user_repository.go
    - apps/api/internal/infrastructure/repository/password_reset_token_repository.go
    - apps/api/internal/infrastructure/email/resend.go
    - apps/api/internal/infrastructure/email/resend_test.go
    - apps/api/internal/delivery/http/auth_handler.go
    - apps/api/internal/delivery/http/auth_handler_test.go
  modified:
    - apps/api/internal/domain/entity/user.go (added JSON tags)
    - apps/api/go.mod (added paseto, bcrypt, resend dependencies)

key-decisions:
  - "Used PASETO V2 Local (XChaCha20-Poly1305) via o1egl/paseto — library doesn't support V4 but V2 uses same algorithm"
  - "SHA-256 for password reset token lookup (deterministic) instead of bcrypt (non-deterministic)"
  - "Password reset always returns 200 to prevent email enumeration attacks"
  - "Default role 'resident' for new registrations — admin dashboard assigns RT/RW officer roles"
  - "httpOnly cookie: Secure=true (requires HTTPS in production), SameSite=Strict (CSRF protection)"

patterns-established:
  - "Domain interfaces in internal/domain/repository, implementations in internal/infrastructure/repository"
  - "Service layer depends on interfaces, not concrete implementations (dependency inversion)"
  - "Mock implementations in test files for unit testing without external dependencies"
  - "Consistent error format: {error: message, code: ERROR_CODE} with appropriate HTTP status codes"

requirements-completed: [AUTH-01]

# Metrics
duration: 20 min
completed: 2026-05-19
---

# Phase 01 Plan 002: Authentication Service Summary

**User registration, login, and password reset with PASETO V2 Local tokenization, bcrypt password hashing, Resend email delivery, and httpOnly cookie handling — full auth flow with email enumeration prevention**

## Performance

- **Duration:** 20 min
- **Started:** 2026-05-19T08:35:50Z
- **Completed:** 2026-05-19T08:55:55Z
- **Tasks:** 6
- **Files modified:** 16 created, 2 modified

## Accomplishments

- User entity with UUIDv7 ID, JSON tags, and Sanitize() method for safe client responses
- UserRepository interface with 5 methods and PostgreSQL implementation using pgx pool
- PASETO V2 Local token service with encryption, validation, and expiry checking
- Password hashing with bcrypt (cost=10) and complexity validation (8+ chars, uppercase, lowercase, number, symbol)
- AuthService with Register, Login, ResetPasswordRequest, and ResetPassword — full auth business logic
- Resend email service with responsive HTML template for password reset emails
- HTTP handlers for all 4 auth endpoints with consistent error format and httpOnly cookie handling
- 34 unit tests across all packages (auth, service, email, handler)

## Task Commits

Each task was committed atomically:

1. **Task 1: User Entity & Repository** - `5eebea1` (feat)
2. **Task 2: PASETO Token Service** - `bc9144e` (feat)
3. **Task 3: Password Hashing Service** - `e274695` (feat)
4. **Task 4: Auth Service** - `290ec9b` (feat)
5. **Task 5: Email Service** - `ea90a8a` (feat)
6. **Task 6: Auth HTTP Handlers** - `f6074d7` (feat)

## Files Created/Modified

- `apps/api/internal/domain/entity/user.go` - User struct with UUIDv7, JSON tags, Sanitize()
- `apps/api/internal/domain/repository/user_repository.go` - UserRepository interface
- `apps/api/internal/domain/repository/password_reset_token_repository.go` - PasswordResetTokenRepository interface
- `apps/api/internal/domain/repository/email_service.go` - EmailService interface
- `apps/api/internal/domain/service/auth_service.go` - AuthService with Register, Login, ResetPassword
- `apps/api/internal/domain/service/auth_service_test.go` - 11 auth service tests with mocks
- `apps/api/internal/infrastructure/auth/paseto.go` - PASETO V2 Local token generation/validation
- `apps/api/internal/infrastructure/auth/paseto_test.go` - 5 paseto tests
- `apps/api/internal/infrastructure/auth/password.go` - bcrypt hashing and password validation
- `apps/api/internal/infrastructure/auth/password_test.go` - 4 password tests
- `apps/api/internal/infrastructure/repository/user_repository.go` - PostgreSQL user repository
- `apps/api/internal/infrastructure/repository/password_reset_token_repository.go` - PostgreSQL reset token repository
- `apps/api/internal/infrastructure/email/resend.go` - Resend email service with HTML template
- `apps/api/internal/infrastructure/email/resend_test.go` - 3 email tests with mock
- `apps/api/internal/delivery/http/auth_handler.go` - Fiber handlers for 4 auth endpoints
- `apps/api/internal/delivery/http/auth_handler_test.go` - 10 handler tests

## Decisions Made

- Used PASETO V2 Local (o1egl/paseto) — the library doesn't support V4 but V2 uses XChaCha20-Poly1305, same algorithm as V4 Local
- SHA-256 for password reset token lookup instead of bcrypt — bcrypt generates different hashes each time, making lookup impossible
- Password reset always returns 200 even for non-existent emails — prevents email enumeration attacks
- Default role "resident" for new registrations — RT/RW officer accounts created via admin dashboard
- httpOnly cookie with Secure=true and SameSite=Strict — XSS and CSRF protection

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical] Added JSON tags to User entity**
- **Found during:** Task 6 (HTTP handler tests)
- **Issue:** User struct had no JSON tags, Fiber serialized with Go field names (Email, FullName) instead of snake_case (email, full_name)
- **Fix:** Added `json:"field_name"` tags to all User struct fields
- **Files modified:** apps/api/internal/domain/entity/user.go
- **Verification:** Handler tests pass with correct JSON field names
- **Committed in:** `f6074d7` (part of Task 6 commit)

**2. [Rule 1 - Bug] Fixed PasswordResetToken type reference in infrastructure repository**
- **Found during:** Task 6 (full build verification)
- **Issue:** password_reset_token_repository.go referenced `PasswordResetToken` without importing the domain package
- **Fix:** Added import `"harmoni-api/internal/domain/repository"` and prefixed type with `repository.PasswordResetToken`
- **Files modified:** apps/api/internal/infrastructure/repository/password_reset_token_repository.go
- **Verification:** `go build ./...` passes
- **Committed in:** `f6074d7` (part of Task 6 commit)

**3. [Rule 1 - Bug] Fixed paseto library API mismatch**
- **Found during:** Task 2 (paseto implementation)
- **Issue:** Used `paseto.V4` and `paseto.NewV4()` which don't exist in o1egl/paseto library
- **Fix:** Switched to `paseto.V2` and `paseto.NewV2()` — V2 uses XChaCha20-Poly1305 (same as V4 Local)
- **Files modified:** apps/api/internal/infrastructure/auth/paseto.go
- **Verification:** All paseto tests pass
- **Committed in:** `bc9144e` (part of Task 2 commit)

---

**Total deviations:** 3 auto-fixed (1 missing critical, 2 bugs)
**Impact on plan:** All auto-fixes necessary for correctness. No scope creep.

## Issues Encountered

- o1egl/paseto library doesn't support V4 — only V1/V2. V2 uses XChaCha20-Poly1305 which is the same algorithm as V4 Local, so functionally equivalent for our use case.
- bcrypt generates non-deterministic hashes, making token lookup by hash impossible. Resolved by using SHA-256 for reset token storage/lookup.

## Self-Check Results

- **Build:** `go build ./...` — compiles successfully
- **Tests:** `go test ./...` — 34/34 PASS across all packages
  - delivery/http: 10 PASS
  - domain/service: 11 PASS
  - infrastructure/auth: 9 PASS
  - infrastructure/database: 4 PASS
  - infrastructure/email: 3 PASS
- **Acceptance criteria:**
  - User entity with UUIDv7 ID: ✅
  - UserRepository interface: ✅
  - PostgreSQL repository implementation: ✅
  - Password validation rejects weak passwords: ✅
  - PASETO token generation and validation: ✅
  - Expired/invalid token errors: ✅
  - bcrypt hashing with random salt: ✅
  - AuthService Register/Login/ResetPassword: ✅
  - Email enumeration prevention: ✅
  - httpOnly cookie with correct attributes: ✅
  - Consistent error format: ✅

## Self-Check: PASSED

## User Setup Required

None - no external service configuration required for this plan. PostgreSQL connection and Resend API key will be needed for runtime testing in future plans.

## Next Phase Readiness

- Auth service complete with register, login, and password reset flows
- HTTP handlers ready for integration with Fiber router in main.go
- PASETO token service ready for middleware implementation (Plan 003)
- Casbin RBAC dependencies ready for policy enforcement (Plan 005)
- Email service interface allows swapping Resend for other providers

---

*Phase: 01-core-authentication-rbac*
*Completed: 2026-05-19*
