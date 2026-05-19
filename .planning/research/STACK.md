---
layout: doc
date: 2026-05-19
---

# Technology Stack

| Component | Technology |
|-----------|------------|
| **Frontend** | React.js (Vite) + Tailwind CSS |
| **Backend** | Go (Fiber framework) |
| **Database** | PostgreSQL |
| **Security** | PASETO (tokenisation) + Casbin (RBAC) |
| **Cache** | Redis (optional) |
| **Architecture** | Clean Architecture (domain‑driven layers) |

**Rationale**
- React + Tailwind gives a fast, mobile‑first UI with low bundle size.
- Go + Fiber provides high‑throughput HTTP handling and easy concurrency for financial transactions.
- PostgreSQL offers strong ACID guarantees critical for monetary data.
- PASETO V4 Local encrypts token payloads, avoiding the pitfalls of JWT replay attacks.
- Casbin supplies fine‑grained, territory‑aware RBAC policies.
- Clean Architecture isolates business rules from framework code, facilitating testing and future migrations.
