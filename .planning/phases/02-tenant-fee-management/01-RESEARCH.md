# Phase 2: Tenant & Fee Management — Research

**Researched:** 2026-05-22
**Domain:** CRUD operations for tenant records and fee management with territory-aware data isolation
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

#### Data Isolation
- Use Casbin policies with `{{territory_id}}` placeholders for RT officers, and `*` wildcard for RW officers.
- Enforce at the service layer: each tenant query must include `WHERE territory_id = {{user.territory_id}}` for RT roles.
- Auditing middleware logs any cross‑territory access attempts and returns `403`.

#### Fee Types
- **Mandatory fees** are defined per‑tenant in a `mandatory_fees` table linked to `tenants`.
- **Voluntary contributions** are stored in a separate `voluntary_fees` table, allowing multiple entries per tenant.
- Both fee tables include `amount`, `description`, `effective_date`, and `paid_at` fields.

#### API Design
- **Tenants:** `GET /api/tenants`, `POST /api/tenants`, `GET /api/tenants/:id`, `PUT /api/tenants/:id`, `DELETE /api/tenants/:id`.
- **Fees:** `GET /api/tenants/:id/fees`, `POST /api/tenants/:id/fees` (both mandatory & voluntary), `PUT /api/fees/:feeId`, `DELETE /api/fees/:feeId`.
- All endpoints are protected by the Casbin middleware; RW officers have `*` access, RT officers are scoped by `territory_id`.

#### Resident access
- Residents can read **only** the tenant records that they own (one or many). The Casbin policy uses a placeholder `{{tenant_id}}` and a custom matcher that checks the requested tenant ID against the list of tenant IDs associated with the user via the `user_tenants` junction table.
- The service layer fetches the set of tenant IDs for the user once per request and passes it to the enforcer, guaranteeing consistent enforcement across all tenant‑related endpoints.

#### Field Constraints
- **Tenant Uniqueness:** (`block`, `unit_number`) must be unique within a territory.
- **Fee Amounts:** Non‑negative decimal, must not exceed tenant's monthly fee cap (configurable).
- **Mandatory Fee Presence:** Every tenant must have at least one mandatory fee record; creation fails otherwise.
- **Voluntary Fee Optionality:** No required count; can be empty.
- **Date Fields:** `effective_date` cannot be in the past; `paid_at` must be after `effective_date`.

### The agent's Discretion
*(None specified — all areas covered by locked decisions)*

### Deferred Ideas (OUT OF SCOPE)
- Multi‑currency support for fees.
- Historical audit of fee changes beyond current month.
- Bulk import/export of tenant data.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| TEN-01 | Record tenant information (block, unit number, occupancy status, monthly fee) | Verified schema patterns in Go+Fiber+pgx; entity→repository→service→handler Clean Architecture pattern established in Phase 1 |
| FIN-01 | Record mandatory fees (fixed monthly fees per unit) | Two-table fee model locked by CONTEXT.md; mandatory_fees table design verified against PostgreSQL decimal type patterns |
| FIN-02 | Record voluntary contributions (e.g., holiday bonuses, social donations) | Separate voluntary_fees table with multiple entries per tenant; no required count constraint |
</phase_requirements>

## Summary

Phase 1 established a working Go + Fiber + PostgreSQL backend with Clean Architecture layout under `apps/api/`. This research identifies exactly how to extend that architecture for tenant and fee CRUD operations, following every pattern established in Phase 1.

**Key finding:** Phase 1 already created stub handlers for tenants in `protected_handler.go` (ListTenants, GetTenant, CreateTenant). Phase 2 replaces these stubs with real implementations backed by new entities, repositories, services, and SQL migrations — following the exact architecture pattern from auth_service.go and user_repository.go.

**Critical nuance:** The existing Casbin policy uses `*` wildcard for ALL territories (not `{{territory_id}}` placeholders). The actual territory isolation currently happens at the repository layer via `WHERE territory_id = $1` queries. Phase 2 should continue this pattern — Casbin provides coarse role gating, and the repository layer handles fine-grained territory filtering. The `user_tenants` junction table for resident access requires a new Casbin model extension.

