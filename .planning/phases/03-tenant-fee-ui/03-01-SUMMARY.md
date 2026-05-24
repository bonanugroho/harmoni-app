---
phase: 03-tenant-fee-ui
plan: 01
subsystem: ui
tags: [tanstack-query, lucide-react, typescript, api-client, vitest]

requires:
  - phase: 02-tenant-fee-management
    provides: Tenant and Fee API endpoints (/api/tenants, /api/tenants/:id/fees)
provides:
  - Shared request() API helper with credentials and error handling
  - Tenant and Fee TypeScript interfaces matching API response shape
  - Tenant and Fee service functions using the shared helper
  - Refactored auth.ts using the shared helper
  - QueryClientProvider wrapping the app with configured defaults
  - Sidebar CSS custom properties in @theme block
affects: [03-02, 03-03, 03-04, 03-05]

tech-stack:
  added: [@tanstack/react-query@^5.100.14, lucide-react@^1.16.0]
  patterns:
    - "Shared request() helper for all API calls with credentials: 'include'"
    - "Vitest mock of request() helper for service-level tests"
    - "QueryClientProvider wrapping BrowserRouter in App.tsx"

key-files:
  created:
    - apps/web/src/services/api.ts
    - apps/web/src/services/api.test.ts
    - apps/web/src/types/tenant.ts
    - apps/web/src/types/fee.ts
    - apps/web/src/services/tenants.ts
    - apps/web/src/services/fees.ts
    - apps/web/src/services/tenants.test.ts
    - apps/web/src/services/fees.test.ts
  modified:
    - apps/web/package.json
    - apps/web/package-lock.json
    - apps/web/src/services/auth.ts
    - apps/web/src/services/auth.test.ts
    - apps/web/src/App.tsx
    - apps/web/src/index.css
  deleted:
    - apps/web/src/App.css

key-decisions:
  - "request() helper auto-detects JSON vs non-JSON via Content-Type header (handles 204 No Content) as specified in D-01/D-02"
  - "auth.ts refactored to use shared request() — removes duplicate API_BASE_URL and parseError"
  - "QueryClient configured with staleTime: 30000, retry: 1, refetchOnWindowFocus: true"
  - "Sidebar CSS tokens follow existing @theme pattern: sidebar, sidebar-hover, sidebar-active"

patterns-established:
  - "Service files import { request } from './api' for all HTTP calls"
  - "Service tests mock './api' with vi.mock and assert endpoint calls via request()"
  - "App.tsx wraps all routes in QueryClientProvider for TanStack Query context"

requirements-completed:
  - TEN-01
  - FIN-01
  - FIN-02

duration: 2 min
completed: 2026-05-24
---

# Phase 3 Plan 01: Foundation Summary

**Shared request() API helper with credentials: 'include' and 204 handling, Tenant/Fee TypeScript interfaces matching API response shapes, tenant/fee service functions, refactored auth.ts, QueryClientProvider wrapper, and sidebar CSS tokens**

## Performance

- **Duration:** 2 min
- **Started:** 2026-05-24T09:47:20Z
- **Completed:** 2026-05-24T09:49:33Z
- **Tasks:** 3 (all auto)
- **Files modified:** 18 (8 created, 8 modified, 1 deleted)

## Accomplishments

- Installed @tanstack/react-query (5.100.14) and lucide-react (1.16.0) as dependencies
- Created shared `services/api.ts` with `request<T>()` helper handling credentials, 204s, and error parsing with 7 passing tests
- Defined `types/tenant.ts` (Tenant, CreateTenantRequest, UpdateTenantRequest) and `types/fee.ts` (Fee, CreateFeeRequest, UpdateFeeRequest, ListFeesResponse)
- Created `services/tenants.ts` (5 functions: list, get, create, update, delete) and `services/fees.ts` (4 functions: list, create, update, delete) using the shared request() helper
- Refactored `services/auth.ts` to use `request()` from `./api`, removing duplicate `API_BASE_URL` and `parseError` definitions (all function signatures unchanged)
- All 62 existing and new tests pass (11 test files)
- `App.tsx` wraps `BrowserRouter` in `QueryClientProvider` with `staleTime: 30_000`, `retry: 1`, `refetchOnWindowFocus: true`
- Added `--color-sidebar`, `--color-sidebar-hover`, `--color-sidebar-active` tokens to `index.css` `@theme` block
- Removed boilerplate `App.css` (Vite template styles)
- Production build succeeds

