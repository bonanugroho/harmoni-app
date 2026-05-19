# Specification Document - Harmoni App Project

## 1. Tech Stack
| Component | Technology |
| :--- | :--- |
| **Frontend** | React.js (Vite), Tailwind CSS |
| **Backend** | Golang (Fiber Framework) |
| **Database** | PostgreSQL |
| **Security** | PASETO (Tokenization), Casbin (RBAC Engine) |
| **Cache** | Redis (Optional/As needed) |

## 2. Backend Architecture (Clean Architecture)
The project structure follows the principle of decoupling business logic from infrastructure:
- `internal/domain`: Core entities and repository interfaces.
- `internal/app`: Business logic and use cases (e.g., arrears calculation).
- `internal/infrastructure`: Implementations for Database, Casbin, and PASETO.
- `internal/interface`: Fiber HTTP Handlers and Middlewares.

## 3. Security & Authorization
- **Tokenization:** PASETO V4 Local for secure, encrypted server-side tokens.
- **RBAC Policy (Casbin):** Utilizing a "Domain" model to restrict access based on territory IDs.
  - Matcher: `g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act`

## 5. Development & Deployment
- **Repository Model:** Monorepo (Path: `/apps/web` and `/apps/api`).
- **Communication:** RESTful API using JSON format.
- **UI/UX:** Mobile-first approach using Tailwind CSS utility classes for high performance on mobile devices.