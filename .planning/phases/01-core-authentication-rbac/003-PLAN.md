---
wave: 3
depends_on:
  - 001
  - 002
files_modified:
  - apps/api/internal/infrastructure/auth/casbin.go
  - apps/api/internal/delivery/middleware/auth.go
  - apps/api/internal/delivery/middleware/casbin.go
  - apps/api/policy.csv
autonomous: true
requirements:
  - AUTH-02
---

# Plan 3: RBAC Middleware & Casbin Policies

## Objective
Implement Casbin RBAC engine with territory-aware policy enforcement and authentication middleware.

## Tasks

### Task 1: Casbin Policy Configuration
<read_first>
- apps/api/policy.csv (create)
- apps/api/rbac_model.conf (create)
- .planning/phases/01-core-authentication-rbac/01-CONTEXT.md (Casbin hybrid approach)
- .planning/phases/01-core-authentication-rbac/01-RESEARCH.md §2 (Casbin integration)
</read_first>

<action>
Create Casbin model and policy files:
- rbac_model.conf: [request_definition], [policy_definition], [role_definition], [policy_effect], [matchers]
- policy.csv: Hybrid policies with {{territory_id}} for RT officers, * for RW officers
- Roles: resident, rt_officer, rw_officer
- Resources: tenant, income, expenditure, report
- Actions: read, write
- Role hierarchy: g, rw_officer, rt_officer (RW inherits RT permissions)
</action>

<acceptance_criteria>
- apps/api/rbac_model.conf defines RBAC model with territory domain
- apps/api/policy.csv contains policies for all 3 roles
- RT officer policies use {{territory_id}} placeholder
- RW officer policies use * wildcard for territory
- Resident policies are read-only for own territory
- `go test ./internal/infrastructure/auth -v -run TestCasbinPolicy` passes
- Policy enforces: rt_officer can read/write rt-01, rw_officer can read/write rt-01 AND rt-02
</acceptance_criteria>

---

### Task 2: Casbin Enforcer Initialization
<read_first>
- apps/api/internal/infrastructure/auth/casbin.go (create)
- apps/api/rbac_model.conf
- apps/api/policy.csv
</read_first>

<action>
Initialize Casbin enforcer:
- Load model from rbac_model.conf
- Load policy from policy.csv (file adapter for MVP, DB adapter later)
- Implement Enforce(role, resource, action, territory) → bool
- Implement AddPolicy, RemovePolicy for runtime policy management
- Cache enforcer instance (singleton pattern)
</action>

<acceptance_criteria>
- apps/api/internal/infrastructure/auth/casbin.go initializes Casbin enforcer
- Enforce returns true for allowed role/resource/action/territory combinations
- Enforce returns false for disallowed combinations
- File adapter loads policy.csv correctly on startup
- `go test ./internal/infrastructure/auth -v -run TestCasbinEnforcer` passes
</acceptance_criteria>

---

### Task 3: Authentication Middleware
<read_first>
- apps/api/internal/delivery/middleware/auth.go (create)
- apps/api/internal/infrastructure/auth/paseto.go
- apps/api/internal/config/env.go
</read_first>

<action>
Implement authentication middleware:
- Extract paseto_token from httpOnly cookie
- Validate token using PASETO service
- Set user claims in request context: user_id, role, territory_id
- Return 401 if token missing, invalid, or expired
- Skip authentication for public routes: /health, /auth/register, /auth/login, /auth/reset
</action>

<acceptance_criteria>
- apps/api/internal/delivery/middleware/auth.go implements Fiber middleware
- Valid token sets user claims in c.Locals("user")
- Missing token returns 401 {"error": "Unauthorized", "code": "MISSING_TOKEN"}
- Invalid token returns 401 {"error": "Unauthorized", "code": "INVALID_TOKEN"}
- Expired token returns 401 {"error": "Unauthorized", "code": "TOKEN_EXPIRED"}
- Public routes (/health, /auth/*) bypass authentication
- `go test ./internal/delivery/middleware -v -run TestAuthMiddleware` passes
</acceptance_criteria>

---

### Task 4: Casbin Authorization Middleware
<read_first>
- apps/api/internal/delivery/middleware/casbin.go (create)
- apps/api/internal/infrastructure/auth/casbin.go
- apps/api/internal/delivery/middleware/auth.go
</read_first>

<action>
Implement Casbin authorization middleware:
- Extract user claims from c.Locals("user")
- Map HTTP method to action: GET→read, POST/PUT/PATCH→write, DELETE→write
- Extract resource from route or request context
- Substitute {{territory_id}} with user's territory_id for RT officers
- Use * for RW officers (automatic access to all territories)
- Return 403 if Casbin.Enforce returns false
- Chain with auth middleware: auth → casbin → handler
</action>

<acceptance_criteria>
- apps/api/internal/delivery/middleware/casbin.go implements Fiber middleware
- RT officer accessing own territory returns 200
- RT officer accessing other territory returns 403
- RW officer accessing any territory returns 200
- Resident accessing own data returns 200
- Resident accessing other user's data returns 403
- Middleware correctly maps HTTP methods to actions
- `go test ./internal/delivery/middleware -v -run TestCasbinMiddleware` passes
</acceptance_criteria>

---

### Task 5: Protected Route Examples
<read_first>
- apps/api/internal/delivery/http/protected_handler.go (create)
- apps/api/internal/delivery/middleware/auth.go
- apps/api/internal/delivery/middleware/casbin.go
</read_first>

<action>
Create example protected routes to verify RBAC:
- GET /api/users → list users (RT officer: own territory, RW officer: all territories)
- GET /api/users/:id → get user details (territory-aware)
- POST /api/users → create user (RT/RW officer only)
- Apply middleware chain: auth → casbin → handler
- Return consistent error format for unauthorized/forbidden requests
</action>

<acceptance_criteria>
- GET /api/users returns users filtered by territory for RT officer
- GET /api/users returns all users for RW officer
- GET /api/users/:id returns 403 if user belongs to different territory
- POST /api/users returns 403 for resident role
- Middleware chain correctly enforces auth → casbin → handler
- Error responses: 401 for unauthenticated, 403 for unauthorized
- `go test ./internal/delivery/http -v -run TestProtectedRoutes` passes
</acceptance_criteria>

---

## Verification

1. **Casbin Policy Enforcement:**
   ```bash
   # RT officer accessing own territory
   curl -b "paseto_token=<rt01_token>" http://localhost:8080/api/users
   ```
   Expected: 200 + users from rt-01 only

2. **Cross-Territory Isolation:**
   ```bash
   # RT officer accessing other territory
   curl -b "paseto_token=<rt01_token>" http://localhost:8080/api/users?territory=rt-02
   ```
   Expected: 403 {"error": "Forbidden"}

3. **RW Officer Access:**
   ```bash
   # RW officer accessing all territories
   curl -b "paseto_token=<rw01_token>" http://localhost:8080/api/users
   ```
   Expected: 200 + users from all territories

4. **Resident Restrictions:**
   ```bash
   # Resident trying to create user
   curl -X POST http://localhost:8080/api/users \
     -b "paseto_token=<resident_token>" \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com",...}'
   ```
   Expected: 403 {"error": "Forbidden"}