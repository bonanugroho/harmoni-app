# Phase 1: Core Authentication & RBAC - Context

**Gathered:** 2026-05-19
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers secure user registration, login, password reset, and role‑based access control (RBAC) for Residents, RT Officers, and RW Officers, with territory‑aware data isolation.
</domain>

<decisions>
## Implementation Decisions

### Session & Token Storage
- **D-01:** Use httpOnly cookies (server‑managed, XSS‑protected) to store PASETO tokens.

### Password Reset Flow
- **D-02:** Implement email‑based password reset links containing one‑time tokens.

### Territory Model
- **D-03:** Associate each user with a single `territory_id` (simplest, clear ownership).

### Role Assignment
- **D-04:** Provide an admin dashboard for creating and assigning RT/RW Officer accounts.

### Password Policy
- **D-05:** Enforce complexity rules requiring uppercase, lowercase, numbers, and symbols.

### Casbin Policy Structure
- **D-06:** Hybrid approach with resource-based policies and `{{territory_id}}` placeholders for RT officers, `*` wildcard for RW officers.

### Database Migration Strategy
- **D-07:** Use golang-migrate CLI (SQL-only, versioned, up/down files).
- **D-08:** All UUID fields use PostgreSQL 18's native `uuidv7()` function as DEFAULT — no application-layer generation needed.
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Scope
- `.planning/PROJECT.md` — project overview, core value, constraints, and key decisions.
- `.planning/REQUIREMENTS.md` — listed v1 requirements, including AUTH‑01 and AUTH‑02 mapped to Phase 1.
- `.planning/ROADMAP.md` — Phase 1 goal and success criteria.

### Technical Foundations
- `.planning/codebase/STACK.md` — current technology stack (JavaScript, Node.js, Go not yet present).
- `.planning/codebase/ARCHITECTURE.md` — placeholder architecture overview and future layer expectations.
- `.planning/codebase/INTEGRATIONS.md` — potential integration points (database, auth services).
</canonical_refs>

<code_context>
## Existing Code Insights

*No application source code exists yet; only tooling and planning artifacts are present.*
</code_context>

<specifics>
## Specific Ideas

No additional specific references were provided beyond the decisions above.
</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.
</deferred>

---

*Phase: 1-Core Authentication & RBAC*
*Context gathered: 2026-05-19*