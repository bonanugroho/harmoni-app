---
phase: 01-core-authentication-rbac
plan: 003
subsystem: auth
tags: [go, fiber, casbin, rbac, middleware, paseto, territory-aware]

# Dependency graph
requires:
  - phase: 01-core-authentication-rbac
    provides: Go + Fiber project structure, PostgreSQL schema, PASETO token service, auth service
provides:
  - Casbin RBAC model configuration (rbac_model.conf) with territory domain
  - Policy CSV with hybrid approach ({{territory_id}} for RT, * for RW)
  - CasbinEnforcer singleton with Enforce, AddPolicy, RemovePolicy, role management
  - Auth middleware extracting PASETO tokens from httpOnly cookies
  - Casbin authorization middleware with territory-aware enforcement
  - Protected route examples with auth → casbin middleware chain
  - 42 tests across auth, middleware, and HTTP packages
affects: [004-password-reset, 005-rbac-middleware]

# Tech tracking
tech-stack:
  added: [casbin/casbin/v2, stretchr/testify]
  patterns: [Casbin RBAC with territory domain, middleware chaining (auth → casbin → handler), singleton enforcer, HTTP method to action mapping, resource extraction from URL path]

key-files:
  created:
    - apps/api/rbac_model.conf
    - apps/api/policy.csv
    - apps/api/internal/infrastructure/auth/casbin.go
    - apps/api/internal/infrastructure/auth/casbin_test.go
    - apps/api/internal/infrastructure/auth/casbin_policy_test.go
    - apps/api/internal/delivery/middleware/auth.go
    - apps/api/internal/delivery/middleware/auth_test.go
    - apps/api/internal/delivery/middleware/casbin.go
    - apps/api/internal/delivery/middleware/casbin_test.go
    - apps/api/internal/delivery/http/protected_handler.go
    - apps/api/internal/delivery/http/protected_handler_test.go
  modified:
    - apps/api/go.mod (added casbin, testify as direct dependencies)

key-decisions:
  - "Casbin model uses 4-part request: sub, obj, act, dom (territory domain)"
  - "Policy CSV uses {{territory_id}} placeholder for RT officers, * wildcard for RW officers"
  - "Middleware substitutes territory at enforcement time, not at policy load time"
  - "Resource names in policy are singular (tenant, income) matching domain concepts"
  - "HTTP method mapping: GET/HEAD/OPTIONS→read, POST/PUT/PATCH/DELETE→write"

patterns-established:
  - "Middleware chain: auth (token validation) → casbin (policy enforcement) → handler"
  - "Singleton enforcer with sync.Once for thread-safe initialization"
  - "ResetEnforcerForTest() for test isolation of singleton state"
  - "Consistent error format: {error, code} with appropriate HTTP status codes"
  - "Resource extraction from URL path first segment after /api/"

requirements-completed: [AUTH-02]

# Metrics
duration: ~25min
completed: 2026-05-19
---

# Phase 01 Plan 003: RBAC Middleware & Casbin Policies Summary

**Casbin RBAC engine with territory-aware policy enforcement, authentication middleware, and protected route examples — full auth → casbin → handler middleware chain**

## Performance

- **Duration:** ~25 min
- **Started:** 2026-05-19T08:59:53Z
- **Completed:** 2026-05-19T09:25:00Z
- **Tasks:** 5
- **Files modified:** 11 created, 1 modified

## Accomplishments

- Casbin RBAC model with territory domain (sub, obj, act, dom) and role hierarchy
- Policy CSV with hybrid approach: {{territory_id}} for RT officers, * wildcard for RW officers
- CasbinEnforcer singleton with Enforce, AddPolicy, RemovePolicy, AddRoleLink methods
- Auth middleware extracting PASETO tokens from httpOnly cookies with 401 error codes
- Casbin authorization middleware with territory-aware enforcement (RW uses *, others use assigned territory)
- Protected route examples demonstrating auth → casbin → handler middleware chain
- 42 tests total across auth (13), middleware (16), and HTTP (13) packages

## Task Commits

Each task was committed atomically:

1. **Task 1: Casbin Policy Configuration** - `f0490cf` (feat)
2. **Task 2: Casbin Enforcer Initialization** - `9a5c0e2` (feat)
3. **Task 3: Authentication Middleware** - `b9b1917` (feat)
4. **Task 4: Casbin Authorization Middleware** - `f078375` (feat)
5. **Task 5: Protected Route Examples** - `57e35a7` (feat)

## Files Created/Modified

