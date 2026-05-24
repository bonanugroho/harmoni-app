---
phase: 03-tenant-fee-ui
plan: 05
subsystem: ui
tags: fees, fee-list, fee-form, tenant-detail, tdd, vitest, react, validation

# Dependency graph
requires:
  - phase: 03-02
    provides: UI components (Input, Select, DatePicker, StatusBadge, ConfirmDialog, PageHeader, LoadingSkeleton, EmptyState)
  - phase: 03-03
    provides: TanStack Query hooks (useTenant, useFees, useCreateFee, useUpdateFee, useDeleteFee)
  - phase: 03-04
    provides: Tenant pages (TenantListPage, TenantCreatePage, TenantEditPage)
provides:
  - FeeList component with mandatory/voluntary sections, fee cards, empty states, edit/delete actions
  - FeeForm component with fee type selector, validation (amount cap, date checks), create/edit modes
  - TenantDetailPage with tenant info header, fee management UI, Record Fee modal, delete confirmation
  - formatIDR and formatDate utility functions
affects: []

# Tech tracking
tech-stack:
  added: []
  patterns:
    - TDD cycle for presentational components (FeeList, FeeForm)
    - Mock custom hooks directly instead of mocking @tanstack/react-query for page-level tests
    - Modal overlay pattern for FeeForm (fixed z-50, backdrop, centered)
    - formatIDR / formatDate as exported utility functions from component file

key-files:
  created:
    - apps/web/src/components/fees/FeeList.tsx - Sectioned fee list with mandatory/voluntary sections, fee cards, empty states, edit/delete actions
    - apps/web/src/components/fees/FeeList.test.tsx - 9 tests for FeeList (sections, empty states, loading, actions, status badges)
    - apps/web/src/components/fees/FeeForm.tsx - Fee form with type selector, validation, create/edit modes, submit/cancel
    - apps/web/src/components/fees/FeeForm.test.tsx - 10 tests for FeeForm (validation, submit, edit mode, errors, cancel)
    - apps/web/src/pages/TenantDetailPage.test.tsx - 8 tests for TenantDetailPage (loading, render, error, modal open, create, delete)
  modified:
    - apps/web/src/pages/TenantDetailPage.tsx - Implemented full fee management page (was stub)

key-decisions:
  - "Force formatIDR as dependency: 'id-ID' with minimumFractionDigits: 0 for Rp X.XXX format"
  - "formatDate uses toLocaleDateString('en-GB', ...) for 'DD MMM YYYY' format"
  - "FeeForm validation compares dates without time component (setHours(0,0,0,0)) for effective_date past check"
  - "TenantDetailPage modal: uses fixed inset-0 z-50 with backdrop overlay approach matching ConfirmDialog pattern"
  - "Page-level tests mock custom hooks directly (useTenant, useFees, useCreateFee, etc.) rather than mocking @tanstack/react-query, since multiple useQuery calls need different return values per call"

patterns-established:
  - "FeeList uses formatIDR and formatDate exported from component file, making them importable by parent pages"
  - "FeeForm validates amount > 0 and amount <= monthlyFee on client side before submit"

requirements-completed:
  - FIN-01
  - FIN-02

# Metrics
duration: 4min
completed: 2026-05-24
---

# Phase 3 Plan 5: Fee Components & Tenant Detail Page

**FeeList and FeeForm components with TDD, TenantDetailPage with modal fee management, all backed by 27 passing tests**

## Performance

- **Duration:** 4 min
- **Started:** 2026-05-24T10:07:54Z
- **Completed:** 2026-05-24T10:11:39Z
- **Tasks:** 3 (2 TDD, 1 standard)
- **Files modified:** 6

## Accomplishments

