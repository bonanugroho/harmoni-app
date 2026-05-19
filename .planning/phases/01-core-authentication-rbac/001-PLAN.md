---
wave: 1
depends_on: []
files_modified:
  - apps/api/go.mod
  - apps/api/cmd/server/main.go
  - apps/api/internal/config/env.go
  - apps/api/.env.example
  - apps/api/migrations/*.sql
autonomous: true
requirements:
  - AUTH-01
  - AUTH-02
---

# Plan 1: Project Setup & Database Schema

## Objective
Initialize Go + Fiber project structure, configure environment variables, and create PostgreSQL schema with golang-migrate.

## Tasks

### Task 1: Initialize Go Module & Dependencies
<read_first>
- apps/api/go.mod (create)
- .planning/PROJECT.md (tech stack reference)
</read_first>

<action>
Initialize Go module in apps/api directory:
- `go mod init harmoni-api`
- Add dependencies: `github.com/gofiber/fiber/v2`, `github.com/jackc/pgx/v5`, `github.com/golang-migrate/migrate/v4`, `github.com/casbin/casbin/v2`, `golang.org/x/crypto`, `github.com/o1egl/paseto`
- Create basic project structure per Clean Architecture layout
</action>

<acceptance_criteria>
- apps/api/go.mod exists with all required dependencies
- apps/api/cmd/server/main.go exists with Fiber app initialization
- `go mod tidy` runs without errors
- `go run cmd/server/main.go` starts Fiber server on port 8080
</acceptance_criteria>

---

### Task 2: Environment Configuration
<read_first>
- apps/api/internal/config/env.go (create)
- apps/api/.env.example (create)
</read_first>

<action>
Create environment configuration:
- Load env vars: DATABASE_URL, PASETO_SECRET_KEY, EMAIL_API_KEY, APP_ENV, APP_PORT
- Generate 32-byte PASETO key: `openssl rand -hex 32`
- Create .env.example with all required variables documented
- Validate required env vars on startup, exit with error if missing
</action>

<acceptance_criteria>
- apps/api/.env.example contains: DATABASE_URL, PASETO_SECRET_KEY, EMAIL_API_KEY, APP_ENV, APP_PORT
- apps/api/internal/config/env.go loads and validates all env vars
- Missing env var causes application to exit with descriptive error
- PASETO_SECRET_KEY validation ensures 32-byte hex string
</acceptance_criteria>

---

### Task 3: Database Migration Setup
<read_first>
- apps/api/migrations/ (create directory)
- .planning/phases/01-core-authentication-rbac/01-RESEARCH.md §4 (golang-migrate patterns)
</read_first>

<action>
Create golang-migrate migration files:
- 001_create_territories_table.up.sql: territories table with id, name, type, parent_id, created_at
- 001_create_territories_table.down.sql: DROP TABLE territories
- 002_create_users_table.up.sql: users table with id (UUID DEFAULT uuidv7()), email (UNIQUE), password_hash, role (CHECK constraint), territory_id (FK), full_name, phone, is_active, created_at, updated_at
- 002_create_users_table.down.sql: DROP TABLE users
- 003_create_sessions_table.up.sql: sessions table with id (UUID DEFAULT uuidv7()), user_id (UUID DEFAULT uuidv7() FK), token_hash, expires_at, created_at
- 003_create_sessions_table.down.sql: DROP TABLE sessions
- 004_create_password_reset_tokens_table.up.sql: password_reset_tokens table with id (UUID DEFAULT uuidv7()), user_id (UUID DEFAULT uuidv7() FK), token_hash, expires_at, used (BOOLEAN), created_at
- 004_create_password_reset_tokens_table.down.sql: DROP TABLE password_reset_tokens
- 005_seed_territories.up.sql: INSERT sample territories (rt-01, rt-02, rw-01)
- 005_seed_territories.down.sql: DELETE FROM territories

**IMPORTANT:** PostgreSQL 18+ provides native `uuidv7()` function. Use `DEFAULT uuidv7()` for all UUID columns. No application-layer UUID generation needed.
</action>

<acceptance_criteria>
- All 10 migration files exist in apps/api/migrations/
- `migrate -path migrations -database "postgres://..." up` applies all migrations successfully
- `migrate -path migrations -database "postgres://..." down` rolls back all migrations successfully
- territories table contains: id (VARCHAR PK), name, type (rt/rw), parent_id (FK to territories)
- users table contains: id (UUID DEFAULT uuidv7()), email (UNIQUE), password_hash, role (CHECK constraint), territory_id (FK), full_name, phone, is_active, timestamps
- sessions table contains: id (UUID DEFAULT uuidv7()), user_id (UUID DEFAULT uuidv7() FK), token_hash, expires_at, created_at
- password_reset_tokens table contains: id (UUID DEFAULT uuidv7()), user_id (UUID DEFAULT uuidv7() FK), token_hash, expires_at, used (BOOLEAN), created_at
- Indexes exist on: users.email, users.territory_id, sessions.user_id, sessions.expires_at, password_reset_tokens.user_id
- All UUID columns use `DEFAULT uuidv7()` (PostgreSQL 18+ native function)
</acceptance_criteria>

---

### Task 4: Database Connection & Migration Runner
<read_first>
- apps/api/internal/infrastructure/database/connection.go (create)
- apps/api/internal/config/env.go
</read_first>

<action>
Implement database connection and migration runner:
- Create PostgreSQL connection using pgx pool
- Implement migration runner using golang-migrate CLI integration
- Add health check endpoint: GET /health → {"status": "ok", "database": "connected"}
- Log migration status on startup
</action>

<acceptance_criteria>
- apps/api/internal/infrastructure/database/connection.go establishes pgx pool connection
- Migration runner applies pending migrations on startup
- GET /health returns 200 with {"status": "ok", "database": "connected"}
- Database connection errors are logged and application exits gracefully
- `go test ./internal/infrastructure/database -v` passes
</acceptance_criteria>

---

## Verification

1. **Project Structure:**
   ```bash
   tree apps/api -L 4
   ```
   Expected: Clean Architecture layout with cmd/, internal/, migrations/, pkg/

2. **Database Setup:**
   ```bash
   cd apps/api && migrate -path migrations -database "$DATABASE_URL" up
   ```
   Expected: All migrations applied, tables created

3. **Server Startup:**
   ```bash
   cd apps/api && go run cmd/server/main.go
   curl http://localhost:8080/health
   ```
   Expected: {"status": "ok", "database": "connected"}

4. **Environment Validation:**
   ```bash
   # Remove PASETO_SECRET_KEY from .env
   cd apps/api && go run cmd/server/main.go
   ```
   Expected: Exit with error "PASETO_SECRET_KEY is required"