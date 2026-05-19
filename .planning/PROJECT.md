---
layout: doc
date: 2026-05-19
---

# Harmoni

## What This Is

Harmoni is a community financial‑management web application for neighborhood‑scale administrations (Rukun Tetangga/RT and Rukun Warga/RW). It provides transparent, accountable income and expenditure reporting for Residents, RT Officers, and RW Officers, accessible via mobile‑first browsers.

## Core Value

Transparency and accountability of community finances – if the reporting layer fails, the whole system loses trust.

## Requirements

### Validated

*(None yet – will be validated as the app is shipped.)*

### Active

- [ ] **AUTH‑01**: User can register, log in, and reset password.
- [ ] **AUTH‑02**: Role‑Based Access Control (Resident, RT Officer, RW Officer).
- [ ] **TEN‑01**: Record tenant information (house block, number, occupancy, monthly fee).
- [ ] **FIN‑01**: Capture mandatory fees, voluntary contributions, and RT→RW transfers.
- [ ] **EXP‑01**: Record operational expenditures (security salaries, cleaning, maintenance).
- [ ] **REP‑01**: Real‑time dashboard showing cash balances and income vs. expenditure trends.
- [ ] **REP‑02**: Accounts‑receivable analysis per RT (on‑time, late > 30 d, > 60 d, > 90 d).

### Out of Scope

- Mobile native app (web‑first only).
- Real‑time chat or messaging between users.
- Advanced analytics beyond basic trend graphs.

## Context

- **Frontend**: React + Vite, Tailwind CSS (mobile‑first).
- **Backend**: Go (Fiber framework), Clean Architecture layout.
- **Database**: PostgreSQL.
- **Security**: PASETO V4 Local tokens, Casbin RBAC engine.
- **Cache**: Optional Redis.
- **Deployment**: Monorepo with `/apps/web` and `/apps/api`.

## Constraints

- **Security**: Must use PASETO tokens and enforce Casbin policies per territory.
- **Data Isolation**: RT 01 officers cannot view RT 02 data.
- **Responsiveness**: UI must work on low‑end mobile browsers.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Tech Stack (React + Go) | Leverages existing team expertise; Go provides strong concurrency for transaction processing. | ✅ Adopted |
| PASETO over JWT | Provides built‑in encryption and mitigates JWT replay attacks. | ✅ Adopted |
| Clean Architecture layout | Keeps business logic testable and infrastructure‑agnostic. | ✅ Adopted |

---
*Last updated: 2026‑05‑19 after initialization*