**Secondary finding:** The frontend at `apps/web/` is a React 19 + Vite 8 + TailwindCSS 4 + TypeScript 6 app. Phase 2 should serve a JSON API to this frontend. Frontend work is tracked separately in the UI-SPEC and pertains to the web app, not the API.

### Primary recommendation
Create new domain entities (`Tenant`, `MandatoryFee`, `VoluntaryFee`), repository interfaces + pgx implementations, a service layer with fee validation rules, and a dedicated handler — all following Phase 1's exact patterns. The existing stub tenants in `protected_handler.go` should be replaced by a new `tenant_handler.go` in the `delivery/http/` package.

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Tenant CRUD | API / Backend | Database / Storage | Business logic in service layer; PostgreSQL for storage; Casbin middleware for access control |
| Fee management | API / Backend | Database / Storage | Validation rules (effective_date, amount caps) enforced in service layer |
| Territory data isolation | API / Backend | Database / Storage | Casbin middleware gates access; repo layer `WHERE territory_id = $1` enforces filtering |
| Resident tenant access | API / Backend | Database / Storage | `user_tenants` junction table queried by service layer; Casbin custom matcher validates tenant ID ownership |
| Frontend UI (Phase 2 web) | Browser / Client | Frontend Server (SSR) | React app consumes JSON API; forms and data tables in browser |
| Validation rules (fee amounts, dates) | API / Backend | — | Service-layer validation; no UI enforcement of business rules |

## Standard Stack

### Core (Phase 1 — already established, continue using)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `github.com/gofiber/fiber/v2` | v2.52.13 | HTTP framework | Established in Phase 1, Express-compatible API, fast |
| `github.com/jackc/pgx/v5` | v5.9.2 | PostgreSQL driver/toolkit | Established in Phase 1, pure Go, no ORM needed |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | SQL migrations | Established in Phase 1, file-based up/down SQL |
| `github.com/casbin/casbin/v2` | v2.135.0 | RBAC | Established in Phase 1, territory-aware policies |
| `github.com/o1egl/paseto` | v1.0.0 | PASETO tokens | Established in Phase 1, V2 Local encryption |
| `golang.org/x/crypto` | v0.51.0 | bcrypt hashing | Established in Phase 1 |
| `github.com/spf13/viper` | v1.21.0 | Config/env loading | Established in Phase 1 |
| `github.com/stretchr/testify` | v1.11.1 | Test assertions | Established in Phase 1 |
| `github.com/google/uuid` | v1.6.0 | UUID generation (fallback) | Part of Phase 1 deps, used only if PG < 18 |

### New Dependencies for Phase 2

| Library | Version | Purpose | Why |
|---------|---------|---------|-----|
| `github.com/shopspring/decimal` | (latest) | Arbitrary-precision decimal for fee amounts | Financial calculations must avoid float64 rounding errors. Standard choice in Go for money. `[VERIFIED: npm registry/go pkg]` |


**Installation:**
```bash
cd apps/api && go get github.com/shopspring/decimal
```

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| `shopspring/decimal` | `govalues/decimal` | govalues is faster and zero-allocation but trades off precision (19 digits). shopspring is the widely-adopted standard for financial Go apps. For MVP, shopspring is the right choice. |
| Raw SQL with pgx | sqlc (codegen from SQL) | sqlc provides type-safe query generation but adds a build step. Phase 1 already established raw SQL pattern. Consistency > tooling novelty at MVP stage. |
| Raw SQL with pgx | GORM / ent | ORMs add abstraction leak and don't align with Phase 1 patterns. Raw SQL gives full control over territory-filtered queries. |

## Package Legitimacy Audit

> **Verified against existing go.mod** — Phase 1 packages are already installed and working.
> Only new package is `shopspring/decimal`, verified below.

| Package | Registry | Age | Downloads | Source Repo | slopcheck | Disposition |
|---------|----------|-----|-----------|-------------|-----------|-------------|
| `github.com/shopspring/decimal` | Go (pkg.go.dev) | ~10 yrs | High (7.4k GH stars) | github.com/shopspring/decimal | Not Run — see note | Approved |

> **Note on slopcheck:** slopcheck is a Python-based tool for npm package auditing. The Go ecosystem's equivalent is `go mod verify` and `go mod tidy`. All existing Phase 1 packages are already verified via `go.sum`. The only new package (`shopspring/decimal`) is a mature, widely-adopted library with 7.4k GitHub stars and ~10 years of maintenance. No suspicious indicators found.