- FeeList component with separate mandatory/voluntary sections, fee cards showing description, formatted IDR amount, effective date, paid/Unpaid status badge, and edit/delete action icons
- Empty state support for each section with context-specific copywriting and CTAs
- FeeForm with fee type selector, description/amount/date fields, comprehensive client-side validation (amount > 0, amount ≤ monthlyFee cap, effective date not in past, payment date after effective date)
- FeeForm supports both create mode and edit mode (pre-populated via initialData prop)
- TenantDetailPage assembling tenant info header, FeeList, modal FeeForm, and ConfirmDialog delete flow — with loading, error, and empty states throughout
- All 3 tests files pass (27 tests total across all new/modified files)
- Full test suite: 28 files, 162 tests — all passing

## Task Commits

Each task was committed atomically:

1. **Task 1: FeeList component (TDD)** — `90caf2e` (test), `7294cb3` (feat)
2. **Task 2: FeeForm component (TDD)** — `577caf6` (test), `586091b` (feat)
3. **Task 3: TenantDetailPage** — `5d78f06` (feat)

**Plan metadata (to follow):** pending

## Files Created/Modified

- `apps/web/src/components/fees/FeeList.tsx` — Sectioned fee list with mandatory/voluntary sections, fee cards, empty states, LoadingSkeleton, edit/delete actions
- `apps/web/src/components/fees/FeeList.test.tsx` — 9 tests for FeeList (sections, empty states per section, loading, edit/delete callbacks, paid/unpaid StatusBadge)
- `apps/web/src/components/fees/FeeForm.tsx` — Fee form with Select/Input/DatePicker components, validation, create/edit modes, error display
- `apps/web/src/components/fees/FeeForm.test.tsx` — 10 tests for FeeForm (fields, options, validation errors, submit, edit mode pre-population, error alert, cancel)
- `apps/web/src/pages/TenantDetailPage.tsx` — Full fee management page with tenant info header, FeeList, modal FeeForm (create/edit), ConfirmDialog delete, loading/error states
- `apps/web/src/pages/TenantDetailPage.test.tsx` — 8 tests for TenantDetailPage (loading skeleton, tenant info render, fee list render, error state, modal open, fee create, delete confirmation, fee delete)

## Decisions Made

- Force `formatIDR` uses `Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 })` for consistent Rp formatting
- `formatDate` uses `toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })` for "01 Jun 2026" format
- FeeForm validation compares dates without time component (via `setHours(0,0,0,0)`) for effective_date past check
- TenantDetailPage uses modal pattern (fixed inset-0 z-50 with backdrop overlay) to display FeeForm, consistent with ConfirmDialog approach
- Page-level tests mock custom hooks directly rather than mocking @tanstack/react-query for cleaner per-call state control

## TDD Gate Compliance

| Plan | RED | GREEN | REFACTOR | Status |
|------|-----|-------|----------|--------|
| Task 1 (FeeList) |  ✓  |   ✓   |    —     | Pass   |
| Task 2 (FeeForm) |  ✓  |   ✓   |    —     | Pass   |

Both TDD tasks followed RED (write failing tests) → GREEN (implement to pass) correctly. No REFACTOR commits needed — implementations were minimal and clean.

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

- TenantDetailPage test for `Rp` text had to use `getAllByText` instead of `getByText` because multiple fee cards display formatted amounts alongside the monthly fee in the header
- "Record Fee" modal test had to check for FeeForm fields instead of checking for the "Record Fee" heading text, since both the button and modal heading contain the same text

## Threat Surface Scan

No new network endpoints, auth paths, file access patterns, or schema changes at trust boundaries were introduced. All data flows through existing hooks and services.

## Known Stubs

None — all components are fully wired.

## Next Phase Readiness

- Phase 3 is complete — all plans (01 through 05) are done
- Ready for Phase 4 (Transaction Engine & Expenditures)
- Tenant management (CRUD + fee management) is fully functional with tests
- Ready for `/gsd-verify-work` before closing Phase 3

---

*Phase: 03-tenant-fee-ui*
*Completed: 2026-05-24*
