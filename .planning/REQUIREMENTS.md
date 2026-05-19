---
layout: doc
date: 2026-05-19
---

# Requirements: Harmoni

**Defined:** 2026-05-19
**Core Value:** Transparency and accountability of community finances

## v1 Requirements

### Authentication & Access Control
- [ ] **AUTH-01**: User can register, log in, and reset password.
- [ ] **AUTH-02**: Role‑Based Access Control (Resident, RT Officer, RW Officer) enforced via Casbin policies.

### Tenant Management
- [ ] **TEN-01**: Record tenant information (house block, unit number, occupancy status, monthly fee).

### Income & Fee Management
- [ ] **FIN-01**: Record mandatory fees (fixed monthly fees per unit).
- [ ] **FIN-02**: Record voluntary contributions (e.g., holiday bonuses, social donations).
- [ ] **FIN-03**: Handle RT → RW cash transfers.

### Expenditure Tracking
- [ ] **EXP-01**: Record operational costs (security salaries, cleaning services, public facility maintenance).

### Reporting & Dashboard
- [ ] **REP-01**: Real‑time dashboard displaying cash balances.
- [ ] **REP-02**: Visualise historical income vs. expenditure trends.
- [ ] **REP-03**: Accounts‑receivable analysis per RT (on‑time, late > 30 d, > 60 d, > 90 d).

## v2 Requirements (deferred)

### Advanced Features
- **ADV-01**: Mobile native app (iOS/Android).
- **ADV-02**: Real‑time chat between residents and officers.
- **ADV-03**: Advanced analytics dashboards with predictive insights.

## Out of Scope

| Feature | Reason |
|---------|--------|
| Mobile native app | Web‑first strategy; native app adds significant scope.
| Real‑time chat | Not core to financial management; deferred to v2.
| Advanced analytics | Requires data‑science pipeline; out of scope for MVP.

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| AUTH-01 | Phase 1 | Pending |
| AUTH-02 | Phase 1 | Pending |
| TEN-01 | Phase 2 | Pending |
| FIN-01 | Phase 2 | Pending |
| FIN-02 | Phase 2 | Pending |
| FIN-03 | Phase 3 | Pending |
| EXP-01 | Phase 3 | Pending |
| REP-01 | Phase 4 | Pending |
| REP-02 | Phase 4 | Pending |
| REP-03 | Phase 4 | Pending |

**Coverage:**
- v1 requirements: 10 total
- Mapped to phases: 10
- Unmapped: 0 ✓

---
*Requirements defined: 2026‑05‑19*
*Last updated: 2026‑05‑19 after initialization*
