---
phase: 03-tenant-fee-ui
plan: 04
subsystem: ui
tags: [react, tanstack-query, tenant-management, form, dynamic-fees, vitest, testing-library]
requires:
  - phase: 03-02
    provides: Reusable UI components (Input, Select, StatusBadge, PageHeader, LoadingSkeleton, EmptyState, ConfirmDialog)
  - phase: 03-03
    provides: TanStack Query hooks (useTenants, useTenant, useCreateTenant, useUpdateTenant, useDeleteTenant)
provides:
  - TenantCard component with occupancy badge, formatted IDR fee, and keyboard accessibility
  - TenantForm with dynamic mandatory fee rows, add/remove, validation, and loading state
  - TenantListPage with responsive card grid (1→2→3 columns), search/filter with 300ms debounce, and loading/empty/error states
  - TenantCreatePage wrapping TenantForm in create mode with mutation and navigation
  - TenantEditPage with pre-populated form, delete confirmation dialog, and error handling
affects:
  - Tenant Detail Page (Plan 05) — TenantCard onClick navigates here
  - App routing — tenant pages already wired in App.tsx

tech-stack:
  added: []
  patterns:
    - Forms with dynamic entry arrays (add/remove fee rows)
    - Client-side search/filter with debounce via useEffect + setTimeout
    - Page-level component composition (PageHeader, LoadingSkeleton, EmptyState, TenantCard)
    - Mock TanStack Query hooks in tests using vi.mock with importActual

key-files:
  created:
    - apps/web/src/components/tenants/TenantCard.tsx — Tenant summary card with occupancy badge and fee info
    - apps/web/src/components/tenants/TenantCard.test.tsx — 7 tests for TenantCard rendering and interaction
    - apps/web/src/components/tenants/TenantForm.tsx — Dynamic mandatory fee form with validation
    - apps/web/src/components/tenants/TenantForm.test.tsx — 10 tests for form fields, dynamic rows, validation, submit
    - apps/web/src/pages/TenantListPage.test.tsx — 7 tests for list page states, search/filter, navigation
    - apps/web/src/pages/TenantCreatePage.test.tsx — 3 tests for create page submission and navigation
    - apps/web/src/pages/TenantEditPage.test.tsx — 6 tests for edit page loading, display, delete flow
  modified:
    - apps/web/src/pages/TenantListPage.tsx — Stub replaced with full implementation
    - apps/web/src/pages/TenantCreatePage.tsx — Stub replaced with full implementation
    - apps/web/src/pages/TenantEditPage.tsx — Stub replaced with full implementation

key-decisions:
  - "TenantCard uses optional mandatoryFeeCount/voluntaryFeeCount props (default 0) since Tenant type doesn't include fee counts"
  - "TenantForm accepts isLoading as prop (from parent mutation), not internal state — matches pattern of parent-driven loading"
  - "Client-side search debounced at 300ms via useEffect + setTimeout for responsive filtering"
  - "Checkpoints omitted — all tasks type='auto' with acceptance criteria verified via automated tests"
  - "Edit page maps CreateTenantRequest to UpdateTenantRequest fields (excludes mandatory_fees which are managed via separate fee endpoints)"

patterns-established:
  - "Dynamic form sections: state array + add/remove handlers + index-based field names for accessibility"
  - "Debounced search: useEffect cleanup pattern with setTimeout on search input state"
  - "Mocked TQ hooks: vi.mock('@tanstack/react-query') with importActual for non-mocked exports + QueryClientProvider wrapper for useQueryClient"

requirements-completed:
  - TEN-01
  - FIN-01

duration: 3min
completed: 2026-05-24
---

# Phase 3 Plan 4: Tenant Pages Summary

**TenantCard, TenantForm (with dynamic mandatory fees), TenantListPage, TenantCreatePage, and TenantEditPage — all with co-located test files consuming existing hooks and UI components**

## Performance

- **Duration:** 3 min
- **Started:** 2026-05-24T10:00:00Z
- **Completed:** 2026-05-24T10:05:22Z
- **Tasks:** 3 (all auto)
- **Files modified:** 10 (7 created, 3 modified from stubs)

## Accomplishments