## Task Commits

Each task was committed atomically:

1. **Task 1: Install dependencies and create shared api.ts helper with tests** - `3109317` (feat)
2. **Task 2: Create TypeScript types, tenant/fee services, refactor auth.ts, and service tests** - `2122419` (feat)
3. **Task 3: Add QueryClientProvider to App.tsx, add sidebar CSS tokens, remove App.css** - `0e427c0` (feat)

**Plan metadata:** *(created after summary)*

## Files Created/Modified

- `apps/web/src/services/api.ts` — Shared request() helper with credentials, 204 handling, error parsing
- `apps/web/src/services/api.test.ts` — 7 tests covering all request scenarios
- `apps/web/src/types/tenant.ts` — Tenant, CreateTenantRequest, UpdateTenantRequest interfaces
- `apps/web/src/types/fee.ts` — Fee, CreateFeeRequest, UpdateFeeRequest, ListFeesResponse interfaces
- `apps/web/src/services/tenants.ts` — Tenant CRUD service functions via request()
- `apps/web/src/services/fees.ts` — Fee CRUD service functions via request()
- `apps/web/src/services/tenants.test.ts` — 5 tests mocking request() for all tenant operations
- `apps/web/src/services/fees.test.ts` — 4 tests mocking request() for all fee operations
- `apps/web/src/services/auth.ts` — Refactored to use request() from ./api
- `apps/web/src/services/auth.test.ts` — Updated mock responses with headers for request() compatibility
- `apps/web/src/App.tsx` — Wraps BrowserRouter in QueryClientProvider
- `apps/web/src/index.css` — Added --color-sidebar, --color-sidebar-hover, --color-sidebar-active tokens
- `apps/web/src/App.css` — Deleted (boilerplate Vite template styles)
- `apps/web/package.json` — Added @tanstack/react-query, lucide-react dependencies

## Decisions Made

- Followed D-01/D-02 exactly: `request()` helper auto-detects JSON vs non-JSON via Content-Type, handles 204 No Content by returning `undefined`
- auth.ts refactored via D-03: removed duplicate `API_BASE_URL` and `parseError`, all function signatures unchanged
- D-07 implemented: `QueryClientProvider` wraps `BrowserRouter` with configured defaults
- Sidebar tokens added to existing `@theme` block following the same naming convention as Phase 1

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical] Added response headers to auth test mock responses**
- **Found during:** Task 2 (auth.ts refactoring)
- **Issue:** Existing auth tests mock `global.fetch` with bare response objects lacking `headers.get('content-type')`. After refactoring auth.ts to use `request()`, these tests fail because `request()` reads `response.headers.get('content-type')` to detect JSON responses.
- **Fix:** Added `headers: new Headers({ 'content-type': 'application/json' })` to all mock response objects via a shared `mockJsonResponse()` helper
- **Files modified:** `apps/web/src/services/auth.test.ts`
- **Verification:** All 8 auth tests pass, all 17 service tests pass
- **Committed in:** `2122419` (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (1 missing critical)
**Impact on plan:** Necessary compatibility fix for auth.ts refactoring. All auth function signatures and behaviors remain unchanged.

## Issues Encountered

- None — plan executed as specified. The auth test fix was a predictable consequence of changing the underlying HTTP helper.

## User Setup Required

None — no external service configuration required.

## Next Phase Readiness

Foundation layer is complete:
- Shared API client (`request()`) ready for all service functions
- TypeScript types defined matching API response shapes
- Tenant and Fee service functions tested
- QueryClientProvider configured and wrapping the app
- Sidebar CSS tokens available for AppLayout in Plan 02

Ready for **Plan 02** (UI components + AppLayout).
