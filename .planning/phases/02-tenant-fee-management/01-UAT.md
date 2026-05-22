---
status: planning
phase: 02-tenant-fee-management
source: 001-SUMMARY.md, 002-SUMMARY.md, 003-SUMMARY.md, 005-SUMMARY.md
started: 2026-05-22T09:00:00Z
updated: 2026-05-22T09:30:00Z
---

## Current Test

[testing complete]

## Tests

### 1. Cold Start Smoke Test
expected: Kill any running server/service. Clear ephemeral state. Start the application from scratch. Server boots without errors, migrations complete, and health check returns live data.
result: pass

### 2. Health Check Endpoint
expected: GET /health returns {"status": "ok", "database": "connected"} with 200 status.
result: pass

### 3. User Registration
expected: POST /auth/register with valid email, password (8+ chars, uppercase, lowercase, number, symbol), and full name returns 201 with user object. Weak passwords are rejected with validation error.
result: pass

### 4. User Login
expected: POST /auth/login with valid credentials returns 200 with user object and sets httpOnly cookie named "token". Invalid credentials return 401 with error message.
result: pass

### 5. Password Reset Request
expected: POST /auth/reset with valid email returns 200 always (even for non-existent emails to prevent enumeration). Email with reset link is sent.
result: pass

### 6. Password Reset Confirmation
expected: POST /auth/reset/confirm with valid token and new password returns 200. Expired or used tokens return 400. New password must meet complexity requirements.
result: pass

### 7. Protected Route Access (API)
expected: GET /api/protected without valid auth cookie returns 401. With valid auth cookie, returns 200 with user data.
result: pass

### 8. Territory Isolation (RT Officers)
expected: RT officer can only access resources from their assigned territory. Attempting to access another RT's data returns 403 Forbidden.
result: pass

### 9. RW Officer Wildcard Access
expected: RW officer can access resources from all territories within their RW jurisdiction. No 403 errors for cross-territory access.
result: pass

### 10. Role-Based Access Control (Residents)
expected: Resident can read data (GET requests) but cannot create, update, or delete (POST/PUT/DELETE) protected resources. Write attempts return 403 Forbidden.
result: pass

### 11. Frontend Auth Pages Render
expected: Navigate to /login, /register, /reset in browser. Each page renders with Harmoni branding, appropriate form fields, and links to other auth pages. Forms submit and show loading states.
result: pass

### 12. Frontend Protected Route Redirect
expected: Navigate to /dashboard without being logged in. Browser redirects to /login. After login, /dashboard is accessible and shows "Welcome to Harmoni" placeholder.
result: pass

### 13. TypeScript Build Verification
expected: `npm run build` in apps/web completes without errors. Output includes JS and CSS bundles. No TypeScript type errors.
result: pass

## Summary

total: 13
passed: 13
issues: 0
pending: 0
skipped: 0
blocked: 0

## Gaps

[none]

### 3. User Registration
expected: POST /auth/register with valid email, password (8+ chars, uppercase, lowercase, number, symbol), and full name returns 201 with user object. Weak passwords are rejected with validation error.
result: [pending]

### 4. User Login
expected: POST /auth/login with valid credentials returns 200 with user object and sets httpOnly cookie named "token". Invalid credentials return 401 with error message.
result: [pending]

### 5. Password Reset Request
expected: POST /auth/reset with valid email returns 200 always (even for non-existent emails to prevent enumeration). Email with reset link is sent.
result: [pending]

### 6. Password Reset Confirmation
expected: POST /auth/reset/confirm with valid token and new password returns 200. Expired or used tokens return 400. New password must meet complexity requirements.
result: [pending]

### 7. Protected Route Access (API)
expected: GET /api/protected without valid auth cookie returns 401. With valid auth cookie, returns 200 with user data.
result: [pending]

### 8. Territory Isolation (RT Officers)
expected: RT officer can only access resources from their assigned territory. Attempting to access another RT's data returns 403 Forbidden.
result: [pending]

### 9. RW Officer Wildcard Access
expected: RW officer can access resources from all territories within their RW jurisdiction. No 403 errors for cross-territory access.
result: [pending]

### 10. Role-Based Access Control (Residents)
expected: Resident can read data (GET requests) but cannot create, update, or delete (POST/PUT/DELETE) protected resources. Write attempts return 403 Forbidden.
result: [pending]

### 11. Frontend Auth Pages Render
expected: Navigate to /login, /register, /reset in browser. Each page renders with Harmoni branding, appropriate form fields, and links to other auth pages. Forms submit and show loading states.
result: [pending]

### 12. Frontend Protected Route Redirect
expected: Navigate to /dashboard without being logged in. Browser redirects to /login. After login, /dashboard is accessible and shows "Welcome to Harmoni" placeholder.
result: [pending]

### 13. TypeScript Build Verification
expected: `npm run build` in apps/web completes without errors. Output includes JS and CSS bundles. No TypeScript type errors.
result: [pending]

## Summary

total: 13
passed: 13
issues: 0
pending: 0
skipped: 0
blocked: 0

## Gaps

[none yet]
