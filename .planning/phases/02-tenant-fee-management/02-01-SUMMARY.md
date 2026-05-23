---
phase: 02-tenant-fee-management
plan: 01
subsystem: database
tags: go, postgresql, shopspring-decimal, migration, entity
requires:
  - phase: 01-core-authentication-rbac
    provides: User entity pattern, migration framework, Go project layout
provides:
  - Tenant, MandatoryFee, VoluntaryFee entity structs with decimal.Decimal for monetary fields
  - SQL migration files for tenants, mandatory_fees, voluntary_fees, user_tenants tables
  - shopspring/decimal dependency for financial precision
affects: plans 02-02 through 02-04 (repository, service, handler layers)

tech-stack:
  added:
    - github.com/shopspring/decimal v1.4.0
  patterns:
    - Entity struct with Sanitize() method (following User entity pattern)
    - Migration naming convention: NNN_description.{up,down}.sql
    - Composite unique constraints for territory-scoped uniqueness
    - Junction table with composite primary key for many-to-many relationships

key-files:
  created:
    - apps/api/internal/domain/entity/tenant.go
    - apps/api/migrations/006_create_tenants_table.up.sql
    - apps/api/migrations/006_create_tenants_table.down.sql
    - apps/api/migrations/007_create_mandatory_fees_table.up.sql
    - apps/api/migrations/007_create_mandatory_fees_table.down.sql
    - apps/api/migrations/008_create_voluntary_fees_table.up.sql
    - apps/api/migrations/008_create_voluntary_fees_table.down.sql
    - apps/api/migrations/009_create_user_tenants_table.up.sql
    - apps/api/migrations/009_create_user_tenants_table.down.sql
  modified:
    - apps/api/go.mod
    - apps/api/go.sum

key-decisions:
  - "Used shopspring/decimal v1.4.0 (upgraded from v1.2.0) for decimal monetary fields — float64 would cause rounding errors in financial calculations"
  - "Following User entity pattern exactly: package entity, JSON tags, Sanitize() method returning all-public-fields copy"
  - "Migration files use TIMESTAMPTZ (with timezone) instead of TIMESTAMP — consistent with plan spec, providing timezone-aware timestamps"
  - "mandatory_fees and voluntary_fees use separate tables (not a single fees table with type column) per CONTEXT.md locked decision"

patterns-established:
  - "Entity pattern: struct with JSON tags + Sanitize() returning non-nil pointer to a copy with all public fields"
  - "Migration pattern: up file creates table with IF NOT EXISTS, FK references, CHECK constraints, indexes; down file drops with CASCADE"
  - "UNIQUE constraint pattern for territory-scoped uniqueness: UNIQUE(block, unit_number, territory_id)"
  - "Junction table pattern: composite PRIMARY KEY(user_id, tenant_id) with separate indexes on each FK"

requirements-completed:
  - TEN-01
  - FIN-01
  - FIN-02

duration: 8min
completed: 2026-05-23
---

# Phase 2 Plan 01: Tenant & Fee Entity Definitions and Migrations

**Entity structs (Tenant, MandatoryFee, VoluntaryFee) with shopspring/decimal monetary precision, plus 8 migration files (4 up, 4 down) for the four new database tables**

## Performance

- **Duration:** 8 min
- **Started:** 2026-05-23T11:45:00Z
- **Completed:** 2026-05-23T11:53:00Z
- **Tasks:** 2
- **Files modified:** 11

## Accomplishments

- Tenant, MandatoryFee, and VoluntaryFee Go structs defined with `decimal.Decimal` for all monetary fields (financial precision)
- Each struct includes a `Sanitize()` method following the existing User entity pattern from Phase 1
- shopspring/decimal v1.4.0 dependency installed and verified
- Migration 006: `tenants` table with `UNIQUE(block, unit_number, territory_id)` constraint and territory FK
- Migration 007: `mandatory_fees` table with FK to tenants, `CHECK(amount >= 0)`, and `ON DELETE CASCADE`
- Migration 008: `voluntary_fees` table with same schema as mandatory_fees (separate table per CONTEXT.md decision)
- Migration 009: `user_tenants` junction table with composite `PRIMARY KEY(user_id, tenant_id)`
- All tables use uuidv7 primary keys, TIMESTAMPTZ timestamps, and appropriate indexes on foreign keys

## Task Commits

Each task was committed atomically:

1. **Task 1: Install shopspring/decimal and create entity types** - `0b09c57` (feat)
2. **Task 2: Create SQL migration files** - `2086c7e` (feat)

**Plan metadata:** (committed separately after SUMMARY)

## Files Created/Modified

- `apps/api/internal/domain/entity/tenant.go` - Tenant, MandatoryFee, VoluntaryFee structs with JSON tags, decimal.Decimal monetary fields, and Sanitize() methods
- `apps/api/migrations/006_create_tenants_table.up.sql` - tenants DDL with UNIQUE constraint
- `apps/api/migrations/006_create_tenants_table.down.sql` - DROP TABLE IF EXISTS tenants CASCADE
- `apps/api/migrations/007_create_mandatory_fees_table.up.sql` - mandatory_fees DDL with FK + CHECK
- `apps/api/migrations/007_create_mandatory_fees_table.down.sql` - DROP TABLE IF EXISTS mandatory_fees CASCADE
- `apps/api/migrations/008_create_voluntary_fees_table.up.sql` - voluntary_fees DDL with FK + CHECK
- `apps/api/migrations/008_create_voluntary_fees_table.down.sql` - DROP TABLE IF EXISTS voluntary_fees CASCADE
- `apps/api/migrations/009_create_user_tenants_table.up.sql` - user_tenants junction table DDL
- `apps/api/migrations/009_create_user_tenants_table.down.sql` - DROP TABLE IF EXISTS user_tenants CASCADE
- `apps/api/go.mod` - Updated with shopspring/decimal v1.4.0
- `apps/api/go.sum` - Updated checksums for shopspring/decimal

## Decisions Made

- Used `shopspring/decimal` (v1.4.0, upgraded from v1.2.0) for all monetary fields — float64 would cause rounding errors in financial calculations
- Followed existing User entity pattern exactly: package entity, JSON tags, Sanitize() returning copy with all public fields
- Used `TIMESTAMPTZ` (with timezone) over `TIMESTAMP` per plan spec for timezone-aware timestamps
- separate `mandatory_fees` and `voluntary_fees` tables (not a combined `fees` table) per CONTEXT.md locked decision — validation rules differ (mandatory requires at least one per tenant)

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

- shopspring/decimal v1.2.0 was already an indirect dependency (from another package). `go get` upgraded it to v1.4.0. No functional impact.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Database foundation complete: entity types and migration SQL ready for repository, service, and handler layers (plans 02-02 through 02-04)
- All downstream layers can now import `entity.Tenant`, `entity.MandatoryFee`, and `entity.VoluntaryFee`
- shopspring/decimal available for financial calculations throughout the codebase

## Self-Check: PASSED

All 9 created files verified on disk. Both commits (`0b09c57`, `2086c7e`) found in git log. No deletions or untracked files.

---

*Phase: 02-tenant-fee-management*
*Completed: 2026-05-23*