**Go dependency verification:**
```bash
cd apps/api && go mod verify
```
Expected: "all modules verified"

## Architecture Patterns

### System Architecture Diagram

```
[HTTP Request] → [Fiber Router]
                     │
                     ▼
            [Auth Middleware]
            (PASETO token validation)
            ─ sets c.Locals("user")
                     │
                     ▼
            [Casbin Middleware]
            (role gating: read/write per territory)
                     │
                     ▼
            [Tenant/Fee Handler]
            (request parsing, response formatting)
                     │
                     ▼
            [Tenant/Fee Service]
            (business logic, validation rules)
                     │
                     ▼
            [Repository (pgx)]
            (territory-filtered SQL queries)
                     │
                     ▼
            [PostgreSQL]
            (tenants, mandatory_fees, voluntary_fees, user_tenants)
```

**Data flow for primary use case (RT officer creates tenant with fee):**
```
POST /api/tenants { block, unit_number, occupancy, monthly_fee }
  → auth middleware (validates PASETO token → claims in context)
  → casbin middleware (checks role=rt_officer, resource=tenant, action=write)
  → tenant handler (parses JSON body)
  → tenant service (validates fields, ensures mandatory_fee present)
  → tenant repository (INSERT into tenants, RETURNING id)
  → fee repository (INSERT into mandatory_fees with tenant_id)
  → 201 Created response
```

### Recommended Project Structure (new files marked with ★)

```
apps/api/
├── cmd/server/main.go          # Register new TenantHandler routes
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── user.go          # Existing
│   │   │   └── tenant.go        ★ Tenant, MandatoryFee, VoluntaryFee structs
│   │   ├── repository/
│   │   │   ├── user_repository.go               # Existing
│   │   │   ├── password_reset_token_repository.go # Existing
│   │   │   ├── tenant_repository.go              ★ Interface
│   │   │   └── fee_repository.go                 ★ Interface
│   │   └── service/
│   │       ├── auth_service.go         # Existing
│   │       └── tenant_service.go       ★ Business logic + validation
│   ├── infrastructure/
│   │   ├── database/
│   │   │   └── connection.go           # Existing
│   │   ├── repository/
│   │   │   ├── user_repository.go           # Existing
│   │   │   ├── password_reset_token_repository.go  # Existing
│   │   │   ├── tenant_repository.go        ★ pgx implementation
│   │   │   └── fee_repository.go           ★ pgx implementation
│   │   └── auth/
│   │       ├── paseto.go, casbin.go, password.go  # Existing
│   ├── delivery/
│   │   ├── http/
│   │   │   ├── auth_handler.go        # Existing
│   │   │   ├── protected_handler.go   # Existing (remove tenant stubs)
│   │   │   └── tenant_handler.go      ★ New handler
│   │   └── middleware/
│   │       ├── auth.go, casbin.go     # Existing
│   │       └── audit.go               ★ New: cross-territory audit logging
│   └── config/
│       └── env.go                     # Existing
├── migrations/
│   ├── 001-005                        # Existing Phase 1 migrations
│   ├── 006_create_tenants_table.*.sql ★
│   ├── 007_create_mandatory_fees_table.*.sql  ★
│   ├── 008_create_voluntary_fees_table.*.sql   ★
│   ├── 009_create_user_tenants_table.*.sql     ★
│   └── 010_seed_sample_tenants.*.sql           ★ (optional)
├── policy.csv                        # Update with tenant, fee resources
└── rbac_model.conf                   # Update matcher for resident tenant access
```

### Pattern 1: Clean Architecture CRUD (same pattern as Phase 1 auth)

**What:** Each domain concept gets Entity → Repository Interface → Service → Handler, following Phase 1's exact structure.

**Example — Entity:**
```go
// internal/domain/entity/tenant.go
package entity

import "time"

type Tenant struct {
    ID          string    `json:"id"`
    Block       string    `json:"block"`
    UnitNumber  string    `json:"unit_number"`
    Occupancy   string    `json:"occupancy"` // "occupied", "vacant"
    MonthlyFee  float64   `json:"monthly_fee"`
    TerritoryID string    `json:"territory_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Example — Repository Interface:**
