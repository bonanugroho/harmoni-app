---
layout: doc
date: 2026-05-19
---

# Roadmap

## Phase 1: Core Authentication & RBAC

**Goal:** Implement secure user registration, login, password reset, and role‑based access control for Residents, RT Officers, and RW Officers.
**Success Criteria:**

1. Users can create accounts, log in, and reset passwords.
2. JWT is replaced by PASETO tokens; tokens are validated on every request.
3. Casbin policies enforce correct permissions per territory.
4. Automated unit & integration tests cover all auth flows.

## Phase 2: Tenant & Fee Management

**Goal:** Provide CRUD operations for tenant data and record mandatory/voluntary fees.
**Success Criteria:**

1. Tenant records include block, unit number, occupancy, and monthly fee.
2. Mandatory fees (e.g., waste, security) can be defined per tenant.
3. Voluntary contributions can be added and reported.
4. Data isolation ensures RT 01 officers cannot access RT 02 tenant data.

**Plans:** 4/4 plans complete

- [x] 02-01-PLAN.md — Database migrations & entity definitions (TEN-01, FIN-01, FIN-02)
- [x] 02-02-PLAN.md — Repository interfaces & pgx implementations (TEN-01, FIN-01, FIN-02)
- [x] 02-03-PLAN.md — Service layer with validation & policy updates (TEN-01, FIN-01, FIN-02)
- [x] 02-04-PLAN.md — HTTP handler, main.go wiring & stub removal (TEN-01, FIN-01, FIN-02)

## Phase 3: Tenant & Fee UI

**Goal:** Build the tenant and fee management frontend based on the 02-UI-SPEC.md design contract — responsive, mobile-first UI for RT/RW officers and residents.
**Success Criteria:**

1. Tenant list page with card layout, search/filter, and responsive grid (1→2→3 columns).
2. Create/edit tenant form with dynamic mandatory fee section.
3. Tenant detail page with fee management (mandatory + voluntary sections).
4. AppLayout with responsive sidebar (collapsible on mobile, fixed on desktop).
5. All pages work on low-end mobile browsers with 44px touch targets.
6. Unit tests cover all new components and pages.

**Plans:** 4/5 plans complete

- [x] 03-01-PLAN.md — Foundation: api helper, types, services, QueryClientProvider
- [x] 03-02-PLAN.md — UI components + AppLayout
- [x] 03-03-PLAN.md — TanStack Query hooks + ProtectedRoute migration
- [x] 03-04-PLAN.md — Tenant pages (list, create, edit)
- [ ] 03-05-PLAN.md — Fee pages (detail, fee management)

## Phase 4: Transaction Engine & Expenditures

**Goal:** Record income (fees, contributions, RT→RW transfers) and expenses (operational costs).
**Success Criteria:**

1. Income entries are stored with audit timestamps.
2. Expenditure entries can be added, categorized, and linked to supporting documents.
3. RT → RW transfers update balances correctly and are reflected in reports.
4. End‑to‑end testing validates financial calculations.

## Phase 5: Dashboards & Reporting

**Goal:** Deliver real‑time financial dashboards and accounts‑receivable analysis.
**Success Criteria:**

1. Dashboard shows current cash balance, income vs. expense trends.
2. Historical charts are interactive and responsive on mobile devices.
3. AR analysis categorises overdue payments (on‑time, >30d, >60d, >90d).
4. Performance benchmarks meet <200 ms response time for dashboard endpoints.

---

*Last updated: 2026‑05‑24*
