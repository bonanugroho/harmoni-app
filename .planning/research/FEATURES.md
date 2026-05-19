---
layout: doc
date: 2026-05-19
---

# Feature Summary

## User Management & Access
- Authentication (login, registration, password reset).
- Role‑Based Access Control (Resident, RT Officer, RW Officer) via Casbin policies.

## Tenant & Fee Management
- Store tenant information (house block, unit number, occupancy, monthly fee).
- Record mandatory fees (e.g., waste management, security) and voluntary contributions (e.g., holiday bonuses, social donations).
- Support RT → RW cash transfers.

## Transaction Recording
- Capture income entries and associate them with tenants.
- Capture expense entries (operational costs, salaries, maintenance).

## Reporting & Dashboard
- Real‑time cash balance display.
- Historical income vs. expenditure trend visualisations.
- Accounts‑receivable analysis per RT (on‑time, late > 30 d, > 60 d, > 90 d).

## Non‑Functional
- Mobile‑first responsive UI (React + Tailwind).
- Secure token handling with PASETO V4 Local.
- Strict data isolation per territory (RBAC enforcement).
