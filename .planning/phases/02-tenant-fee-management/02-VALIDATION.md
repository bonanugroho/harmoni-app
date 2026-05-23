---
phase: 2
slug: tenant-fee-management
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-05-22
---

# Phase 2 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — Go standard testing |
| **Quick run command** | `cd apps/api && go test ./internal/domain/... -short` |
| **Full suite command** | `cd apps/api && go test ./... -count=1` |
| **Estimated runtime** | ~30 seconds |

---

## Sampling Rate

- **After every task commit:** Run `cd apps/api && go test ./internal/domain/... -short`
- **After every plan wave:** Run `cd apps/api && go test ./... -count=1`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 02-Plan01-T01 | 01 | 1 | TEN-01 | T-02-01 / — | Data isolation: tenant queries include territory_id filter | unit | `go test ./internal/... -run TestTenantRepository -v` | ❌ W0 | ⬜ pending |
| 02-Plan01-T02 | 01 | 1 | FIN-01, FIN-02 | T-02-02 / — | Fee records respect tenant isolation | unit | `go test ./internal/... -run TestFeeRepository -v` | ❌ W0 | ⬜ pending |
| 02-Plan01-T03 | 01 | 1 | TEN-01 | — | N/A | integration | `go test ./internal/... -run TestTenantAPI -v` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `apps/api/internal/domain/tenant_test.go` — stubs for TEN-01
- [ ] `apps/api/internal/domain/fee_test.go` — stubs for FIN-01, FIN-02
- [ ] `apps/api/internal/infrastructure/repository/tenant_repository_test.go` — repository test stubs
- [ ] `apps/api/internal/infrastructure/repository/fee_repository_test.go` — repository test stubs

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Data isolation: RT-01 officer cannot view RT-02 tenant data | TEN-01 | Requires real Casbin enforcer + DB state | Create tenants in two territories, authenticate as RT-01, verify 403 on RT-02 tenant read |
| Fee cap enforcement | FIN-01, FIN-02 | Depends on configurable tenant cap | Set cap to 500000, attempt fee of 600000, verify validation error |

*If none: "All phase behaviors have automated verification."*

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 30s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
