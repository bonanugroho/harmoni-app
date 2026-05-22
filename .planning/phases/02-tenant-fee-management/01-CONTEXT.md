# Phase 2: Tenant & Fee Management – Context

**Gathered:** 2026-05-22
**Status:** Planning

<domain>
## Phase Boundary
This phase delivers CRUD operations for tenant records (block, unit number, occupancy, monthly fee) and fee management (mandatory fees per tenant, voluntary contributions). Data isolation ensures RT 01 officers cannot view RT 02 tenant data.
</domain>

<decisions>
## Implementation Decisions (locked)

### Data Isolation
- Use Casbin policies with `{{territory_id}}` placeholders for RT officers, and `*` wildcard for RW officers.
- Enforce at the service layer: each tenant query must include `WHERE territory_id = {{user.territory_id}}` for RT roles.
- Auditing middleware logs any cross‑territory access attempts and returns `403`.

### Fee Types
- **Mandatory fees** are defined per‑tenant in a `mandatory_fees` table linked to `tenants`.
- **Voluntary contributions** are stored in a separate `voluntary_fees` table, allowing multiple entries per tenant.
- Both fee tables include `amount`, `description`, `effective_date`, and `paid_at` fields.

### API Design
- **Tenants:** `GET /api/tenants`, `POST /api/tenants`, `GET /api/tenants/:id`, `PUT /api/tenants/:id`, `DELETE /api/tenants/:id`.
- **Fees:** `GET /api/tenants/:id/fees`, `POST /api/tenants/:id/fees` (both mandatory & voluntary), `PUT /api/fees/:feeId`, `DELETE /api/fees/:feeId`.
- All endpoints are protected by the Casbin middleware; RW officers have `*` access, RT officers are scoped by `territory_id`.

### Resident access
- Residents can read **only** the tenant records that they own (one or many). The Casbin policy uses a placeholder `{{tenant_id}}` and a custom matcher that checks the requested tenant ID against the list of tenant IDs associated with the user via the `user_tenants` junction table.
- The service layer fetches the set of tenant IDs for the user once per request and passes it to the enforcer, guaranteeing consistent enforcement across all tenant‑related endpoints.

- **Tenant Uniqueness:** (`block`, `unit_number`) must be unique within a territory.
- **Fee Amounts:** Non‑negative decimal, must not exceed tenant’s monthly fee cap (configurable).
- **Mandatory Fee Presence:** Every tenant must have at least one mandatory fee record; creation fails otherwise.
- **Voluntary Fee Optionality:** No required count; can be empty.
- **Date Fields:** `effective_date` cannot be in the past; `paid_at` must be after `effective_date`.

</decisions>

<canonical_refs>
## References
- `.planning/ROADMAP.md` – Phase 2 goal and success criteria.
- `.planning/REQUIREMENTS.md` – Tenant & Fee Management requirement IDs (e.g., TENANT‑01, FEE‑01).
- `.planning/codebase/STACK.md` – Current tech stack (Go + Fiber, PostgreSQL).
</canonical_refs>

<deferred>
## Deferred Ideas
- Multi‑currency support for fees.
- Historical audit of fee changes beyond current month.
- Bulk import/export of tenant data.
</deferred>

---

*Phase: 2‑Tenant & Fee Management*
*Context gathered: 2026-05-22*