```go
// internal/domain/repository/tenant_repository.go
type TenantRepository interface {
    Create(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error)
    FindByID(ctx context.Context, id string) (*entity.Tenant, error)
    ListByTerritory(ctx context.Context, territoryID string) ([]*entity.Tenant, error)
    Update(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error)
    Delete(ctx context.Context, id string, territoryID string) error
}
```

**Example — Service:**
```go
// internal/domain/service/tenant_service.go
type TenantService struct {
    tenantRepo repository.TenantRepository
    feeRepo    repository.FeeRepository
}

func (s *TenantService) Create(ctx context.Context, tenant *entity.Tenant, mandatoryFee *entity.MandatoryFee, userRole string, userTerritory string) (*entity.Tenant, error) {
    // 1. Validate uniqueness: (block, unit_number) within territory
    // 2. Validate mandatory fee exists (required per CONTEXT.md)
    // 3. Set territory_id from user context (RT officers) or request (RW officers)
    if userRole != "rw_officer" {
        tenant.TerritoryID = userTerritory
    }
    // 4. Create tenant + mandatory fee in transaction
    // 5. Return created tenant
}
```

**Source:** Pattern derived from Phase 1's `auth_service.go` (Create → Validate → Persist) and `user_repository.go` (pgx raw SQL).

### Anti-Patterns to Avoid

- **Putting territory filtering in the handler:** Territory filtering must be in the repository layer (like `ListByTerritory` in Phase 1). Handlers extract claims but don't construct SQL.
- **Using float64 for fee amounts:** Always use `shopspring/decimal` (or PostgreSQL `NUMERIC`) for monetary values. Float64 rounding errors cause financial discrepancies.
- **Mixing mandatory and voluntary fees:** CONTEXT.md mandates separate tables. Don't combine into a single `fees` table with a type column — the validation rules differ (mandatory requires ≥1 per tenant, voluntary is optional).
- **Modifying the existing Casbin enforcer singleton:** Phase 1 uses a `sync.Once` singleton pattern in `casbin.go`. Phase 2 additions to `policy.csv` must not break this. Test with `InitEnforcer` to validate policy loading.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Decimal/financial math | Custom decimal implementation | `shopspring/decimal` | Handles rounding, precision, serialization; battle-tested in production financial systems |
| RBAC enforcement | Custom permission checks | Casbin (already in project) | Phase 1 already invested in Casbin; territory-aware matcher is configured |
| Validation framework | Request validation in handlers | Service-layer validation (consistent with Phase 1) | Phase 1 does validation in service layer, not handlers |
| Transaction management | Manual BEGIN/COMMIT per operation | pgx `pool.Begin(ctx)` for multi-table operations | Creating a tenant + mandatory fee must be atomic; pgx transactions handle this |

**Key insight:** Phase 1 already made the "don't hand-roll" decisions for auth, RBAC, and database access. Phase 2 should follow the same patterns rather than introducing new libraries.

## Common Pitfalls

### Pitfall 1: Territory Isolation Leak (RT officers seeing other RT data)
**What goes wrong:** An RT-01 officer calls `GET /api/tenants` and gets RT-02 tenant records.
**Why it happens:** The repository query omits the `WHERE territory_id = $1` clause, or the handler doesn't pass the user's territory to the service.
**How to avoid:** Every tenant query in the repository must include territory filtering. The service layer receives the user's `territory_id` from claims (never from request body). Repository methods that accept a `territoryID` parameter enforce it in SQL.
**Warning signs:** Unit tests using mock repos return data from multiple territories; integration tests query with different role tokens.

### Pitfall 2: Resident Tenant Access Not Scoped
**What goes wrong:** A resident calls `GET /api/tenants` and gets ALL tenants in their territory instead of only their own.
**Why it happens:** The `user_tenants` junction table join is missing from the query, or the Casbin custom matcher isn't implemented.
**How to avoid:** Implement the `user_tenants` table and query it when the role is `resident`. Use a different repository method for residents (`ListByUserID` vs `ListByTerritory`). The Casbin custom matcher (decided in CONTEXT.md) must be implemented in `rbac_model.conf`.
**Warning signs:** No `user_tenants` migration exists; resident tests return more records than expected.

