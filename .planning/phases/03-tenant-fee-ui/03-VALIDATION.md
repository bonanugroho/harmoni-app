---
phase: 3
slug: tenant-fee-ui
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-05-24
---

# Phase 3 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Vitest v4.1.6 + Testing Library React v16.3.2 |
| **Config file** | `apps/web/vitest.config.js` |
| **Quick run command** | `npm test -- --run` (or `npx vitest run` from `apps/web/`) |
| **Full suite command** | `npm test` (or `npx vitest` for watch mode) |
| **Estimated runtime** | ~60 seconds |

---

## Sampling Rate

- **After every task commit:** Run `npm test -- --run`
- **After every plan wave:** Run `npm test -- --run`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 60 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 03-01-01 | 01 | 1 | TEN-01, FIN-01, FIN-02 | T-03-01 | Auth header forwarded via credentials:include | unit | `npx vitest run services/api.test.ts` | ❌ W0 | ⬜ pending |
| 03-01-02 | 01 | 1 | TEN-01, FIN-01, FIN-02 | T-03-04 | Validation before API calls | unit | `npx vitest run services/tenants.test.ts services/fees.test.ts` | ❌ W0 | ⬜ pending |
| 03-01-03 | 01 | 1 | TEN-01, FIN-01, FIN-02 | — | QueryClientProvider wraps routes | unit | `npx vitest run` | ❌ W0 | ⬜ pending |
| 03-02-01 | 02 | 2 | TEN-01 | — | All inputs 44px min touch target | unit | `npx vitest run components/ui/` | ❌ W0 | ⬜ pending |
| 03-02-02 | 02 | 2 | TEN-01, FIN-01 | — | Accessible & responsive | unit | `npx vitest run components/ui/` | ❌ W0 | ⬜ pending |
| 03-02-03 | 02 | 2 | — | T-03-05, T-03-06 | Role-based sidebar visibility | unit | `npx vitest run components/layout/` | ❌ W0 | ⬜ pending |
| 03-03-01 | 03 | 2 | TEN-01, FIN-01, FIN-02 | T-03-02 | Error states handled | unit | `npx vitest run` | ❌ W0 | ⬜ pending |
| 03-03-02 | 03 | 2 | — | T-03-03 | Token refresh via useQuery | unit | `npx vitest run routes/` | ❌ W0 | ⬜ pending |
| 03-04-01 | 04 | 3 | TEN-01 | — | Card renders all fields | unit | `npx vitest run components/tenants/TenantCard.test.tsx` | ❌ W0 | ⬜ pending |
| 03-04-02 | 04 | 3 | TEN-01 | T-03-07, T-03-08 | Form validation + submit | unit | `npx vitest run pages/TenantListPage.test.tsx` | ❌ W0 | ⬜ pending |
| 03-04-03 | 04 | 3 | TEN-01 | T-03-09 | Edit pre-populates + delete | unit | `npx vitest run components/tenants/TenantForm.test.tsx` | ❌ W0 | ⬜ pending |
| 03-05-01 | 05 | 3 | FIN-01, FIN-02 | T-03-10 | Fee sections display correctly | unit | `npx vitest run components/fees/FeeList.test.tsx` | ❌ W0 | ⬜ pending |
| 03-05-02 | 05 | 3 | FIN-01, FIN-02 | T-03-11, T-03-12, T-03-13 | FeeType selector + validation | unit | `npx vitest run components/fees/FeeForm.test.tsx` | ❌ W0 | ⬜ pending |
| 03-05-03 | 05 | 3 | FIN-01, FIN-02 | T-03-14, T-03-15, T-03-16 | Complete fee management page | unit | `npx vitest run pages/TenantDetailPage.test.tsx` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `services/api.test.ts` — shared helper behavior (base URL, credentials, error handling, 204 handling)
- [ ] `services/tenants.test.ts` — calls correct endpoints, handles responses
- [ ] `services/fees.test.ts` — calls correct endpoints, handles responses
- [ ] `components/ui/Input.test.tsx` — renders label, shows error, aria attributes, 44px touch target
- [ ] `components/ui/Select.test.tsx` — renders options, shows error, handles change
- [ ] `components/ui/DatePicker.test.tsx` — renders date input, handles change
- [ ] `components/ui/StatusBadge.test.tsx` — Occupied vs Vacant, Paid vs Unpaid
- [ ] `components/ui/ConfirmDialog.test.tsx` — shows on trigger, Cancel dismisses, Delete confirms
- [ ] `components/ui/PageHeader.test.tsx` — renders title + CTA
- [ ] `components/ui/LoadingSkeleton.test.tsx` — renders skeleton placeholders
- [ ] `components/ui/EmptyState.test.tsx` — shows message + CTA
- [ ] `components/layout/AppLayout.test.tsx` — sidebar visible/hidden, responsive, navigation
- [ ] `components/tenants/TenantCard.test.tsx` — renders tenant info, occupancy badge, fee summary
- [ ] `components/tenants/TenantForm.test.tsx` — field validation, add fee row, submit, edit mode
- [ ] `components/fees/FeeList.test.tsx` — mandatory + voluntary sections, empty states
- [ ] `components/fees/FeeForm.test.tsx` — type selector, validation, submit
- [ ] `pages/TenantListPage.test.tsx` — loading skeleton, tenant cards, empty state, search filter
- [ ] `pages/TenantCreatePage.test.tsx` — renders form, cancel goes back
- [ ] `pages/TenantEditPage.test.tsx` — pre-populated fields, delete button
- [ ] `pages/TenantDetailPage.test.tsx` — tenant info header, fee sections, record fee button

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Sidebar responsive behavior (collapse/expand on resize) | SC-4 | Requires viewport resize interaction in browser | Open app on mobile viewport (375px). Hamburger icon should be visible. Click it — sidebar slides in. Click backdrop — sidebar closes. Resize to 1024px+ — sidebar fixed visible. |
| Touch target physical size | SC-5 | Cannot be measured in unit tests | Use browser DevTools to verify all buttons, sidebar links, interactive cards have ≥44px computed height. |
| Pull-to-refresh on tenant list (mobile) | SC-1 | Requires native mobile browser behavior | Open tenant list on mobile device or Chrome DevTools mobile emulation. Pull down — page should refresh. |
| Form validation visual feedback | SC-2 | Visual confirmation of error states | Submit empty tenant form — each field should show inline error message in red, focus should move to first error field. |
| Delete confirmation dialog | SC-3 | Visual/interaction verification | Click trash icon on a fee — ConfirmDialog should appear with "Delete" and "Cancel" buttons. Press Escape — dialog closes. Click backdrop — dialog closes. |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 60s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