- TenantCard with occupancy badge, formatted IDR monthly fee, fee summary, keyboard accessibility (Enter key)
- TenantForm with Block, Unit Number, Occupancy, Monthly Fee fields plus dynamic Mandatory Fees section with add/remove rows
- Client-side validation on submit: block/unit_number required, monthly_fee positive, >=1 mandatory fee, each fee description/amount/date required, fee amount <= monthly_fee cap
- TenantForm shows submitError alert from server, spinner + "Saving..." during submission
- TenantListPage with useTenants hook, responsive grid (1→2→3 columns), search/filter with 300ms debounce, sorted by block→unit_number, loading/empty/error states
- TenantCreatePage wraps TenantForm, calls useCreateTenant mutation, navigates to /tenants on success
- TenantEditPage fetches tenant via useTenant, renders pre-populated form, shows delete button with ConfirmDialog, navigates on success
- All 33 tests passing across 5 test files

## Task Commits

Each task was committed atomically:

1. **Task 1: Create TenantCard and TenantForm with tests** - `29104b1` (feat) — 17 tests
2. **Task 2: Create TenantListPage with search/filter** - `486381f` (feat) — 7 tests
3. **Task 3: Create TenantCreatePage and TenantEditPage** - `a19bbb7` (feat) — 9 tests

**Plan metadata:** pending

## Files Created/Modified

### Created (7)
- `apps/web/src/components/tenants/TenantCard.tsx` — TenantCard component with div + role="button" pattern
- `apps/web/src/components/tenants/TenantCard.test.tsx` — 7 tests for card rendering and interaction
- `apps/web/src/components/tenants/TenantForm.tsx` — Tenant form with dynamic mandatory fees section
- `apps/web/src/components/tenants/TenantForm.test.tsx` — 10 tests for form fields, validation, dynamic rows
- `apps/web/src/pages/TenantListPage.test.tsx` — 7 tests for list page states
- `apps/web/src/pages/TenantCreatePage.test.tsx` — 3 tests for create page
- `apps/web/src/pages/TenantEditPage.test.tsx` — 6 tests for edit page

### Modified (3)
- `apps/web/src/pages/TenantListPage.tsx` — Full implementation replacing stub
- `apps/web/src/pages/TenantCreatePage.tsx` — Full implementation replacing stub
- `apps/web/src/pages/TenantEditPage.tsx` — Full implementation replacing stub

## Decisions Made

- **TenantCard fee counts as optional props**: Since Tenant type doesn't include fee count information from the API, mandatoryFeeCount and voluntaryFeeCount are optional props defaulting to 0. The card shows "No fees configured" when no counts provided.
- **TenantForm isLoading as prop**: Loading state for the submit button is driven by the parent page's mutation isPending state, not internal form state. The form calls onSubmit(data) and catches errors to display submitError.
- **Debounced search**: 300ms debounce via useEffect + setTimeout pattern, with cleanup to handle rapid typing. Filtered results via useMemo for performance.
- **Edit page data mapping**: When editing, TenantForm pre-populates with initialData and the page maps CreateTenantRequest fields to UpdateTenantRequest (omitting mandatory_fees since they're managed separately via fee endpoints per Backend API design).
- **QueryClientProvider in tests**: Used to satisfy useQueryClient() calls in hooks. Creates a fresh QueryClient with retry: false for predictable test behavior.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - all 33 tests pass. The TenantForm loading state test was adjusted: originally tested internal loading during submit, but since isLoading is a parent-driven prop, the test was changed to verify the visual state when isLoading={true} is provided.

## Known Stubs

- **formatIDR helper**: Defined locally in TenantCard.tsx. If reused by FeeList or other components in Plan 05, it should be extracted to a shared utility file (`utils/format.ts`).

## Threat Flags

None — all components consume data through existing hooks (which use the request() helper from Plan 01). No new network endpoints or auth paths introduced.

## Next Phase Readiness

- Ready for Plan 05 (Tenant Detail Page / Fee management)
- TenantCard onClick navigates to `/tenants/:id` — the detail page
- TenantListPage wired and rendering in App.tsx

---

*Phase: 03-tenant-fee-ui*
*Completed: 2026-05-24*
