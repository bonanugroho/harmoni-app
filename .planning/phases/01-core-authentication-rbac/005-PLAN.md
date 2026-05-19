---
wave: 5
depends_on:
  - 002
  - 003
  - 004
files_modified:
  - apps/api/internal/delivery/http/auth_handler_test.go
  - apps/api/internal/delivery/middleware/auth_test.go
  - apps/api/internal/delivery/middleware/casbin_test.go
  - apps/api/internal/domain/service/auth_service_test.go
  - apps/api/internal/infrastructure/auth/paseto_test.go
  - apps/api/internal/infrastructure/auth/password_test.go
  - apps/api/internal/infrastructure/repository/user_repository_test.go
autonomous: true
requirements:
  - AUTH-01
  - AUTH-02
---

# Plan 5: Integration Tests & Verification

## Objective
Create comprehensive test suite covering all authentication flows and RBAC enforcement.

## Tasks

### Task 1: Unit Tests - Auth Service
<read_first>
- apps/api/internal/domain/service/auth_service_test.go (create)
- apps/api/internal/domain/service/auth_service.go
- apps/api/internal/domain/repository/user_repository.go (mock)
</read_first>

<action>
Write unit tests for auth service:
- TestRegister_Success: valid data creates user with hashed password
- TestRegister_DuplicateEmail: returns error for existing email
- TestRegister_WeakPassword: returns error for weak password
- TestLogin_Success: valid credentials return user + token
- TestLogin_InvalidEmail: returns error for non-existent email
- TestLogin_WrongPassword: returns error for incorrect password
- TestResetPasswordRequest_Success: creates reset token
- TestResetPassword_Success: updates password, marks token used
- TestResetPassword_InvalidToken: returns error for invalid token
- TestResetPassword_ExpiredToken: returns error for expired token
- Mock UserRepository interface for isolation
</action>

<acceptance_criteria>
- All 10 test cases pass
- Mock repository returns expected responses
- Password hashing verified (hash != plaintext)
- Token generation verified (valid PASETO token)
- Reset token expiry verified (1 hour)
- `go test ./internal/domain/service -v -cover` shows >= 80% coverage
</acceptance_criteria>

---

### Task 2: Unit Tests - PASETO & Password
<read_first>
- apps/api/internal/infrastructure/auth/paseto_test.go (create)
- apps/api/internal/infrastructure/auth/password_test.go (create)
- apps/api/internal/infrastructure/auth/paseto.go
- apps/api/internal/infrastructure/auth/password.go
</read_first>

<action>
Write unit tests for PASETO and password services:
- TestGenerateToken_Success: creates valid token with correct claims
- TestValidateToken_Success: validates token and returns claims
- TestValidateToken_Expired: returns error for expired token
- TestValidateToken_Invalid: returns error for tampered token
- TestHashPassword_Success: creates bcrypt hash
- TestComparePassword_Success: verifies correct password
- TestComparePassword_Wrong: returns error for incorrect password
- TestValidatePassword_Weak: rejects weak passwords
- TestValidatePassword_Strong: accepts strong passwords
</action>

<acceptance_criteria>
- All 9 test cases pass
- Token claims match input values
- Expired token detection works correctly
- Password complexity rules enforced
- `go test ./internal/infrastructure/auth -v -cover` shows >= 85% coverage
</acceptance_criteria>

---

