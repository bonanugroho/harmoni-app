---
layout: doc
date: 2026-05-19
---

# Architecture Overview

The backend follows **Clean Architecture** principles, separating concerns into distinct layers:

- **Domain Layer (`internal/domain`)** – Core business entities (User, Tenant, Transaction, Report) and repository interfaces.
- **Use‑Case / Application Layer (`internal/app`)** – Business logic orchestrating domain entities (e.g., `CreateTenant`, `RecordIncome`, `GenerateReport`).
- **Infrastructure Layer (`internal/infrastructure`)** – Implementations for persistence (PostgreSQL), caching (Redis), security (PASETO generation/validation, Casbin policy loading).
- **Interface Layer (`internal/interface`)** – HTTP handlers built with Fiber, request/response DTOs, and middleware (auth, logging, RBAC enforcement).

**Data Flow**
1. HTTP request → Fiber handler (Interface).
2. Handler extracts authentication token, validates via PASETO, enforces Casbin RBAC.
3. Calls appropriate use‑case in Application layer.
4. Use‑case interacts with Domain entities and repository interfaces.
5. Infrastructure provides concrete repository implementations (PostgreSQL) and auxiliary services (Redis cache).

**Deployment**
- Monorepo layout with separate directories for the web frontend (`/apps/web`) and API backend (`/apps/api`).
- API compiled as a single binary (`go build`) and served behind a reverse proxy (e.g., Nginx) in production.
- Database migrations managed via a Go migration tool (e.g., golang‑migrate).
