# Phase 3 Security Audit — Tenant & Fee UI

**Audit date:** 2026-05-24
**Phase:** 03-tenant-fee-ui
**Scope:** Frontend-only — `apps/web/src/`
**ASVS Level:** 1

## Threat Verification

Total threats in register: 17 (10 mitigate, 7 accepted)
**Closed:** 10/10 | **Open:** 0/10

### Closed — Mitigate

| Threat ID | Category | Disposition | Evidence |
|-----------|----------|-------------|----------|
| T-03-01 | Tampering | mitigate | `apps/web/src/services/api.ts:12` — `credentials: 'include'` set on `fetch()` calls. PASETO cookie sent only to same-origin API. |
| T-03-03 | Tampering | mitigate | `apps/web/package.json:14-15` — `@tanstack/react-query@^5.100.14` and `lucide-react@^1.16.0` declared as dependencies. |
| T-03-SC | Tampering | mitigate | `apps/web/package.json:14` — `@tanstack/react-query@^5.100.14` declared. |
| T-03-04 | Information Disclosure | mitigate | `apps/web/src/components/layout/AppLayout.tsx:20,29-32` — Settings nav item has `roles: ['rt_officer', 'rw_officer']` constraint; `visibleItems` filters by `user?.role` from `useAuth()`. |
| T-03-05 | Tampering | mitigate | `apps/web/src/components/ui/ConfirmDialog.tsx:4,9-10` — Component interface declares `isOpen: boolean`, `onConfirm: () => void`, `onCancel: () => void` as controlled props. |
| T-03-09 | Tampering | mitigate | All 6 mutation hooks call `invalidateQueries` in `onSuccess`: `useCreateTenant.ts:10`, `useUpdateTenant.ts:10`, `useDeleteTenant.ts:9`, `useCreateFee.ts:10`, `useUpdateFee.ts:10`, `useDeleteFee.ts:9`. |
| T-03-10 | Denial of Service | mitigate | `apps/web/src/App.tsx:16` — `QueryClient` configured with `staleTime: 30_000` in `defaultOptions.queries`. |
| T-03-11 | Tampering | mitigate | `apps/web/src/components/tenants/TenantForm.tsx:55-109` — `validate()` function checks `block`, `unit_number`, `monthly_fee`, and mandatory fee entries before submit. |
| T-03-14 | Tampering | mitigate | `apps/web/src/components/fees/FeeForm.tsx:50-92` — `validate()` checks: amount > 0, amount ≤ monthlyFee, total mandatory fees ≤ monthlyFee, effective date not in past, payment date ≥ effective date. |
| T-03-16 | Tampering | mitigate | `apps/web/src/pages/TenantDetailPage.tsx:166-173` — Fee delete uses `<ConfirmDialog>` with `isOpen`/`onConfirm`/`onCancel`. `apps/web/src/pages/TenantEditPage.tsx:85-94` — Tenant delete uses same pattern. |

### Closed — Accepted (documented risks, no code verification required)

| Threat ID | Category | Disposition | Evidence |
|-----------|----------|-------------|----------|
| T-03-02 | Information Disclosure | accept | Types match API — no PII beyond block/unit_number. Documented in threat register. |
| T-03-06 | Information Disclosure | accept | Sidebar only hides links visually; ProtectedRoute + Casbin enforce server-side. |
| T-03-07 | Information Disclosure | accept | AuthContext user has role/territory, no sensitive data. |
| T-03-08 | Tampering | accept | Client-side role check is UX only; backend Casbin enforces. |
| T-03-12 | Tampering | accept | Delete button visible but backend Casbin prevents unauthorized deletes. |
| T-03-13 | Information Disclosure | accept | Search filters client-side only; backend enforces territory isolation. |
| T-03-15 | Denial of Service | accept | No rate limiting on frontend; backend should handle. |

## Unregistered Flags

None — no new attack surface was introduced beyond what the threat register covers.

## Verdict

**SECURED** — All 10 declared mitigations are verified present in the implementation code. All 7 accepted risks are documented. No implementation gaps found.

**Phase may proceed.**