### Task 3: Unit Tests - Repository
<read_first>
- apps/api/internal/infrastructure/repository/user_repository_test.go (create)
- apps/api/internal/infrastructure/repository/user_repository.go
- apps/api/migrations/*.sql
</read_first>

<action>
Write integration tests for user repository:
- Setup: create test database, run migrations
- TestCreate_Success: inserts user, returns created user
- TestFindByEmail_Success: finds user by email
- TestFindByEmail_NotFound: returns sql.ErrNoRows
- TestListByTerritory_Success: filters users by territory
- TestUpdatePassword_Success: updates password hash
- Teardown: drop test database
- Use testcontainers-go for isolated PostgreSQL instance
</action>

<acceptance_criteria>
- All 5 test cases pass
- Test database created and migrated before each test
- Test database dropped after each test (isolation)
- Territory filtering returns correct subset of users
- Password update persists to database
- `go test ./internal/infrastructure/repository -v -tags=integration` passes
</acceptance_criteria>

---

### Task 4: Integration Tests - Auth Handlers
<read_first>
- apps/api/internal/delivery/http/auth_handler_test.go (create)
- apps/api/internal/delivery/http/auth_handler.go
- apps/api/internal/domain/service/auth_service.go (mock)
</read_first>

<action>
Write integration tests for auth HTTP handlers:
- Setup: Fiber app with mocked auth service
- TestRegisterHandler_Success: 201 + user object
- TestRegisterHandler_DuplicateEmail: 409 + error
- TestRegisterHandler_WeakPassword: 400 + error
- TestLoginHandler_Success: 200 + httpOnly cookie
- TestLoginHandler_InvalidCredentials: 401 + error
- TestResetHandler_Success: 200 + email sent
- TestResetConfirmHandler_Success: 200 + password updated
- TestResetConfirmHandler_InvalidToken: 400 + error
- Verify cookie attributes: HttpOnly, Secure, SameSite, Path, MaxAge
</action>

<acceptance_criteria>
- All 8 test cases pass
- Response status codes match expected values
- Response body contains correct error messages
- httpOnly cookie set with correct attributes
- Mock service called with correct arguments
- `go test ./internal/delivery/http -v -run TestAuthHandler` passes
</acceptance_criteria>

---

### Task 5: Integration Tests - Middleware
<read_first>
- apps/api/internal/delivery/middleware/auth_test.go (create)
- apps/api/internal/delivery/middleware/casbin_test.go (create)
- apps/api/internal/delivery/middleware/auth.go
- apps/api/internal/delivery/middleware/casbin.go
</read_first>

<action>
Write integration tests for auth and casbin middleware:
- TestAuthMiddleware_MissingToken: 401 + error
- TestAuthMiddleware_InvalidToken: 401 + error
- TestAuthMiddleware_ExpiredToken: 401 + error
- TestAuthMiddleware_ValidToken: sets user claims, calls next
- TestCasbinMiddleware_RTOwnerOwnTerritory: 200
- TestCasbinMiddleware_RTOwnerOtherTerritory: 403
- TestCasbinMiddleware_RWOfficerAnyTerritory: 200
- TestCasbinMiddleware_ResidentWrite: 403
- Test public routes bypass authentication
</action>

<acceptance_criteria>
- All 9 test cases pass
- Missing/invalid/expired tokens return 401
- Valid token sets user claims in context
- Territory isolation enforced for RT officers
- RW officers have access to all territories
- Residents cannot perform write operations
- Public routes (/health, /auth/*) bypass auth
- `go test ./internal/delivery/middleware -v` passes
</acceptance_criteria>

---

### Task 6: E2E Test Script
<read_first>
- apps/api/tests/e2e/auth_flow.sh (create)
- apps/api/.env.example
</read_first>

<action>
Create end-to-end test script:
- Start test server with test database
- Register new user
- Login with credentials
- Access protected route with token
- Request password reset
- Confirm password reset
- Login with new password
- Test RBAC: RT officer access, RW officer access, resident restrictions
- Clean up test data
- Exit 0 if all tests pass, exit 1 if any fail
</action>

<acceptance_criteria>
- apps/api/tests/e2e/auth_flow.sh is executable
- Script runs all auth flows end-to-end
- Script exits 0 on success, 1 on failure
- Test data cleaned up after execution
- Script outputs pass/fail for each step
- `bash apps/api/tests/e2e/auth_flow.sh` passes
</acceptance_criteria>

---

## Verification

1. **Unit Tests:**
   ```bash
   cd apps/api && go test ./... -v -cover
   ```
   Expected: All tests pass, >= 80% coverage

2. **Integration Tests:**
   ```bash
   cd apps/api && go test ./... -v -tags=integration
   ```
   Expected: All integration tests pass

3. **E2E Test:**
   ```bash
   cd apps/api && bash tests/e2e/auth_flow.sh
   ```
   Expected: All steps pass, exit 0

4. **Coverage Report:**
   ```bash
   cd apps/api && go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out -o coverage.html
   ```
   Expected: coverage.html generated, >= 80% line coverage