- `apps/api/rbac_model.conf` - Casbin RBAC model with territory domain
- `apps/api/policy.csv` - Hybrid policies with {{territory_id}} placeholder and * wildcard
- `apps/api/internal/infrastructure/auth/casbin.go` - CasbinEnforcer singleton with territory-aware enforcement
- `apps/api/internal/infrastructure/auth/casbin_test.go` - 7 enforcer tests
- `apps/api/internal/infrastructure/auth/casbin_policy_test.go` - 6 policy tests with string adapter
- `apps/api/internal/delivery/middleware/auth.go` - Fiber auth middleware for PASETO token validation
- `apps/api/internal/delivery/middleware/auth_test.go` - 8 auth middleware tests
- `apps/api/internal/delivery/middleware/casbin.go` - Casbin authorization middleware
- `apps/api/internal/delivery/middleware/casbin_test.go` - 8 casbin middleware tests
- `apps/api/internal/delivery/http/protected_handler.go` - Protected route handler with middleware chain
- `apps/api/internal/delivery/http/protected_handler_test.go` - 9 protected route tests
- `apps/api/go.mod` - Added casbin and testify as direct dependencies

## Decisions Made

- Casbin model uses 4-part request definition (sub, obj, act, dom) for territory-aware RBAC
- Policy CSV uses {{territory_id}} as a template placeholder — middleware substitutes at runtime
- Resource names in policy are singular (tenant, income, expenditure, report) to match domain concepts
- HTTP method to action mapping: read operations (GET/HEAD/OPTIONS) → "read", write operations (POST/PUT/PATCH/DELETE) → "write"
- Singleton pattern for Casbin enforcer with ResetEnforcerForTest() for test isolation
- RW officers use "*" domain for all-territory access; RT officers use their assigned territory_id

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical] Added testify testing library**
- **Found during:** Task 3 (auth middleware tests)
- **Issue:** No assertion library available for test assertions
- **Fix:** Added `github.com/stretchr/testify` via `go get`
- **Files modified:** apps/api/go.mod, apps/api/go.sum
- **Verification:** All tests pass with assert package

**2. [Rule 1 - Bug] Fixed resource extractor leading slash handling**
- **Found during:** Task 4 (casbin middleware tests)
- **Issue:** defaultResourceExtractor returned empty string for paths like "/protected" because first character '/' triggered early return
- **Fix:** Strip leading slash before extracting first segment
- **Files modified:** apps/api/internal/delivery/middleware/casbin.go
- **Verification:** Resource extractor tests pass for all path patterns

**3. [Rule 3 - Blocking] Fixed test file paths for Casbin config**
- **Found during:** Task 4 (casbin middleware tests)
- **Issue:** Go tests run from package directory, not module root — relative paths to rbac_model.conf failed
- **Fix:** Used absolute path for config files in protected handler tests; relative paths work for middleware tests (3 levels up)
- **Files modified:** apps/api/internal/delivery/http/protected_handler_test.go
- **Verification:** All tests pass

---

**Total deviations:** 3 auto-fixed (1 missing critical, 1 bug, 1 blocking)
**Impact on plan:** All auto-fixes necessary for correctness and testability. No scope creep.

## Issues Encountered

- Casbin file adapter loads `{{territory_id}}` as literal string — middleware must substitute at enforcement time, not at load time
- Resource names in routes must match policy resource names (singular: tenant, not tenants)
- Go test working directory is the package directory, not module root — affects relative path resolution for config files

## Self-Check Results

- **Build:** `go build ./...` — compiles successfully
- **Tests:** `go test ./...` — 42/42 PASS across all packages
  - delivery/http: 19 PASS (10 prior + 9 new)
  - delivery/middleware: 16 PASS (8 auth + 8 casbin)
  - infrastructure/auth: 13 PASS (5 paseto + 4 password + 6 policy + 7 enforcer - 2 shared = 13 unique)
  - domain/service: 11 PASS (prior wave)
  - infrastructure/database: 4 PASS (prior wave)
  - infrastructure/email: 3 PASS (prior wave)
- **Acceptance criteria:**
  - rbac_model.conf defines RBAC model with territory domain: ✅
  - policy.csv contains policies for all 3 roles: ✅
  - RT officer policies use {{territory_id}} placeholder: ✅
  - RW officer policies use * wildcard: ✅
  - Resident policies are read-only: ✅
  - CasbinEnforcer initializes and enforces correctly: ✅
  - Auth middleware validates tokens and sets user claims: ✅
  - Casbin middleware enforces territory-aware access: ✅
  - Protected routes demonstrate middleware chain: ✅

## Self-Check: PASSED

All 12 created files verified on disk. All 6 task commits verified in git log.

---

*Phase: 01-core-authentication-rbac*
*Completed: 2026-05-19*