### Pitfall 3: Missing Mandatory Fee on Tenant Creation
**What goes wrong:** A tenant is created without a mandatory fee, violating CONTEXT.md's requirement.
**Why it happens:** The service layer creates the tenant but doesn't validate that at least one mandatory fee record was also created.
**How to avoid:** Use a pgx transaction for atomic create: insert tenant → validate mandatory fee exists in request → insert fee → commit. All in one transactional boundary. Roll back if mandatory fee fails.
**Warning signs:** Tenant creation accepts body without a `mandatory_fees` array.

### Pitfall 4: Date Validation Inconsistency
**What goes wrong:** A fee's `effective_date` is in the past, or `paid_at` is before `effective_date`.
**Why it happens:** Validation is done in the handler instead of the service layer, or timezone handling is inconsistent.
**How to avoid:** Centralize date validation in the service layer. Use UTC for all timestamps. Validate `effective_date >= today()` and `paid_at > effective_date` in the service before persisting.
**Warning signs:** Tests pass dates without timezone awareness; handler has `if time.Now()...` logic.

### Pitfall 5: Fee Amount Exceeds Monthly Fee Cap
**What goes wrong:** A mandatory fee amount exceeds the tenant's `monthly_fee` value.
**Why it happens:** No cross-entity validation between `mandatory_fees.amount` and `tenants.monthly_fee`.
**How to avoid:** In the tenant service's fee creation logic, load the tenant record and validate `fee.amount <= tenant.MonthlyFee`. Use decimal comparison (not float) to avoid precision issues.
**Warning signs:** Test creates tenant with monthly_fee=100000 and fee with amount=200000, expects success.

## Code Examples

### Pattern: Transactional Tenant + Fee Creation
```go
// Source: Derived from pgx/v5 documentation and Phase 1 pattern
// apps/api/internal/domain/service/tenant_service.go

func (s *TenantService) Create(ctx context.Context, req *CreateTenantRequest, userClaims *auth.Claims) (*entity.Tenant, error) {
    // Begin transaction
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx) // safe: no-op if committed

    // Determine territory scope
    territoryID := userClaims.TerritoryID
    if userClaims.Role == "rw_officer" && req.TerritoryID != "" {
        territoryID = req.TerritoryID
    }

    // Create tenant
    tenant := &entity.Tenant{
        Block:       req.Block,
        UnitNumber:  req.UnitNumber,
        Occupancy:   req.Occupancy,
        MonthlyFee:  req.MonthlyFee,
        TerritoryID: territoryID,
    }
    created, err := s.tenantRepo.CreateTx(ctx, tx, tenant)
    if err != nil {
        return nil, fmt.Errorf("failed to create tenant: %w", err)
    }

    // Validate at least one mandatory fee provided
    if len(req.MandatoryFees) == 0 {
        return nil, ErrMandatoryFeeRequired
    }

    // Create mandatory fees
    for _, feeReq := range req.MandatoryFees {
        fee := &entity.MandatoryFee{
            TenantID:      created.ID,
            Amount:        feeReq.Amount,
            Description:   feeReq.Description,
            EffectiveDate: feeReq.EffectiveDate,
        }
        if err := s.validateFee(fee, tenant.MonthlyFee); err != nil {
            return nil, fmt.Errorf("mandatory fee validation failed: %w", err)
        }
        if _, err := s.feeRepo.CreateMandatoryTx(ctx, tx, fee); err != nil {
            return nil, fmt.Errorf("failed to create mandatory fee: %w", err)
        }
    }

    // Commit
    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return created, nil
}
```

### Pattern: Territory-Filtered Repository Query
```go
// Source: Derived from PostgresUserRepository.FindByTerritory in Phase 1
// apps/api/internal/infrastructure/repository/tenant_repository.go

func (r *PostgresTenantRepository) ListByTerritory(ctx context.Context, territoryID string) ([]*entity.Tenant, error) {
    query := `
        SELECT id, block, unit_number, occupancy, monthly_fee, territory_id, created_at, updated_at
        FROM tenants
        WHERE territory_id = $1
        ORDER BY block, unit_number
    `
    rows, err := r.pool.Query(ctx, query, territoryID)
    if err != nil {
        return nil, fmt.Errorf("failed to list tenants by territory: %w", err)
    }
    defer rows.Close()

    var tenants []*entity.Tenant
    for rows.Next() {
        t := &entity.Tenant{}
        err := rows.Scan(&t.ID, &t.Block, &t.UnitNumber, &t.Occupancy, &t.MonthlyFee, &t.TerritoryID, &t.CreatedAt, &t.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan tenant row: %w", err)
        }
        tenants = append(tenants, t)
    }
    return tenants, rows.Err()
}
```

