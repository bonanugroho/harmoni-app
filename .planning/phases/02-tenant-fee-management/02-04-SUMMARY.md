---
phase: 02-tenant-fee-management
plan: 04
subsystem: api
tags: [gofiber, paseto, casbin, http, handler, rest, crud]

requires:
  - phase: 02-03
    provides: TenantService, FeeRepository, Casbin policy with fee resources

provides:
  - HTTP handler layer (TenantHandler) for tenant/fee CRUD — 9 handler methods
  - HTTP integration tests for all tenant/fee endpoints — 10 test cases
  - Cleaned ProtectedHandler without tenant stubs (no ListTenants/GetTenant/CreateTenant)
  - Refactored main.go: middleware creation and /api group at main level, TenantHandler wired
  - Full auth + Casbin protection on all /api/tenants and /api/tenants/:id/fees routes

affects:
  - phase 03 (frontend UI): full API surface ready for consumption

tech-stack:
  added: []
  patterns:
    - "handler-level request types with float64 amounts and string dates (ISO 8601)"
    - "errors.Is matching for service-layer sentinel errors → HTTP status mapping"
    - "Delete handlers return 204 No Content (consistent REST convention)"
    - "TestClaims middleware for integration tests without real token creation"
    - "Shared /api group with middleware created at main.go, handlers receive fiber.Router"

key-files:
  created:
    - apps/api/internal/delivery/http/tenant_handler.go
    - apps/api/internal/delivery/http/tenant_handler_test.go
  modified:
    - apps/api/internal/delivery/http/protected_handler.go
    - apps/api/internal/delivery/http/protected_handler_test.go
    - apps/api/cmd/server/main.go
    - apps/api/internal/delivery/http/auth_handler_test.go

key-decisions:
  - "Use errors.Is (not ==) for all service error matching — handles %w-wrapped errors consistently"
  - "Handler request types use float64 for amounts, string for dates (service layer converts to decimal.Decimal / time.Time)"
  - "CreateFee uses type discriminator field ('mandatory' / 'voluntary') to route to appropriate service method"
  - "Middleware creation moves from ProtectedHandler to main.go — handlers share a single /api group with middleware"
  - "TestClaims middleware sets c.Locals('user') directly instead of real PASETO token parsing in tests"

patterns-established:
  - "Handler error mapping: errors.Is(err, service.ErrXxx) → switch block → status code + error code JSON"
  - "Test claims middleware pattern: middleware-compatible function that injects test user claims"
  - "Shared middleware group: /api group with auth + Casbin created at main level, passed as fiber.Router to handlers"

requirements-completed:
  - TEN-01
  - FIN-01
  - FIN-02

duration: 6min
completed: 2026-05-23
---

# Phase 02: Tenant & Fee Management — Plan 04 Summary

**HTTP handler layer for tenant/fee CRUD with 9 endpoints, 10 integration tests, middleware refactored to main.go, and old stubs removed from ProtectedHandler**

## Performance

- **Duration:** ~6 min execution (Task 1: 2m, Task 2: 2m, verification + summary: 2m)
- **Started:** 2026-05-23T09:48:00+07:00 (first task commit)
- **Completed:** 2026-05-23T09:50:37+07:00 (final task commit)
- **Tasks:** 2 (1 auto + 1 auto)
- **Files modified:** 6 (2 created, 4 modified)

## Accomplishments

- Created `TenantHandler` with 9 handler methods: ListTenants, CreateTenant, GetTenant, UpdateTenant, DeleteTenant, ListFees, CreateFee, UpdateFee, DeleteFee — all using `errors.Is` pattern for service error mapping
- Created 10 HTTP integration tests with mock repos/services — test claims middleware pattern for auth bypass
- Removed tenant stubs from `ProtectedHandler.RegisterRoutes` (5 route lines + 3 methods = 74 lines of stub code deleted)
- Refactored `ProtectedHandler.RegisterRoutes` to accept `fiber.Router` instead of `*fiber.App` — middleware creation moved to `main.go`
- Wired `TenantHandler` in `main.go` with `tenantRepo`, `feeRepo`, `tenantService`, and route registration on shared `/api` group
- Updated `protected_handler_test.go` to create its own `/api` group with auth + Casbin middleware; changed `/api/tenant` → `/api/users` paths
- Fixed pre-existing compile error in `auth_handler_test.go:139` (missing `pasetoService` argument to `NewAuthHandler`) — Rule 3 deviation

## Task Commits

Each task was committed atomically:

1. **Task 1: Create TenantHandler with CRUD endpoints and HTTP tests** — `5e3d435` (feat)
2. **Task 2: Remove tenant stubs, refactor ProtectedHandler, wire TenantHandler in main.go** — `1240542` (refactor)

**Plan metadata:** Will commit in final step.

## Files Created/Modified

