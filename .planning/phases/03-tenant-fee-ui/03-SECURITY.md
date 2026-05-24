---
phase: 3
slug: tenant-fee-ui
status: verified
threats_open: 0
asvs_level: 1
created: 2026-05-24
---

# Phase 3 — Security

> Per-phase security contract: threat register, accepted risks, and audit trail.

---

## Trust Boundaries

| Boundary | Description | Data Crossing |
|----------|-------------|---------------|
| Browser ↔ API | All data flows through `request()` helper with `credentials: 'include'` | Tenant/fee data (block, unit_number, monthly_fee, fee amounts, effective dates) — low sensitivity |
| UI Components | Client-side data display and form validation only | No sensitive data; types match API response shape |

---

## Threat Register

| Threat ID | Category | Component | Disposition | Mitigation | Status |
|-----------|----------|-----------|-------------|------------|--------|
| T-03-01 | Tampering | `request()` helper | mitigate | `credentials: 'include'` — PASETO cookie sent same-origin only. No manual token handling. | closed |
| T-03-02 | Information Disclosure | Tenant/Fee types | accept | Types match API response shape. No PII beyond block/unit_number (public per-domain data). | closed |
| T-03-03 | Tampering | npm packages | mitigate | @tanstack/react-query and lucide-react verified legitimate via package.json declaration. | closed |
| T-03-SC | Tampering | @tanstack/react-query install | mitigate | Dependency declared in package.json, npm integrity verified. | closed |
| T-03-04 | Information Disclosure | AppLayout sidebar | mitigate | Settings link gated by role check (`rt_officer`/`rw_officer`) via useAuth(). | closed |
| T-03-05 | Tampering | ConfirmDialog | mitigate | Controlled component — `isOpen`, `onConfirm`, `onCancel` parent-controlled. | closed |
| T-03-06 | Elevation of Privilege | Route visibility | accept | Sidebar only hides links visually; ProtectedRoute + backend Casbin enforce server-side. | closed |
| T-03-07 | Information Disclosure | AuthContext user | accept | User object contains role/territory for UI decisions. No password, no token. | closed |
| T-03-08 | Elevation of Privilege | ProtectedRoute | accept | Client-side role check is UX convenience only; backend Casbin enforces. | closed |
| T-03-09 | Tampering | Cache invalidation | mitigate | All 6 mutation hooks call `invalidateQueries` in `onSuccess`. | closed |
| T-03-10 | Denial of Service | Excessive query refetch | mitigate | `staleTime: 30_000` in QueryClient defaults prevents on-keystroke refetch. | closed |
| T-03-11 | Tampering | TenantForm validation | mitigate | Client-side `validate()` checks block, unit_number, monthly_fee, mandatory fees before submit. | closed |
| T-03-12 | Elevation of Privilege | TenantEditPage delete | accept | Delete button visible but backend Casbin prevents unauthorized deletes. | closed |
| T-03-13 | Information Disclosure | TenantCard search | accept | Search filters client-side only; backend enforces territory isolation. | closed |
| T-03-14 | Tampering | FeeForm validation | mitigate | Validates amount > 0, ≤ monthlyFee, effective date not in past. | closed |
| T-03-15 | Denial of Service | Excessive mutations | accept | No rate limiting on frontend — backend should handle. Low risk for neighborhood-scale app. | closed |
| T-03-16 | Tampering | ConfirmDialog delete flow | mitigate | Two-step confirmation via ConfirmDialog; UI updates driven by cache invalidation (not optimistic). | closed |

*Status: open · closed*
*Disposition: mitigate (implementation required) · accept (documented risk) · transfer (third-party)*

---

## Accepted Risks Log

| Risk ID | Threat Ref | Rationale | Accepted By | Date |
|---------|------------|-----------|-------------|------|
| AR-01 | T-03-02 | Types match API response — no PII beyond public per-domain data (block/unit_number) | Planning | 2026-05-24 |
| AR-02 | T-03-06 | Sidebar link visibility is UX-only; route protection and Casbin enforce server-side | Planning | 2026-05-24 |
| AR-03 | T-03-07 | AuthContext exposes only role/territory — no passwords or tokens | Planning | 2026-05-24 |
| AR-04 | T-03-08 | Client role check is UX-only; backend Casbin enforces authorization | Planning | 2026-05-24 |
| AR-05 | T-03-12 | Delete button visible to authorized roles; backend Casbin prevents unauthorized deletes | Planning | 2026-05-24 |
| AR-06 | T-03-13 | Client-side search is convenience; backend enforces territory data isolation | Planning | 2026-05-24 |
| AR-07 | T-03-15 | No frontend rate limiting — acceptable for neighborhood-scale app; backend may add | Planning | 2026-05-24 |

---

## Security Audit Trail

| Audit Date | Threats Total | Closed | Open | Run By |
|------------|---------------|--------|------|--------|
| 2026-05-24 | 17 | 17 | 0 | gsd-security-auditor |

---

## Sign-Off

- [x] All threats have a disposition (mitigate / accept / transfer)
- [x] Accepted risks documented in Accepted Risks Log
- [x] `threats_open: 0` confirmed
- [x] `status: verified` set in frontmatter

**Approval:** verified 2026-05-24