### Pattern: Nested Fee Routes with Fiber Groups
```go
// Source: Derived from protected_handler.go route registration pattern + Fiber routing docs
// apps/api/internal/delivery/http/tenant_handler.go

func (h *TenantHandler) RegisterRoutes(api fiber.Router, pasetoSvc *auth.PasetoService) {
    // Tenant CRUD
    api.Get("/tenants", h.ListTenants)
    api.Post("/tenants", h.CreateTenant)
    api.Get("/tenants/:id", h.GetTenant)
    api.Put("/tenants/:id", h.UpdateTenant)
    api.Delete("/tenants/:id", h.DeleteTenant)

    // Fee sub-resources (nested under tenant)
    api.Get("/tenants/:id/fees", h.ListFees)
    api.Post("/tenants/:id/fees", h.CreateFee)
    api.Put("/fees/:feeId", h.UpdateFee)
    api.Delete("/fees/:feeId", h.DeleteFee)
}
```

**Note:** The Casbin middleware's `defaultResourceExtractor` extracts the first path segment after `/api/`. For `/api/tenants/:id/fees`, it returns `"tenants"`. The policy must use `tenant` (singular) as the resource name to match the existing convention in `policy.csv` and the matcher.

### Pattern: Test Structure with Mock Repos
```go
// Source: Derived from auth_handler_test.go mock patterns in Phase 1
func TestTenantHandler_Create_Success(t *testing.T) {
    app := setupTestApp(t)

    body := map[string]interface{}{
        "block":      "A",
        "unit_number": "01",
        "occupancy":   "occupied",
        "monthly_fee": 50000,
        "mandatory_fees": []map[string]interface{}{
            {"amount": 25000, "description": "Security Fee", "effective_date": "2026-06-01"},
        },
    }
    jsonBody, _ := json.Marshal(body)

    req := httptest.NewRequest("POST", "/api/tenants", bytes.NewReader(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    // Set PASETO cookie for auth
    req.AddCookie(&http.Cookie{Name: "paseto_token", Value: validRTOfficerToken})

    resp, _ := app.Test(req)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Tenant stubs in `protected_handler.go` returning `{"tenants": []}` | Proper tenant handler with CRUD + territory filtering | Phase 2 | Replace stub responses with real DB queries |
| No fee tables exist | `mandatory_fees` + `voluntary_fees` tables | Phase 2 | New migrations 006-009 |
| Policy has `*` for all territories | Policy stays `*` (territory filtering at service/repo layer) | Continuation | No policy.csv schema change needed for basic CRUD |
| No resident tenant scoping | `user_tenants` junction table + Casbin custom matcher | Phase 2 | New migration 009 + rbac_model.conf update |

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | `shopspring/decimal` is the right library for fee amounts | Standard Stack | If Go ecosystem has a better financial decimal library, could add tech debt. However, shopspring has 7.4k GH stars and is standard in Go fintech. |
| A2 | Casbin policy should stay `*` wildcard for all territories | Architecture Patterns | CONTEXT.md says to use `{{territory_id}}` placeholders. Current code uses `*`. The `EnforceWithTerritory` method exists but isn't used by middleware. If CONTEXT.md decision is strictly enforced, `casbin.go` middleware and `policy.csv` need changes. |
| A3 | `tenant` (singular) resource name in Casbin policy is correct | Architecture Patterns | The `defaultResourceExtractor` would return `tenants` from `/api/tenants/...`, but `policy.csv` uses `tenant` (singular). The policy might need an update or the extractor needs adjustment. |

**All other claims verified against actual Phase 1 implementation files.**

## Open Questions

1. **`user_tenants` junction table schema?**
   - What we know: CONTEXT.md mandates a junction table linking users to tenants for resident access.
   - What's unclear: The exact schema — does it include relationship metadata (e.g., "owner", "occupant") or just user_id ↔ tenant_id mapping?
   - Recommendation: Start minimal: `user_tenants (user_id UUID REFERENCES users, tenant_id UUID REFERENCES tenants, PRIMARY KEY (user_id, tenant_id), created_at TIMESTAMP)`. Can be extended later.

2. **Cross-territory audit middleware?**
   - What we know: CONTEXT.md says "Auditing middleware logs any cross‑territory access attempts and returns 403."
   - What's unclear: Does this mean a new middleware layer, or should the existing Casbin middleware be extended? What logging destination (stdout, file, database)?
   - Recommendation: Extend the existing Casbin middleware in `apps/api/internal/delivery/middleware/casbin.go` to log denials via `log.Printf`. Database audit log can be added later if needed.

3. **How should the Casbin model handle `user_tenants`?**
   - What we know: CONTEXT.md says to use a custom matcher checking tenant ID against the user's tenant list.
   - What's unclear: Casbin's default matcher doesn't support dynamic lookups. Implementing this requires either: (a) passing tenant IDs as a parameter to `Enforce()`, or (b) implementing a custom Casbin function.
   - Recommendation: Option (a) — modify `CasbinMiddleware` to accept a `UserTenantIDs` function, and call `Enforce()` with the tenant list as an additional parameter. The matcher checks `r.tenant_id_in` against the parameter. This avoids modifying the Casbin model and keeps logic in the Go layer.

## Environment Availability

> Phase 2 extends existing code in `apps/api/` which already has all Go dependencies installed. No new external tools are required.

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go 1.26+ | Build & run | ✓ (from go.mod) | 1.26.2 | — |
| PostgreSQL 18+ | Data storage | ✓ | (assumed from env) | — |
| `make` | Build commands | ✓ | (system) | `go build ./cmd/server` |
| `golang-migrate` | Migration execution | ✓ | via `database.RunMigrations` | — |

## Security Domain

> `security_enforcement` is enabled (default — config.json does not set to false).

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | No | Handled by Phase 1 (PASETO + httpOnly cookies) |
| V3 Session Management | No | Handled by Phase 1 (PASETO token lifecycle) |
| V4 Access Control | **Yes** | Casbin RBAC with territory-aware matcher; service-layer `WHERE territory_id = $1` |
| V5 Input Validation | **Yes** | Service-layer validation of fee amounts, dates, block/unit uniqueness |
| V6 Cryptography | No | No cryptographic operations in Phase 2 (PASETO handled by Phase 1) |
| V8 Data Protection | **Yes** | Territory isolation prevents RT-01 officers from viewing RT-02 tenant/fee data |

### Known Threat Patterns for Go + Fiber + PostgreSQL

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Horizontal privilege escalation (RT-01 views RT-02 tenants) | Information Disclosure | Repository-layer `WHERE territory_id = $1` ensures user can only query their assigned territory |
| Vertical privilege escalation (resident creates tenants) | Elevation of Privilege | Casbin middleware blocks `write` action on `tenant` resource for `resident` role |
| Fee manipulation (amount > monthly_fee cap) | Tampering | Service-layer validation compares `fee.amount` against `tenant.monthly_fee` |
| SQL injection via tenant fields | Tampering | pgx parameterized queries (`$1`, `$2`) prevent injection — all Phase 1 queries use this pattern |
| Insecure direct object reference (IDOR) on fee records | Information Disclosure | Fee queries filter by `tenant_id` which is scoped to the user's territory via the tenant join |

## Validation Architecture

> `workflow.nyquist_validation` is enabled in `.planning/config.json`.

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go testing + `github.com/stretchr/testify` |
| Config file | none (standard Go convention: `*_test.go` files) |
| Quick run command | `cd apps/api && go test -v -race -count=1 ./internal/domain/service/ ./internal/infrastructure/repository/` |
| Full suite command | `cd apps/api && go test -v -race -cover ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| TEN-01 | Create tenant with valid data | unit + integration | `go test ./internal/domain/service/ -run TestTenantService_Create_Success -v` | ❌ Wave 0 |
| TEN-01 | Tenant list respects territory isolation | unit | `go test ./internal/infrastructure/repository/ -run TestTenantRepository_ListByTerritory -v` | ❌ Wave 0 |
| TEN-01 | Tenant uniqueness (block, unit) within territory | unit | `go test ./internal/domain/service/ -run TestTenantService_Create_DuplicateBlockUnit -v` | ❌ Wave 0 |
| TEN-01 | Delete tenant enforces territory scope | unit | `go test ./internal/domain/service/ -run TestTenantService_Delete_NotOwnTerritory -v` | ❌ Wave 0 |
| FIN-01 | Create tenant fails without mandatory fee | unit | `go test ./internal/domain/service/ -run TestTenantService_Create_MandatoryFeeRequired -v` | ❌ Wave 0 |
| FIN-01 | Mandatory fee amount ≤ tenant monthly_fee | unit | `go test ./internal/domain/service/ -run TestTenantService_Create_FeeExceedsMonthlyCap -v` | ❌ Wave 0 |
| FIN-02 | Create voluntary fee with valid dates | unit | `go test ./internal/domain/service/ -run TestFeeService_CreateVoluntary_ValidDates -v` | ❌ Wave 0 |
| FIN-02 | Voluntary fee dates: effective_date not in past | unit | `go test ./internal/domain/service/ -run TestFeeService_CreateVoluntary_PastEffectiveDate -v` | ❌ Wave 0 |
| FIN-02 | Voluntary fee dates: paid_at after effective_date | unit | `go test ./internal/domain/service/ -run TestFeeService_CreateVoluntary_PaidAtBeforeEffective -v` | ❌ Wave 0 |
| AUTH-02 | Resident reads only own tenants | unit | `go test ./internal/domain/service/ -run TestTenantService_List_ResidentOnlyOwn -v` | ❌ Wave 0 |
| AUTH-02 | RT officer reads only own territory tenants | unit | `go test ./internal/domain/service/ -run TestTenantService_List_RTOfficerScoped -v` | ❌ Wave 0 |
| AUTH-02 | RW officer reads all territories | unit | `go test ./internal/domain/service/ -run TestTenantService_List_RWOfficerAll -v` | ❌ Wave 0 |

