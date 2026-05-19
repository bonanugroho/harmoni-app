---
layout: doc
date: 2026-05-19
---

# Project Pitfalls & Technical Debt

## Security
- **PASETO key rotation** – no process defined yet; must implement regular key rollover to avoid token reuse.
- **Casbin policy management** – policies are stored in a YAML file; editing manually can introduce misconfigurations that break RBAC.

## Data Isolation
- Territory‑level isolation relies on correct `domain` values in Casbin policies; any mismatch will leak data across RTs.
- No database row‑level security enforced at the PostgreSQL level – all isolation is application‑level, which is a risk if a bug bypasses Casbin.

## Performance
- Redis is optional; without caching, high‑traffic dashboards may suffer latency.
- Fiber’s default settings are tuned for moderate load; large‑scale deployments may need connection‑pool tuning.

## Maintainability
- Clean Architecture adds indirection; developers unfamiliar with the pattern may find the folder structure confusing.
- Go modules are currently pinned to specific versions; upgrades will require manual vetting.

## UI/UX
- Tailwind‑based mobile‑first design can lead to class‑name bloat; consider using a component library to enforce consistency.
- No offline support or progressive‑web‑app features – users on unreliable networks may have a degraded experience.

## Recommendations
- Define a key‑rotation schedule and store PASETO keys in a secret manager.
- Add automated tests for Casbin policy loading and enforce policy linting.
- Implement row‑level security or a database view layer for an additional safety net.
- Benchmark dashboard endpoints and add Redis caching where latency is > 200 ms.
- Draft a developer guide explaining the Clean Architecture layout.
