---
phase: 02-tenant-fee-management
plan: 02
subsystem: api
tags: go, pgx, repository, tenant, fee, postgres
requires:
  - phase: 02-01
    provides: Tenant/MandatoryFee/VoluntaryFee entities, migrations
provides:
  - TenantRepository interface with 7 CRUD methods
  - FeeRepository interface with 9 methods split by mandatory/voluntary
  - PostgresTenantRepository pgx implementation with territory-filtered SQL
  - PostgresFeeRepository pgx implementation with separate table methods
affects:
  - 02-03 (service layer depends on these repositories)
  - 02-04 (handler wiring depends on repository constructors)

tech-stack:
  added: []
  patterns:
    - Repository interface + pgx implementation pattern matching Phase 1
    - pgtype.Numeric conversion for decimal.Decimal financial fields
    - Scan helpers for code reuse across single-row and multi-row queries
    - CreateTx variants accepting pgx.Tx for transactional operations
    - Territory-filtered queries at repository layer for data isolation

key-files:
  created:
    - apps/api/internal/domain/repository/tenant_repository.go
    - apps/api/internal/domain/repository/fee_repository.go
    - apps/api/internal/infrastructure/repository/tenant_repository.go
    - apps/api/internal/infrastructure/repository/fee_repository.go
  modified: []

key-decisions:
  - "Used pgtype.Numeric scan helper for NUMERIC columns to convert to decimal.Decimal — pgx v5 scans NUMERIC into pgtype.Numeric by default, requiring explicit conversion"
  - "Extracted scanTenant/scanTenantFromRow helpers to avoid duplicating column scan logic between QueryRow and multi-row Query patterns"
  - "Used Exec + RowsAffected for Delete operations, RETURNING + QueryRow for Update operations — Delete's simpler return type (error) allows Exec pattern, while Update returns the full entity"

requirements-completed:
  - TEN-01
  - FIN-01
  - FIN-02

duration: 8min
completed: 2026-05-23
---

# Phase 2 Plan 2: Repository Interfaces & pgx Implementations Summary

**TenantRepository and FeeRepository interfaces with pgx-backed PostgreSQL implementations matching Phase 1's user_repository pattern**

## Performance

- **Duration:** 8 min
- **Started:** 2026-05-23T02:40:00Z
- **Completed:** 2026-05-23T02:48:00Z
- **Tasks:** 3
- **Files modified:** 4

## Accomplishments

- **TenantRepository interface** with 7 methods including `CreateTx` for atomic tenant+fee creation and `ListByUserID` for resident-scoped queries through `user_tenants` junction table
- **FeeRepository interface** with 9 methods cleanly separating mandatory fee operations (CreateMandatory, CreateMandatoryTx, ListMandatoryByTenant, UpdateMandatory, DeleteMandatory) from voluntary fee operations (CreateVoluntary, ListVoluntaryByTenant, UpdateVoluntary, DeleteVoluntary)
- **PostgresTenantRepository** implementing all 7 TenantRepository methods with territory-filtered SQL (`WHERE territory_id = $N`) in ListByTerritory, Update, and Delete — preventing cross-territory data access at the database level
- **PostgresFeeRepository** implementing all 9 FeeRepository methods with parameterized SQL for separate `mandatory_fees` and `voluntary_fees` tables, including `CreateMandatoryTx` for transactional tenant+fee creation
- Both repositories include compile-time interface assertions (`var _ repository.X = (*PostgresX)(nil)`) ensuring full interface compliance

## Task Commits

Each task was committed atomically:

1. **Task 1: Create TenantRepository and FeeRepository interfaces** - `0b7024a` (feat)
2. **Task 2: Create PostgresTenantRepository implementing TenantRepository** - `0ce7c4a` (feat)
3. **Task 3: Create PostgresFeeRepository implementing FeeRepository** - `02ce353` (feat)

## Files Created

- `apps/api/internal/domain/repository/tenant_repository.go` — TenantRepository interface with 7 methods (Create, CreateTx, FindByID, ListByTerritory, ListByUserID, Update, Delete)
- `apps/api/internal/domain/repository/fee_repository.go` — FeeRepository interface with 9 methods split by mandatory/voluntary fee types
- `apps/api/internal/infrastructure/repository/tenant_repository.go` — PostgresTenantRepository with pgxpool-based implementation, scan helpers for decimal.Decimal conversion, territory-filtered queries
- `apps/api/internal/infrastructure/repository/fee_repository.go` — PostgresFeeRepository with separate SQL for mandatory_fees and voluntary_fees tables, Update/DELETE with RowsAffected checking

## Decisions Made

- Used `pgtype.Numeric` conversion vs `float64` scan: Scanned NUMERIC columns into `pgtype.Numeric` via pgx's built-in type mapping, then converted to `decimal.Decimal` via `Float64Value()` — preserving financial precision through the domain layer while keeping the scan implementation clean
- Shared scan helpers: Extracted `scanTenant`, `scanTenantFromRow`, `scanMandatoryFee`, and `scanVoluntaryFee` to handle the common column scan pattern across single-row (QueryRow RETURNING) and multi-row (Query) result sets
- Update with RETURNING: Used PostgreSQL `RETURNING` clause on UPDATE statements to return the full updated row (including server-generated `updated_at` timestamp) in a single round-trip, matching the existing Phase 1 pattern

## Deviations from Plan

None — plan executed exactly as written.

## Issues Encountered

- Pre-existing `auth_handler_test.go` vet error (`not enough arguments in call to NewAuthHandler`) — unrelated to this plan's changes, originates from a test file not modified in this phase.

## Next Phase Readiness

- All repository interfaces and pgx implementations complete
- Ready for **02-03**: Service layer with tenant/fee validation, transactional creation flows, and Casbin policy updates
