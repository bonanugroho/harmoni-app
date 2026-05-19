# Graph Report - harmoni-app  (2026-05-19)

## Corpus Check
- 66 files · ~22,858 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 378 nodes · 596 edges · 27 communities (22 shown, 5 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 82 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f93e6aec`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]

## God Nodes (most connected - your core abstractions)
1. `newTestAuthService()` - 19 edges
2. `setupTestHandler()` - 16 edges
3. `ResetEnforcerForTest()` - 16 edges
4. `InitEnforcer()` - 16 edges
5. `ProtectedHandler` - 14 edges
6. `newTestPasetoService()` - 13 edges
7. `main()` - 12 edges
8. `getUserClaims()` - 12 edges
9. `newTestEnforcerForProtected()` - 11 edges
10. `newEnforcerForTerritory()` - 11 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewPostgresUserRepository()`  [INFERRED]
  apps/api/cmd/server/main.go → apps/api/internal/infrastructure/repository/user_repository.go
- `main()` --calls--> `NewPostgresPasswordResetTokenRepository()`  [INFERRED]
  apps/api/cmd/server/main.go → apps/api/internal/infrastructure/repository/password_reset_token_repository.go
- `main()` --calls--> `NewPasetoService()`  [INFERRED]
  apps/api/cmd/server/main.go → apps/api/internal/infrastructure/auth/paseto.go
- `main()` --calls--> `NewResendEmailService()`  [INFERRED]
  apps/api/cmd/server/main.go → apps/api/internal/infrastructure/email/resend.go
- `main()` --calls--> `NewAuthService()`  [INFERRED]
  apps/api/cmd/server/main.go → apps/api/internal/domain/service/auth_service.go

## Communities (27 total, 5 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.05
Nodes (29): getPasswordStrength(), PasswordStrength, RegisterForm(), TERRITORIES, Territory, passwordInput, select, strengthText (+21 more)

### Community 1 - "Community 1"
Cohesion: 0.08
Nodes (24): generateRandomToken(), hashToken(), NewAuthService(), newMockResetRepo(), newMockUserRepo(), newTestAuthService(), TestAuthService_Login_InactiveUser(), TestAuthService_Login_InvalidCredentials() (+16 more)

### Community 2 - "Community 2"
Cohesion: 0.07
Nodes (14): Config, LoadEnv(), validatePasetoKey(), EnvValidationError, NewConnection(), RunMigrations(), TestNewConnection_EmptyURL(), TestNewConnection_InvalidURL() (+6 more)

### Community 3 - "Community 3"
Cohesion: 0.13
Nodes (28): GetEnforcer(), InitEnforcer(), newEnforcerForTerritoryTest(), TestCasbinEnforcer_AddRemovePolicy(), TestCasbinEnforcer_AddRemoveRoleLink(), TestCasbinEnforcer_Enforce(), TestCasbinEnforcer_EnforceWithTerritory(), TestCasbinEnforcer_EnforceWithTerritory_RWOfficer() (+20 more)

### Community 4 - "Community 4"
Cohesion: 0.1
Nodes (18): containsStr(), newMockResetRepo(), newMockUserRepo(), setupTestHandler(), TestAuthHandler_Login_InvalidCredentials(), TestAuthHandler_Login_Success(), TestAuthHandler_Register_DuplicateEmail(), TestAuthHandler_Register_MissingFields() (+10 more)

### Community 5 - "Community 5"
Cohesion: 0.16
Nodes (10): buildPolicyWithTerritory(), newEnforcerForTerritory(), parsePolicyLines(), TestCasbinPolicy_ResidentOtherTerritory(), TestCasbinPolicy_ResidentReadOnly(), TestCasbinPolicy_RoleInheritance(), TestCasbinPolicy_RTOfficerOtherTerritory(), TestCasbinPolicy_RTOfficerOwnTerritory() (+2 more)

### Community 6 - "Community 6"
Cohesion: 0.11
Nodes (18): Architecture, Constraints, Conventions, Core language & runtime, Developer Profile, Development dependencies (absent), Ecosystem summary, Expected future layers (placeholders) (+10 more)

### Community 7 - "Community 7"
Cohesion: 0.3
Nodes (15): DefaultPublicRoutes(), NewAuthMiddleware(), doRequest(), newTestPasetoService(), readBody(), setupTestApp(), TestAuthMiddleware_AuthRouteBypass(), TestAuthMiddleware_CustomPublicRoutes() (+7 more)

### Community 8 - "Community 8"
Cohesion: 0.2
Nodes (5): getFilterType(), getUserClaims(), NewProtectedHandler(), TestGetFilterType(), ProtectedHandler

### Community 9 - "Community 9"
Cohesion: 0.18
Nodes (12): Claims, NewPasetoService(), newTestService(), TestPaseto_ExpiredToken(), TestPaseto_GenerateAndValidateToken(), TestPaseto_InvalidToken(), TestPaseto_NewServiceInvalidKey(), TestPaseto_TamperedToken() (+4 more)

### Community 10 - "Community 10"
Cohesion: 0.56
Nodes (14): ResetEnforcerForTest(), doProtectedRequest(), newTestEnforcerForProtected(), newTestPasetoServiceForProtected(), readProtectedBody(), setupProtectedApp(), TestProtectedRoutes_ErrorFormat(), TestProtectedRoutes_GetUser() (+6 more)

### Community 11 - "Community 11"
Cohesion: 0.14
Nodes (6): NewAuthHandler(), AuthHandler, LoginRequest, RegisterRequest, ResetConfirmRequest, ResetRequest

### Community 12 - "Community 12"
Cohesion: 0.25
Nodes (7): MockEmailService, NewResendEmailService(), passwordResetHTML(), TestMockEmailService(), TestResend_SendPasswordResetEmail_ContainsResetLink(), TestResend_SendPasswordResetEmail_MissingAPIKey(), ResendEmailService

### Community 13 - "Community 13"
Cohesion: 0.25
Nodes (7): 1. Introduction, 2.1 User Management & Access, 2.2 Data & Transaction Management, 2.3 Reporting & Dashboard, 2. Functional Requirements, 3. Non-Functional Requirements, Requirements Document - Harmoni Project

### Community 16 - "Community 16"
Cohesion: 0.33
Nodes (5): 1. Tech Stack, 2. Backend Architecture (Clean Architecture), 3. Security & Authorization, 5. Development & Deployment, Specification Document - Harmoni App Project

### Community 18 - "Community 18"
Cohesion: 0.5
Nodes (3): Expanding the ESLint configuration, React Compiler, React + Vite

## Knowledge Gaps
- **53 isolated node(s):** `ResetRequest`, `ConfirmResetRequest`, `Territory`, `PasswordStrength`, `TERRITORIES` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `Community 2` to `Community 1`, `Community 3`, `Community 8`, `Community 9`, `Community 11`, `Community 12`?**
  _High betweenness centrality (0.251) - this node is a cross-community bridge._
- **Why does `NewPasetoService()` connect `Community 9` to `Community 1`, `Community 2`, `Community 4`, `Community 7`, `Community 10`?**
  _High betweenness centrality (0.160) - this node is a cross-community bridge._
- **Why does `InitEnforcer()` connect `Community 3` to `Community 2`, `Community 10`?**
  _High betweenness centrality (0.138) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `newTestAuthService()` (e.g. with `NewPasetoService()` and `NewAuthService()`) actually correct?**
  _`newTestAuthService()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **Are the 3 inferred relationships involving `setupTestHandler()` (e.g. with `NewPasetoService()` and `NewAuthService()`) actually correct?**
  _`setupTestHandler()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **Are the 15 inferred relationships involving `ResetEnforcerForTest()` (e.g. with `newTestEnforcerForMiddleware()` and `TestCasbinMiddleware_RTOfficerOwnTerritory()`) actually correct?**
  _`ResetEnforcerForTest()` has 15 INFERRED edges - model-reasoned connections that need verification._
- **Are the 15 inferred relationships involving `InitEnforcer()` (e.g. with `main()` and `newTestEnforcerForMiddleware()`) actually correct?**
  _`InitEnforcer()` has 15 INFERRED edges - model-reasoned connections that need verification._