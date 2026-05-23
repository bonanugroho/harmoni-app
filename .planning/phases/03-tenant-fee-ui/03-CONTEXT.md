# Phase 3: Tenant & Fee UI — Context

**Gathered:** 2026-05-23
**Status:** Ready for planning

<domain>
## Phase Boundary

This phase delivers the tenant and fee management frontend — responsive, mobile-first UI for RT/RW officers and residents to view, create, edit, and delete tenant records and their associated mandatory/voluntary fees. Consumes the Phase 2 backend API.

The UI-SPEC (`02-UI-SPEC.md`) provides a complete design contract: 10 components, 4 pages, copywriting, colors, typography, spacing, responsive breakpoints, interaction patterns, and file structure. All WHAT decisions are locked there. This context captures HOW implementation decisions.
</domain>

<decisions>
## Implementation Decisions

### API Client Abstraction
- **D-01:** Create a shared `request()` helper in `services/api.ts` with base URL from `VITE_API_URL`, `credentials: 'include'`, and unified error handling
- **D-02:** Auto-detect JSON vs non-JSON responses via Content-Type header (handles 204 No Content on delete)
- **D-03:** Refactor existing `services/auth.ts` to use the shared `request()` helper for consistency
- **D-04:** New service files (`services/tenants.ts`, `services/fees.ts`) use the shared helper

### Data Fetching Pattern
- **D-05:** Use TanStack Query (`@tanstack/react-query`) for all server state management
- **D-06:** Install `@tanstack/react-query` as a dependency
- **D-07:** Add `QueryClientProvider` in `App.tsx` wrapping routes
- **D-08:** Create custom hooks (`useTenants`, `useTenant`, `useFees`, `useCreateFee`, etc.) wrapping `useQuery`/`useMutation` and calling the service layer
- **D-09:** Migrate the auth/me fetch in `ProtectedRoute.tsx` to use `useQuery`

### Reusable Form Fields
- **D-10:** Create `components/ui/Input.tsx`, `components/ui/Select.tsx`, `components/ui/DatePicker.tsx` as reusable components
- **D-11:** Use native `<input type="date">` for DatePicker (styled with Tailwind, no extra dependency)
- **D-12:** Each field component handles: label, error message, aria attributes, 44px minimum touch target
- **D-13:** Existing auth forms (LoginForm, RegisterForm, ResetPasswordForm) stay as-is — no refactor
- **D-14:** TenantForm and FeeForm use the new shared field components

### Testing Approach
- **D-15:** Mock TanStack Query hooks directly (`vi.mock('@tanstack/react-query')`) — no QueryClientProvider in test wrappers
- **D-16:** Continued use of Vitest + Testing Library (matching Phase 1 pattern)
- **D-17:** Custom hooks tested via the services they call; components tested via mocked query hooks

### Agent's Discretion (areas not explicitly discussed)
- **Sidebar structure:** Hardcode sidebar items in AppLayout for now. Make configurable when Phase 4 adds Reports.
- **Search/filter:** Client-side debounced filter (300ms) on tenant list by block and unit number, as specified in UI-SPEC
- **Delete confirmation:** Follow UI-SPEC ConfirmDialog pattern (optimistic removal with error fallback)
- **Pagination:** Eager load all tenants for MVP. Server-side pagination deferred if data grows.
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### UI Design Contract (Phase 3 scope definition)
- `.planning/phases/02-tenant-fee-management/02-UI-SPEC.md` — Complete UI design contract: 10 components, 4 pages, copywriting, colors, typography, spacing, responsive rules, interaction patterns, accessibility requirements, file structure

### Project Scope & Requirements
- `.planning/ROADMAP.md` §Phase 3 — Goal and 6 success criteria for this phase
- `.planning/REQUIREMENTS.md` — Requirement IDs TEN-01, FIN-01, FIN-02 (complete); FIN-03, EXP-01 (future phases)
- `.planning/PROJECT.md` — Core value, constraints (PASETO, Casbin, data isolation, mobile-first)

### Prior Phase Decisions
- `.planning/phases/01-core-authentication-rbac/01-CONTEXT.md` — Auth patterns: httpOnly cookies, territory model, Casbin policy structure, password reset flow
- `.planning/phases/02-tenant-fee-management/01-CONTEXT.md` — Phase 2 API decisions: endpoint structure, data isolation, fee types, validation rules
- `.planning/phases/02-tenant-fee-management/02-VERIFICATION.md` — Phase 2 verification: 11/11 UAT tests passing, API ready for frontend consumption

### Codebase Maps
- `.planning/codebase/STACK.md` — Tech stack (note: outdated — Phase 1 & 2 code exists)
- `.planning/codebase/CONVENTIONS.md` — Code conventions (outdated — Phase 1 set React+TS patterns)
- `.planning/codebase/STRUCTURE.md` — Repository layout (outdated — apps/web now exists)
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `apps/web/src/services/auth.ts` — Service pattern using `fetch` + `credentials: 'include'` + `parseError`. To be refactored into shared `services/api.ts` helper
- `apps/web/src/routes/ProtectedRoute.tsx` — AuthContext provider, route guard with role check. To be extended with TanStack Query for `/auth/me`
- `apps/web/src/components/auth/LoginForm.tsx` — Form patterns: inline validation, loading/error state, 44px touch targets, aria attributes. Field markup patterns to guide Input/Select/DatePicker design
- `apps/web/src/components/auth/LoginForm.test.tsx` — Test pattern: `vi.mock` service modules, `render` with MemoryRouter, `fireEvent` + `waitFor`
- `apps/web/src/types/auth.ts` — Type pattern (`User`, request/response interfaces). Follow same pattern for `types/tenant.ts` and `types/fee.ts`

### Established Patterns
- **CSS:** Tailwind v4 with `@theme` custom properties for colors. Phase 3 adds sidebar tokens (`--color-sidebar`, `--color-sidebar-hover`, `--color-sidebar-active`)
- **Components:** Hand-rolled with Tailwind utility classes, no component library
- **Routing:** React Router v7, BrowserRouter in App.tsx, protected routes via ProtectedRoute wrapper
- **Pages:** Standalone full-screen pages (auth), new AppLayout shell wrapping all authenticated pages

### Integration Points
- **App.tsx:** Add `QueryClientProvider` wrapping routes; add /tenants route group
- **ProtectedRoute:** Wrap authenticated layout + pages (not individual pages)
- **AppLayout:** New shell component wrapping all authenticated pages — sidebar + header + content
- **API:** Consume Phase 2 endpoints (`/api/tenants`, `/api/tenants/:id`, `/api/tenants/:id/fees`, `/api/fees/:id`)
- **Auth context:** Sidebar items visibility depends on user role (RT officer sees all; resident sees read-only)
</code_context>

<specifics>
## Specific Ideas

The 02-UI-SPEC.md design contract is comprehensive — all visual and interaction specifics are documented there. No additional preferences beyond the implementation decisions above.
</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.
</deferred>

---

*Phase: 3-Tenant & Fee UI*
*Context gathered: 2026-05-23*