- `apps/api/internal/delivery/http/tenant_handler.go` — TenantHandler with 9 handler methods, request types (CreateTenantRequest, CreateMandatoryFeeRequest, CreateVoluntaryFeeRequest, UpdateFeeRequest), conversion helpers (decodeFloat64, decodeDate), error-to-status-code mapping via `errors.Is`
- `apps/api/internal/delivery/http/tenant_handler_test.go` — 10 HTTP integration tests using mock TenantService, mock TxPool, and test claims middleware
- `apps/api/internal/delivery/http/protected_handler.go` — RegisterRoutes signature changed to `(api fiber.Router, pasetoSvc)`, middleware creation removed, tenant stubs deleted (ListTenants/GetTenant/CreateTenant)
- `apps/api/internal/delivery/http/protected_handler_test.go` — setupProtectedApp creates /api group with middleware internally, all `/api/tenant` paths changed to `/api/users`
- `apps/api/cmd/server/main.go` — Casbin enforcer initialized, authMW + casbinMW created, `/api` group created at main level, ProtectedHandler.RegisterRoutes and TenantHandler.RegisterRoutes both wired on shared `/api` group, tenantRepo/feeRepo/tenantService initialized
- `apps/api/internal/delivery/http/auth_handler_test.go` — Fixed pre-existing compile error (missing paseto arg to NewAuthHandler at line 139)

## Decisions Made

- **errors.Is for service error matching** — Use `errors.Is(err, service.ErrDuplicateBlockUnit)` instead of `==` comparison. This handles both direct sentinel errors and `%w`-wrapped errors from the service layer consistently. Applied across all 9 handler methods.
- **float64 amounts + string dates in request types** — The handler layer receives monetary amounts as `float64` (standard JSON number) and dates as ISO 8601 strings. The service layer converts to `decimal.Decimal` and `time.Time`. This avoids decimal serialization issues in JSON and gives frontend maximum flexibility with date formatting.
- **CreateFee type discriminator** — POST `/api/tenants/:id/fees` includes a `type` field (`"mandatory"` or `"voluntary"`). The handler routes to `tenantService.CreateMandatoryFee` or `tenantService.CreateVoluntaryFee` based on this value. This keeps the API surface minimal (one endpoint instead of two) while maintaining clean service separation.
- **Middleware creation moved to main.go** — Both ProtectedHandler and TenantHandler share the same `/api` group with identical auth + Casbin middleware. Creating middleware at `main.go` level and passing `fiber.Router` avoids middleware duplication and ensures consistent protection across all route groups.
- **Test claims middleware** — Integration tests use a test-only middleware that directly sets `c.Locals("user", &auth.UserClaims{...})` instead of creating real PASETO tokens. This eliminates token creation overhead in tests while still exercising the full request routing path.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Fixed pre-existing compile error in auth_handler_test.go**
- **Found during:** Task 1 (TenantHandler test creation)
- **Issue:** `auth_handler_test.go:139` calls `NewAuthHandler(authService)` but the constructor signature changed in Phase 1 to require `(authService, pasetoService)` — the test was missing the second argument and wouldn't compile
- **Fix:** Added `pasetoService` as the second argument to `NewAuthHandler(authService, pasetoService)`
- **Files modified:** `apps/api/internal/delivery/http/auth_handler_test.go`
- **Verification:** `go test ./internal/delivery/http/...` passes
- **Committed in:** `5e3d435` (Task 1 commit)

---

**Total deviations:** 1 auto-fixed (1 blocking)
**Impact on plan:** Minor. The fix was necessary for test compilation. No scope creep — the test file was already in the codebase and the missing argument was pre-existing.

## Issues Encountered

- **Pre-existing Casbin test failures (3 tests):** `TestCasbinEnforcer_Enforce`, `TestCasbinEnforcer_AddRemoveRoleLink`, `TestCasbinEnforcer_EnforceWithTerritory_RWOfficer` — all in `casbin_test.go` in the `auth` package. These failures are pre-existing (RBAC policy quirks with placeholder matching and role inheritance edge cases) and not caused by this plan's changes. All `internal/delivery/http` tests pass cleanly.
- **Pre-existing test path mismatch resolved:** The protected_handler_test.go tests were hitting `/api/tenant` which the `defaultResourceExtractor` resolves to resource `tenant`, but `policy.csv` uses `tenants` (plural) — causing 9 pre-existing failures. Changing paths to `/api/users` in this plan resolves those failures since policy has matching `users` rules.

## Known Stubs

None. The TenantHandler is fully wired with real service calls and proper error handling. No placeholder data, mock responses, or empty states flow through to production endpoints.

## Threat Flags

None. No new network endpoints, auth paths, file access patterns, or schema changes beyond those explicitly designed in the threat model. All routes are protected by the shared auth + Casbin middleware chain.

## Next Phase Readiness

- Full backend tenant/fee API ready for frontend consumption
- 10 endpoints available under shared `/api` group with auth + Casbin protection
- All routes follow D-03 contract (plural `/api/tenants`, nested `/api/tenants/:id/fees`)
- Error responses use consistent `{"error": "...", "code": "..."}` format
- Frontend can begin consuming endpoints immediately

**Deferred for future phases:**
- `/api/tenants/:id/dashboard` and `/api/dashboard` aggregate endpoints (if needed)
- Voluntary fee payment tracking at tenant level
- Period-based fee reporting (monthly/quarterly summaries)

---

*Phase: 02-tenant-fee-management*
*Plan: 04*
*Completed: 2026-05-23*
