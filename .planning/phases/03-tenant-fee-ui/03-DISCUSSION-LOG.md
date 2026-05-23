# Phase 3: Tenant & Fee UI - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-05-23
**Phase:** 3-Tenant & Fee UI
**Areas discussed:** API client abstraction, Data fetching pattern, Reusable form fields, Testing approach

---

## API Client Abstraction

| Option | Description | Selected |
|--------|-------------|----------|
| Shared request helper | Extract a request() function into services/api.ts with base URL, credentials, and error handling | ✓ |
| Standalone services (as-is) | Keep each service file fully self-contained | |
| Full client class | Create an ApiClient class with interceptors, typed responses | |

**User's choice:** Shared request helper
**Notes:** Also chose to put it in `services/api.ts`, refactor existing `auth.ts` to use it, and auto-detect JSON vs non-JSON responses via Content-Type header.

## Data Fetching Pattern

| Option | Description | Selected |
|--------|-------------|----------|
| Custom hooks (useTenants, useFees) | Each domain gets a hook managing loading/error/data state | |
| Inline useEffect per page | Copy the auth pattern: useEffect + fetch + local state | |
| Add TanStack Query | External library with caching, refetching, loading/error states built-in | ✓ |

**User's choice:** Add TanStack Query
**Notes:** Migrate auth/me fetch in ProtectedRoute too. QueryClientProvider in App.tsx. Mock TanStack Query hooks directly in tests.

## Reusable Form Fields

| Option | Description | Selected |
|--------|-------------|----------|
| Reusable Input/Select/DatePicker | Components in components/ui/ | ✓ |
| Inline per form (as-is) | Keep LoginForm's pattern repeated | |
| Single polymorphic Field component | One Field component with type prop | |

**User's choice:** Reusable Input/Select/DatePicker
**Notes:** Native `<input type="date">` for DatePicker. Auth forms left as-is. Only new forms use shared components.

## Testing Approach

| Option | Description | Selected |
|--------|-------------|----------|
| Follow Phase 1 pattern | Mock service modules, add QueryClientProvider wrapper | |
| Mock TanStack Query directly | Mock useQuery/useMutation instead of services | ✓ |

**User's choice:** Mock TanStack Query directly
**Notes:** Mock at `@tanstack/react-query` library level. Mock custom hooks level was also considered.

## Agent's Discretion

- Sidebar structure: hardcode items for now, make configurable in Phase 4
- Search/filter: client-side debounced filter per UI-SPEC
- Delete confirmation: follow UI-SPEC ConfirmDialog pattern
- Pagination: eager load all tenants for MVP

## Deferred Ideas

None — discussion stayed within phase scope.
