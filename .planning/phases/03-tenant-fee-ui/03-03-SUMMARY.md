---
phase: 03-tenant-fee-ui
plan: 03
subsystem: ui
tags: tanstack-query, hooks, protected-route, query-keys, cache-invalidation

requires:
  - phase: 03-01
    provides: api helper, types, services, QueryClientProvider
  - phase: 02-tenant-fee-management
    provides: API endpoints for tenants and fees

provides:
  - "9 custom TanStack Query hooks (useTenants, useTenant, useFees, useCreateTenant, useUpdateTenant, useDeleteTenant, useCreateFee, useUpdateFee, useDeleteFee)"
  - "ProtectedRoute migrated to useQuery for auth/me session check"

affects:
  - 03-04
  - 03-05

tech-stack:
  added: []
  patterns:
    - "Custom hooks wrapping useQuery/useMutation with consistent query key structure"
    - "Mutations always invalidate related query caches in onSuccess callback"
    - "ProtectedRoute uses useQuery with queryKey ['auth', 'me'] for session check"
    - "TQ hooks tested implicitly via component tests with mocked useQuery"

key-files:
  created:
    - apps/web/src/hooks/useTenants.ts
    - apps/web/src/hooks/useTenant.ts
    - apps/web/src/hooks/useFees.ts
    - apps/web/src/hooks/useCreateTenant.ts
    - apps/web/src/hooks/useUpdateTenant.ts
    - apps/web/src/hooks/useDeleteTenant.ts
    - apps/web/src/hooks/useCreateFee.ts
    - apps/web/src/hooks/useUpdateFee.ts
    - apps/web/src/hooks/useDeleteFee.ts
  modified:
    - apps/web/src/routes/ProtectedRoute.tsx
    - apps/web/src/routes/ProtectedRoute.test.tsx

key-decisions:
  - "Hooks tested implicitly via component tests (D-17) — no dedicated hook test files"
  - "useQuery(['auth', 'me']) replaces manual fetch + useEffect in ProtectedRoute"
  - "useQueryClient.invalidateQueries in onSuccess prevents stale cache after mutations"
  - "Auth query uses staleTime: 5 min since user session rarely changes mid-session"

patterns-established:
  - "Query hooks: useQuery<T>({ queryKey: ['resource'], queryFn: serviceFn })"
  - "Query hooks with params: enabled: !!param prevents fetch on empty params"
  - "Mutation hooks: useQueryClient() + invalidateQueries in onSuccess"
  - "MutationFn signatures match service function contracts exactly"

requirements-completed:
  - TEN-01
  - FIN-01
  - FIN-02

duration: 7min
completed: 2026-05-24
---

# Phase 3: TanStack Query hooks + ProtectedRoute migration to useQuery

**9 custom TanStack Query hooks for tenant and fee CRUD operations, plus ProtectedRoute migration from manual fetch + useEffect to useQuery-based auth/me session check**

## Performance

- **Duration:** 7 min
- **Started:** 2026-05-24T09:51:00Z
- **Completed:** 2026-05-24T09:58:46Z
- **Tasks:** 2 (TDD: RED + GREEN)
- **Files modified:** 11

## Accomplishments

- Created 9 custom TanStack Query hooks with consistent query key patterns and cache invalidation
- Migrated ProtectedRoute to use useQuery for /auth/me session check — removed manual fetch + useEffect + useState
- Updated ProtectedRoute tests to mock @tanstack/react-query with vi.mock, removing global.fetch dependency
- All 6 ProtectedRoute tests pass with mocked useQuery (2-pass → 6-pass after migration)

## Task Commits

Each task was committed atomically:

1. **Task 1: Create all 9 custom TanStack Query hooks** - `8184178` (feat)
2. **Task 2 (RED): Update test with useQuery mocks** - `4dc1ead` (test)
3. **Task 2 (GREEN): Migrate ProtectedRoute to useQuery** - `0d6bed5` (feat)

**Plan metadata:** (committed below)

## Files Created/Modified

- `apps/web/src/hooks/useTenants.ts` — useQuery hook with `['tenants']` key, calls `listTenants`
- `apps/web/src/hooks/useTenant.ts` — useQuery hook with `['tenants', id]` key, calls `getTenant(id)`, enabled guard
- `apps/web/src/hooks/useFees.ts` — useQuery hook with `['fees', tenantId]` key, calls `listFees(tenantId)`, enabled guard
- `apps/web/src/hooks/useCreateTenant.ts` — useMutation with cache invalidation of `['tenants']` on success
- `apps/web/src/hooks/useUpdateTenant.ts` — useMutation with cache invalidation of `['tenants']` on success
- `apps/web/src/hooks/useDeleteTenant.ts` — useMutation with cache invalidation of `['tenants']` on success
- `apps/web/src/hooks/useCreateFee.ts` — useMutation with cache invalidation of `['fees', tenantId]` on success
- `apps/web/src/hooks/useUpdateFee.ts` — useMutation with cache invalidation of `['fees', tenantId]` on success
- `apps/web/src/hooks/useDeleteFee.ts` — useMutation with cache invalidation of `['fees', tenantId]` on success
- `apps/web/src/routes/ProtectedRoute.tsx` — Migrated from manual fetch + useEffect to useQuery(['auth', 'me'])
- `apps/web/src/routes/ProtectedRoute.test.tsx` — Updated to mock @tanstack/react-query useQuery and services/api request

## Decisions Made

- **D-17 followed:** Custom hooks have no dedicated test files — they're tested implicitly via the services they call; components are tested via mocked query hooks (per CONTEXT.md D-17)
- **Mutation invalidation pattern:** All mutation hooks call `queryClient.invalidateQueries` in their `onSuccess` callback to prevent stale data — following the pitfall-avoidance pattern from RESEARCH.md §Pitfall 2
- **Auth staleTime:** The auth/me query uses `staleTime: 5 * 60 * 1000` (5 minutes) since user session rarely changes mid-session — per threat model T-03-10 mitigation
- **Payment field update contract:** `useUpdateFee` accepts `{ feeId, data }` as single argument (not positional) to match TanStack Query's mutationFn single-arg requirement

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - all tasks completed cleanly.

## TDD Gate Compliance

- **Task 1:** TDD flag present but hooks have no dedicated test files (per D-17 exemption documented in plan). All 9 hook files created directly with no test gate.
- **Task 2:** Full TDD cycle followed — RED commit (test update with mocks → 4 tests fail as expected), GREEN commit (ProtectedRoute migration → all 6 tests pass).

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

Ready for Plans 04 (Tenant pages: list, create, edit) and 05 (Fee pages: detail, fee management). All data-fetching hooks are in place with correct query keys and cache invalidation. ProtectedRoute is migrated to useQuery for the session check, providing consistent loading, error, and data states for the authenticated layout.

---

*Phase: 03-tenant-fee-ui*
*Completed: 2026-05-24*