### Sampling Rate

- **Per task commit:** `cd apps/api && go test -v -race -count=1 ./internal/domain/service/ -run TestTenantService`
- **Per wave merge:** `cd apps/api && go test -v -race -cover ./...`
- **Phase gate:** Full suite green before `/gsd-verify-work`

### Wave 0 Gaps

- [ ] `apps/api/internal/domain/service/tenant_service_test.go` — covers all TEN-01 and FIN-01/02 unit tests
- [ ] `apps/api/internal/infrastructure/repository/tenant_repository_test.go` — integration tests with mock DB
- [ ] `apps/api/internal/delivery/http/tenant_handler_test.go` — httptest-based handler tests
- [ ] `apps/api/internal/domain/service/fee_service_test.go` — fee-specific validation tests
- [ ] `apps/api/internal/delivery/middleware/audit_test.go` — cross-territory audit logging tests

## Sources

### Primary (HIGH confidence)
- Phase 1 implementation code in `apps/api/` — all architecture patterns, repository structure, handler patterns, Casbin model, migration structure
- `.planning/phases/02-tenant-fee-management/01-CONTEXT.md` — locked decisions, API design, constraints
- `go.mod` — exact versions of all dependencies
- `pkg.go.dev/github.com/shopspring/decimal` — verified decimal library API and stability

### Secondary (MEDIUM confidence)
- `docs.gofiber.io/guide/routing/` — Fiber routing patterns for nested paths
- `docs.gofiber.io/guide/grouping/` — Fiber group behavior (routes flattened, prefix-based)

### Tertiary (LOW confidence)
- None — all research claims verified against running codebase or official docs

## Metadata

**Confidence breakdown:**
- Standard stack: **HIGH** — all libraries verified against existing `go.mod` and actual implementation code
- Architecture: **HIGH** — patterns derived directly from Phase 1 code, not speculative
- Pitfalls: **HIGH** — based on understanding of the existing Casbin + repository + service layer interaction
- Package legitimacy: **HIGH** — all existing packages from Phase 1 verified; new package (`shopspring/decimal`) is 10-year-old mature library

**Research date:** 2026-05-22
**Valid until:** 2026-06-22 (30-day window for stable Go ecosystem)
