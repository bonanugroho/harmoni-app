---
status: testing
phase: 02-tenant-fee-management
source: 02-01-SUMMARY.md, 02-02-SUMMARY.md, 02-03-SUMMARY.md, 02-04-SUMMARY.md
started: 2026-05-23T14:15:00Z
updated: 2026-05-23T14:15:00Z
---

## Current Test

number: 10
name: Cross-Territory Access Blocked
expected: |
  RT 01 officer accessing (GET/PUT/DELETE) a tenant in RT 02's territory returns 403 Forbidden.
awaiting: user response

## Tests

### 1. Cold Start Smoke Test
expected: Kill any running server. Clear ephemeral state (temp DBs, caches). Start the application from scratch. Server boots without errors, migrations 006–009 complete, and GET /health returns 200.
result: pass

### 2. Create Tenant
expected: POST /api/tenants with block, unit_number, occupant_name, monthly_fee (float), and territory_id returns 201 with full tenant object including id, timestamps. Missing required fields return 422.
result: pass

### 3. List Tenants (Role-Aware)
expected: RT officer GET /api/tenants returns only tenants within their territory. RW officer returns all tenants in jurisdiction. Resident returns only tenants linked via user_tenants.
result: pass

### 4. Get / Update / Delete Tenant
expected: GET /api/tenants/:id returns tenant. PUT /api/tenants/:id with updated fields returns updated object. DELETE /api/tenants/:id returns 204.
result: pass

### 5. Create Mandatory Fee
expected: POST /api/tenants/:id/fees with type=mandatory, amount, effective_date, description returns 201 with fee object including type: "mandatory". Missing type discriminator returns 400.
result: pass

### 6. Create Voluntary Fee
expected: POST /api/tenants/:id/fees with type=voluntary, amount (nullable), description, effective_date returns 201 with fee object including type: "voluntary". Optional amount field may be null.
result: pass

### 7. List / Update / Delete Fee
expected: GET /api/tenants/:id/fees returns combined list of mandatory+voluntary fees. PUT /api/tenants/:id/fees/:feeId updates amount/description. DELETE returns 204.
result: pass

### 8. Duplicate Block+Unit Rejected
expected: Creating a second tenant with same block, unit_number, and territory_id returns 409 Conflict with duplicate error code.
result: pass

### 9. Fee Validation Rules
expected: Fee amount > monthly_fee cap returns 400. Negative amount returns 400. Effective date before today returns 400. paid_at before effective_date returns 400.
result: pass

### 10. Cross-Territory Access Blocked
expected: RT 01 officer accessing (GET/PUT/DELETE) a tenant in RT 02's territory returns 403 Forbidden.
result: pass

### 11. Resident Write Protection
expected: Resident role can GET tenants but POST/PUT/DELETE any tenant or fee resource returns 403 Forbidden.
result: pass

## Summary

total: 11
passed: 11
issues: 0
pending: 0
skipped: 0
blocked: 0

## Gaps

- truth: "Register endpoint accepts role field and creates user with specified role"
  status: fixed
  reason: "User reported: register ignores role field, always creates resident. Fixed by adding Role to RegisterRequest and Register() service function."
  severity: major
  test: setup
  root_cause: "RegisterRequest struct had no Role field; AuthService.Register() hardcoded Role: 'resident'"
  artifacts:
    - path: apps/api/internal/delivery/http/auth_handler.go
      issue: "Missing Role field in RegisterRequest"
    - path: apps/api/internal/domain/service/auth_service.go
      issue: "Hardcoded Role: 'resident' in Register()"
  missing:
    - "Add Role field to RegisterRequest struct"
    - "Add role parameter to AuthService.Register()"
    - "Validate role against valid roles (resident, rt_officer, rw_officer)"
