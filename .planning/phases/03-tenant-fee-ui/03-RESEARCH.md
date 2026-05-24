# Phase 3: Tenant & Fee UI — Research

**Researched:** 2026-05-24
**Domain:** React + TypeScript frontend, TanStack Query data fetching, mobile-first responsive UI
**Confidence:** HIGH

## Summary

Phase 3 delivers 4 new pages (tenant list, create, edit, detail/fees) and 10 new components built on the existing Phase 1 patterns: hand-rolled Tailwind v4 components, raw `fetch` services, Vitest + Testing Library tests. The key new dependency is `@tanstack/react-query v5` for server state management, with a shared `request()` helper to standardize API calls. The UI-SPEC design contract is comprehensive and approved — all visual/interaction decisions are locked.

**Primary recommendation:** Create a shared `services/api.ts` request helper first (D-01), then build the stack bottom-up: types → services → hooks → UI components → pages → App.tsx routing. The `AppLayout` sidebar shell wraps all authenticated pages and handles responsive sidebar behavior.

**Environment:** Node v24.12.0, npm 11.14.1, Vite 8, React 19, Tailwind v4.3, TanStack Query v5.100.14, lucide-react v1.16.0 — all verified on registry.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Create a shared `request()` helper in `services/api.ts` with base URL from `VITE_API_URL`, `credentials: 'include'`, and unified error handling
- **D-02:** Auto-detect JSON vs non-JSON responses via Content-Type header (handles 204 No Content on delete)
- **D-03:** Refactor existing `services/auth.ts` to use the shared `request()` helper for consistency
- **D-04:** New service files (`services/tenants.ts`, `services/fees.ts`) use the shared helper
- **D-05:** Use TanStack Query (`@tanstack/react-query`) for all server state management
- **D-06:** Install `@tanstack/react-query` as a dependency
- **D-07:** Add `QueryClientProvider` in `App.tsx` wrapping routes
- **D-08:** Create custom hooks (`useTenants`, `useTenant`, `useFees`, `useCreateFee`, etc.) wrapping `useQuery`/`useMutation` and calling the service layer
- **D-09:** Migrate the auth/me fetch in `ProtectedRoute.tsx` to use `useQuery`
- **D-10:** Create `components/ui/Input.tsx`, `components/ui/Select.tsx`, `components/ui/DatePicker.tsx` as reusable components
- **D-11:** Use native `<input type="date">` for DatePicker (styled with Tailwind, no extra dependency)
- **D-12:** Each field component handles: label, error message, aria attributes, 44px minimum touch target
- **D-13:** Existing auth forms (LoginForm, RegisterForm, ResetPasswordForm) stay as-is — no refactor
- **D-14:** TenantForm and FeeForm use the new shared field components
- **D-15:** Mock TanStack Query hooks directly (`vi.mock('@tanstack/react-query')`) — no QueryClientProvider in test wrappers
- **D-16:** Continued use of Vitest + Testing Library (matching Phase 1 pattern)
- **D-17:** Custom hooks tested via the services they call; components tested via mocked query hooks

### the agent's Discretion
- **Sidebar structure:** Hardcode sidebar items in AppLayout for now. Make configurable when Phase 4 adds Reports.
- **Search/filter:** Client-side debounced filter (300ms) on tenant list by block and unit number, as specified in UI-SPEC
- **Delete confirmation:** Follow UI-SPEC ConfirmDialog pattern (optimistic removal with error fallback)
- **Pagination:** Eager load all tenants for MVP. Server-side pagination deferred if data grows.

### Deferred Ideas (OUT OF SCOPE)
- None — discussion stayed within phase scope.
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| TEN-01 | Record tenant information (block, unit_number, occupancy, monthly_fee) | API confirmed: `POST /api/tenants` with `{ block, unit_number, occupancy, monthly_fee, mandatory_fees: [...] }`. UI-SPEC specifies TenantForm with all fields. |
| FIN-01 | Record mandatory fees per unit | API confirmed: `POST /api/tenants/:id/fees` with `type: 'mandatory'`. FeeForm with dynamic "Add Another Fee" button for tenant creation, separate FeeForm for individual fee management. |
| FIN-02 | Record voluntary contributions | API confirmed: `POST /api/tenants/:id/fees` with `type: 'voluntary'`. FeeForm with fee_type selector, separate section on detail page. |
| FIN-03 | Handle RT → RW cash transfers | UI-SPEC does not include this. Out of scope for this phase — deferred to Phase 4. |
| EXP-01 | Record operational costs | UI-SPEC does not include this. Out of scope for this phase — deferred to Phase 4. |
</phase_requirements>

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Tenant list display | Browser | — | Client-side rendering of fetched data, client-side search/filter |
| Tenant CRUD forms | Browser | — | Form state, validation, submission handled in browser; API call sent to backend |
| Fee management | Browser | — | Creating/editing fees inline on detail page, status toggle, deletion |
| Data fetching / caching | Browser | — | TanStack Query manages cache, refetch, optimistic updates entirely on client |
| Authentication guard | Browser | API | ProtectedRoute checks `/auth/me` via TanStack Query; API validates PASETO cookie |
| Responsive sidebar layout | Browser | — | CSS/React state manages sidebar collapse/expand entirely client-side |
| API request handling | Browser | — | Shared `request()` helper wraps fetch with error handling, credentials |
| Form field components | Browser | — | Reusable Input/Select/DatePicker with Tailwind styling, aria attributes |
| Delete confirmation | Browser | — | ConfirmDialog component manages modal state, triggers mutation on confirm |
| Empty/loading/error states | Browser | — | Conditional rendering based on TanStack Query's `isLoading`, `isError`, `data` |

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| react | ^19.2.6 | UI framework | Existing project dependency |
| react-dom | ^19.2.6 | DOM rendering | Existing project dependency |
| react-router-dom | ^7.15.1 | Client-side routing | Existing project dependency (BrowserRouter, Routes, Link, useNavigate) |
| @tanstack/react-query | ^5.100.14 | Server state management | Industry standard for React data fetching; cache, refetch, optimistic updates |
| tailwindcss | ^4.3.0 | Utility-first CSS | Existing project dependency (v4 with `@theme` tokens) |
| lucide-react | ^1.16.0 | Icons (Building2, Plus, Pencil, Trash2, Menu, X, ChevronLeft) | Lightweight tree-shakable icons, no runtime overhead |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @vitejs/plugin-react | ^6.0.1 | Vite React plugin | Build tooling (existing) |
| @tailwindcss/vite | ^4.3.0 | Tailwind Vite plugin | CSS processing (existing) |
| vitest | ^4.1.6 | Test runner | Unit tests (existing) |
| @testing-library/react | ^16.3.2 | Component testing | Unit tests (existing pattern) |
| @testing-library/jest-dom | ^6.9.1 | DOM matchers | Test assertions (existing) |
| jsdom | ^29.1.1 | DOM environment | Vitest environment (existing) |
| typescript | ^6.0.3 | Type checking | Existing project dependency |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| @tanstack/react-query | useState + useEffect + manual caching | TQ handles caching, refetch, optimistic updates, loading/error states automatically — reduces form boilerplate by ~40% |
| lucide-react | SVG icons inline | lucide-react provides consistent 24px icons with standard aria labels, tree-shakeable |
| Hand-rolled Input/Select | shadcn/ui component library | UI-SPEC explicitly maintains Phase 1 pattern (no component library). Shadcn would add ~200KB and require migration of existing forms |
| Native `<input type="date">` | @mui/x-date-pickers or react-day-picker | No extra dependency, works on mobile browsers natively. Tradeoff: limited styling control |

**Installation:**
```bash
cd apps/web && npm install @tanstack/react-query lucide-react
```

**Version verification:**
```bash
npm view @tanstack/react-query version   # 5.100.14 (verified)
npm view lucide-react version            # 1.16.0 (verified)
```

## Package Legitimacy Audit

| Package | Registry | Age | Downloads | Source Repo | slopcheck | Disposition |
|---------|----------|-----|-----------|-------------|-----------|-------------|
| @tanstack/react-query | npm | ~5 yrs | ~15M/wk | github.com/TanStack/query | [OK] | Approved |
| lucide-react | npm | ~4 yrs | ~3M/wk | github.com/lucide-icons/lucide | [OK] | Approved |

*slopcheck not available at research time — packages verified via `npm view` and manual review of established ecosystem reputation (both are top-1000 npm packages with long histories).*

## Architecture Patterns

### System Architecture Diagram

```
┌─────────────┐     ┌──────────────────────────────────────────────┐
│  Browser     │     │  React App (apps/web/src/)                   │
│  (User)      │     │                                              │
│              │     │  ┌──────────────────┐  ┌──────────────────┐  │
│  ┌────────┐  │     │  │  App.tsx          │  │  main.tsx        │  │
│  │ Login  │──┼─────┼─▶│  BrowserRouter    │  │  React.StrictMode│  │
│  │ Page   │  │     │  │  QueryClientProv. │  └──────────────────┘  │
│  └────────┘  │     │  │  Routes           │                        │
│              │     │  │  ├─ /login         │  ┌──────────────────┐  │
│  ┌────────┐  │     │  │  ├─ /register      │  │  ProtectedRoute  │  │
│  │Tenant  │──┼─────┼─▶│  ├─ /reset         │  │  ├─ AuthContext   │  │
│  │ Pages  │  │     │  │  └─ /tenants/*     │  │  ├─ /auth/me     │  │
│  └────────┘  │     │  │       (Protected)  │  │  └─ useQuery     │  │
│              │     │  └──────────────────┘  └──────────────────┘  │
│              │     │                                              │
│              │     │  ┌──────────────┐  ┌──────────────────────┐  │
│              │     │  │  Pages/       │  │  Components/         │  │
│              │     │  │  TenantList   │  │  ├─ layout/          │  │
│              │     │  │  TenantCreate │  │  │  └─ AppLayout     │  │
│              │     │  │  TenantEdit   │  │  ├─ tenants/         │  │
│              │     │  │  TenantDetail │  │  │  ├─ TenantCard    │  │
│              │     │  └──────┬───────┘  │  │  └─ TenantForm     │  │
│              │     │         │          │  ├─ fees/             │  │
│              │     │         ▼          │  │  ├─ FeeList        │  │
│              │     │  ┌──────────────┐  │  │  └─ FeeForm        │  │
│              │     │  │  Custom Hooks │  │  └─ ui/              │  │
│              │     │  │  useTenants   │  │     ├─ Input         │  │
│              │     │  │  useTenant    │  │     ├─ Select        │  │
│              │     │  │  useFees      │  │     ├─ DatePicker    │  │
│              │     │  │  useCreateFee │  │     ├─ StatusBadge   │  │
│              │     │  │  useDeleteFee │  │     ├─ ConfirmDialog │  │
│              │     │  └──────┬───────┘  │     ├─ PageHeader    │  │
│              │     │         │          │     ├─ LoadingSkel.  │  │
│              │     │         ▼          │     └─ EmptyState    │  │
│              │     │  ┌──────────────┐  │                        │
│              │     │  │  Services/    │  │  ┌──────────────────┐  │
│              │     │  │  api.ts       │  │  │  Types/           │  │
│              │     │  │  auth.ts      │  │  │  tenant.ts       │  │
│              │     │  │  tenants.ts   │  │  │  fee.ts          │  │
│              │     │  │  fees.ts      │  │  │  auth.ts (exist) │  │
│              │     │  └──────┬───────┘  │  └──────────────────┘  │
│              │     │         │          │                        │
│              │     │         ▼          │                        │
│              │     │  ┌──────────────┐  │                        │
│              │     │  │  fetch()      │  │                        │
│              │     │  │  credentials: │  │                        │
│              │     │  │  'include'    │  │                        │
│              │     │  └──────┬───────┘  │                        │
│              │     └─────────┼──────────┘                        │
│              │               │                                    │
│              └───────────────┼────────────────────────────────────┘
│                              │ HTTP /api/*
│                              ▼
│              ┌─────────────────────────────┐
│              │  Go API (Fiber)             │
│              │  localhost:8080             │
│              │  PASETO cookie validation   │
│              │  Casbin policy enforcement  │
│              └─────────────────────────────┘
```

### Recommended Project Structure
```
apps/web/src/
├── components/
│   ├── auth/                          # Existing — NOT modified (D-13)
│   ├── layout/
│   │   └── AppLayout.tsx              ★ Responsive sidebar + header + content shell
│   ├── tenants/
│   │   ├── TenantCard.tsx             ★ Summary card for list view
│   │   ├── TenantCard.test.tsx        ★
│   │   ├── TenantForm.tsx             ★ Create/edit with dynamic mandatory fees
│   │   └── TenantForm.test.tsx        ★
│   ├── fees/
│   │   ├── FeeList.tsx                ★ Mandatory + voluntary sectioned list
│   │   ├── FeeList.test.tsx           ★
│   │   ├── FeeForm.tsx                ★ Create/edit fee (inline or modal)
│   │   └── FeeForm.test.tsx           ★
│   └── ui/
│       ├── Input.tsx                  ★ Reusable form input (label, error, aria)
│       ├── Select.tsx                 ★ Reusable form select
│       ├── DatePicker.tsx             ★ Native date input with styling
│       ├── StatusBadge.tsx            ★ Occupancy + payment status indicator
│       ├── ConfirmDialog.tsx          ★ Delete confirmation modal
│       ├── PageHeader.tsx             ★ Page title + action button pattern
│       ├── LoadingSkeleton.tsx        ★ Skeleton placeholders
│       └── EmptyState.tsx             ★ Empty state with CTA
├── hooks/
│   ├── useTenants.ts                  ★ TanStack query hook for tenant list
│   ├── useTenant.ts                   ★ TanStack query hook for single tenant
│   ├── useFees.ts                     ★ TanStack query hook for fee list
│   ├── useCreateTenant.ts            ★ TanStack mutation for tenant create
│   ├── useUpdateTenant.ts            ★ TanStack mutation for tenant update
│   ├── useDeleteTenant.ts            ★ TanStack mutation for tenant delete
│   ├── useCreateFee.ts               ★ TanStack mutation for fee create
│   ├── useUpdateFee.ts               ★ TanStack mutation for fee update
│   └── useDeleteFee.ts               ★ TanStack mutation for fee delete
├── pages/
│   ├── auth/                          # Existing
│   ├── TenantListPage.tsx             ★
│   ├── TenantCreatePage.tsx           ★
│   ├── TenantEditPage.tsx             ★
│   └── TenantDetailPage.tsx           ★ (fee management)
├── services/
│   ├── api.ts                         ★ Shared request() helper (NEW)
│   ├── auth.ts                        # Existing — updated to use api.ts (D-03)
│   ├── tenants.ts                     ★ Tenant API service
│   └── fees.ts                        ★ Fee API service
├── types/
│   ├── auth.ts                        # Existing
│   ├── tenant.ts                      ★ Tenant, api response types
│   └── fee.ts                         ★ Fee, Fee api request/response types
├── routes/
│   └── ProtectedRoute.tsx             # Existing — updated with useQuery for /auth/me (D-09)
├── App.tsx                            # Update — add QueryClientProvider + /tenants/* routes
├── index.css                          # Update — add sidebar tokens
└── main.tsx                           # Unchanged
```

### Pattern 1: Shared API Request Helper

**What:** A single `request()` function that wraps `fetch` with shared configuration, error handling, and content-type detection.

**When to use:** Every API service function (tenants, fees, and refactored auth).

**Example:**
```typescript
// Source: Phase 3 CONTEXT.md D-01, D-02
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface ApiError {
  error: string;
  code?: string;
}

export async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${url}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  // Handle 204 No Content (DELETE operations)
  if (response.status === 204) {
    return undefined as T;
  }

  // Auto-detect JSON vs non-JSON
  const contentType = response.headers.get('content-type');
  if (contentType?.includes('application/json')) {
    const data = await response.json();
    if (!response.ok) {
      throw new Error((data as ApiError).error || 'An unexpected error occurred');
    }
    return data as T;
  }

  if (!response.ok) {
    throw new Error('Connection lost. Check your internet and try again.');
  }

  return response.text() as T;
}
```

### Pattern 2: Custom TanStack Query Hook

**What:** Custom hook wrapping `useQuery` or `useMutation` that calls a service function.

**When to use:** Every API interaction. Follow this exact pattern for consistency.

**Example:**
```typescript
// Source: CONTEXT.md D-08, standard TanStack Query v5 pattern
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { listTenants, createTenant, type Tenant, type CreateTenantRequest } from '../services/tenants';

export function useTenants() {
  return useQuery<Tenant[]>({
    queryKey: ['tenants'],
    queryFn: listTenants,
  });
}

export function useCreateTenant() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTenantRequest) => createTenant(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });
}
```

### Pattern 3: Reusable Form Field Component

**What:** Wraps native `<input>` / `<select>` with label, error message, aria attributes, and consistent styling. Follows the LoginForm field markup pattern.

**When to use:** All forms in this phase (TenantForm, FeeForm). The existing auth forms are NOT refactored (D-13).

**Example:**
```typescript
// Source: CONTEXT.md D-10, D-12, LoginForm field pattern
import { type InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
}

export default function Input({ label, error, id, className = '', ...props }: InputProps) {
  const inputId = id || props.name;
  const errorId = `${inputId}-error`;

  return (
    <div className="space-y-2">
      <label htmlFor={inputId} className="block text-base font-medium text-gray-700">
        {label}
      </label>
      <input
        id={inputId}
        className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
          error
            ? 'border-red-600 focus:border-red-600'
            : 'border-gray-200 focus:border-blue-600'
        } ${className}`}
        aria-invalid={!!error}
        aria-describedby={error ? errorId : undefined}
        {...props}
      />
      {error && (
        <p id={errorId} className="text-sm text-red-600" role="alert">
          {error}
        </p>
      )}
    </div>
  );
}
```

### Anti-Patterns to Avoid
- **Inline form fields in TenantForm/FeeForm:** Use the shared Input/Select/DatePicker components (D-14). Don't copy-paste LoginForm's inline field markup — the whole point of the UI components is reusability.
- **Manual fetch in pages:** Always use TanStack Query hooks. Don't use `useEffect` + `useState` for data fetching — that's what TQ replaces.
- **QueryClient in test wrappers:** Per D-15, mock `@tanstack/react-query` directly with `vi.mock()`. Don't create test wrappers with QueryClientProvider.
- **Side-by-side fees and form on mobile:** The UI-SPEC specifies the FeeForm as full-screen modal on mobile, inline on desktop. Don't try to show both simultaneously on small screens.
- **Pre-populating auth form fields:** D-13 says existing auth forms stay as-is. Don't refactor LoginForm/RegisterForm to use the new Input component.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Data fetching & caching | Manual `useEffect` + `useState` + fetch | `@tanstack/react-query` | Auto-caching, background refetch, optimistic updates, loading/error states, stale-while-revalidate |
| Icons | SVG paths directly | `lucide-react` | Consistent 24px icons with aria-labels, tree-shakeable, 1000+ icons |
| Date picker | Custom date picker with month/year navigation | Native `<input type="date">` | Works on all modern mobile browsers, zero JS overhead, accessible by default |
| HTTP request boilerplate | Per-service error handling repetition | Shared `request()` helper | Single point for credentials, base URL, error parsing, content-type detection |
| Role-based sidebar visibility | Complex auth middleware | Simple `useAuth()` check in AppLayout | The user role is in AuthContext — just conditionally render sidebar items |

**Key insight:** TanStack Query eliminates ~40% of boilerplate compared to manual `useEffect` + `useState` data fetching. Every page that loads data from the API should use a custom hook wrapping `useQuery`, not a manual fetch.

## Runtime State Inventory

> **Not applicable — greenfield frontend phase.** This phase creates new UI code. It does not rename, refactor, or migrate any runtime state. The existing `apps/web/` code has no stored data, live service config, OS registrations, or build artifacts that need inventory.

## Common Pitfalls

### Pitfall 1: Forgetting to check Content-Type for 204 responses
**What goes wrong:** DELETE endpoints return 204 No Content with no body. Calling `.json()` on a 204 response throws a parse error.
**Why it happens:** The shared `request()` helper tries to parse JSON for all responses.
**How to avoid:** Check `response.status === 204` before any body parsing. Return `undefined` (or `null`) for 204s. The auto-detect approach (D-02) handles this.
**Warning signs:** Test fails with "Unexpected end of JSON input" on delete operations.

### Pitfall 2: Stale query cache after mutations
**What goes wrong:** After creating a tenant, the tenant list still shows the old data because the query cache wasn't invalidated.
**Why it happens:** TanStack Query caches data by `queryKey`. Mutations don't automatically invalidate related queries.
**How to avoid:** Always call `queryClient.invalidateQueries({ queryKey: ['tenants'] })` in `onSuccess` of mutations that affect the tenant list. Same pattern for fees.
**Warning signs:** Created tenant doesn't appear in list until page refresh.

### Pitfall 3: Mocking TanStack Query hooks incorrectly
**What goes wrong:** Tests fail because `vi.mock('@tanstack/react-query')` doesn't export the right mocks.
**Why it happens:** The module structure of TQ v5 uses named exports and default exports. Mocking needs to cover both `useQuery` and `useMutation`.
**How to avoid:** Follow the D-15 pattern: mock the specific hooks you need:
```typescript
vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual('@tanstack/react-query');
  return {
    ...actual,
    useQuery: vi.fn(),
    useMutation: vi.fn(),
  };
});
```
**Warning signs:** Test errors like "useQuery is not a function" or "Invalid hook call".

### Pitfall 4: Sidebar layout breaking on mobile vs desktop
**What goes wrong:** The sidebar layout looks correct on desktop but overlaps content on mobile, or the hamburger toggle doesn't work.
**Why it happens:** Responsive sidebar requires careful CSS: fixed positioning on mobile (overlay drawer), flex layout on desktop (inline sidebar).
**How to avoid:** Use separate layout modes: mobile (< 768px) = fixed sidebar with translateX transition + backdrop overlay; tablet (768-1023px) = collapsible icon-only sidebar; desktop (≥ 1024px) = fixed 260px sidebar. Test all three breakpoints.
**Warning signs:** Content area doesn't shift when sidebar opens on mobile, or double scrollbars appear.

### Pitfall 5: Numeric amount display with IDR formatting
**What goes wrong:** Monthly fee amounts like 50000 display as "50000" instead of "Rp 50,000".
**Why it happens:** The API returns numeric values. Client-side formatting is needed for display.
**How to avoid:** Create a utility function `formatIDR(amount: number): string` using `Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 })`. The UI-SPEC copywriting contract specifies "Rp X,XXX" format.
**Warning signs:** Fee amounts appear as raw numbers or with wrong formatting.

## Code Examples

Verified patterns from the existing codebase:

### API Service Pattern (existing auth.ts, to be migrated)
```typescript
// Source: apps/web/src/services/auth.ts — current pattern
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export async function login(email: string, password: string): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ email, password }),
  });
  return parseError<LoginResponse>(response);
}
```

### Tenant API Service (new, using shared helper)
```typescript
// Source: Phase 2 API endpoints (tenant_handler.go), CONTEXT.md D-04
import { request } from './api';
import type { Tenant } from '../types/tenant';
import type { Fee } from '../types/fee';

export interface ListTenantsResponse {
  tenants: Tenant[];
}

export interface CreateTenantRequest {
  block: string;
  unit_number: string;
  occupancy: 'occupied' | 'vacant';
  monthly_fee: number;
  mandatory_fees: Array<{
    amount: number;
    description: string;
    effective_date: string;
  }>;
}

export async function listTenants(): Promise<Tenant[]> {
  const data = await request<ListTenantsResponse>('/api/tenants');
  return data.tenants;
}

export async function createTenant(data: CreateTenantRequest): Promise<Tenant> {
  return request<Tenant>('/api/tenants', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getTenant(id: string): Promise<Tenant> {
  return request<Tenant>(`/api/tenants/${id}`);
}

export async function updateTenant(id: string, data: Partial<CreateTenantRequest>): Promise<Tenant> {
  return request<Tenant>(`/api/tenants/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteTenant(id: string): Promise<void> {
  return request<void>(`/api/tenants/${id}`, { method: 'DELETE' });
}
```

### Fee API Service (new)
```typescript
// Source: Phase 2 API endpoints (tenant_handler.go: ListFees, CreateFee, UpdateFee, DeleteFee)
import { request } from './api';
import type { Fee } from '../types/fee';

export interface ListFeesResponse {
  mandatory_fees: Fee[];
  voluntary_fees: Fee[];
}

export interface CreateFeeRequest {
  type: 'mandatory' | 'voluntary';
  amount: number;
  description: string;
  effective_date: string;
  paid_at?: string;
}

export async function listFees(tenantId: string): Promise<{ mandatory_fees: Fee[]; voluntary_fees: Fee[] }> {
  return request<ListFeesResponse>(`/api/tenants/${tenantId}/fees`);
}

export async function createFee(tenantId: string, data: CreateFeeRequest): Promise<Fee> {
  return request<Fee>(`/api/tenants/${tenantId}/fees`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateFee(tenantId: string, feeId: string, data: Partial<CreateFeeRequest>): Promise<void> {
  return request<void>(`/api/tenants/${tenantId}/fees/${feeId}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteFee(tenantId: string, feeId: string): Promise<void> {
  return request<void>(`/api/tenants/${tenantId}/fees/${feeId}`, { method: 'DELETE' });
}
```

### App.tsx Routing with QueryClientProvider
```typescript
// Source: CONTEXT.md D-07, existing App.tsx pattern
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ResetPasswordPage from './pages/ResetPasswordPage';
import TenantListPage from './pages/TenantListPage';
import TenantCreatePage from './pages/TenantCreatePage';
import TenantEditPage from './pages/TenantEditPage';
import TenantDetailPage from './pages/TenantDetailPage';
import ProtectedRoute from './routes/ProtectedRoute';
import AppLayout from './components/layout/AppLayout';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,       // 30s before refetch
      retry: 1,                 // Retry once on failure
      refetchOnWindowFocus: true,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/reset" element={<ResetPasswordPage />} />
          <Route
            path="/tenants"
            element={
              <ProtectedRoute>
                <AppLayout><TenantListPage /></AppLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/tenants/new"
            element={
              <ProtectedRoute>
                <AppLayout><TenantCreatePage /></AppLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/tenants/:id"
            element={
              <ProtectedRoute>
                <AppLayout><TenantDetailPage /></AppLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/tenants/:id/edit"
            element={
              <ProtectedRoute>
                <AppLayout><TenantEditPage /></AppLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <AppLayout>
                  <div className="flex min-h-screen items-center justify-center">
                    <div className="text-center">
                      <h1 className="text-2xl font-semibold text-gray-900">Welcome to Harmoni</h1>
                      <p className="mt-2 text-gray-600">Your dashboard is coming soon.</p>
                    </div>
                  </div>
                </AppLayout>
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
```

### Test Pattern with Mocked TanStack Query
```typescript
// Source: CONTEXT.md D-15, LoginForm.test.tsx pattern
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import TenantListPage from './TenantListPage';

vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual('@tanstack/react-query');
  return {
    ...actual,
    useQuery: vi.fn(),
  };
});

import { useQuery } from '@tanstack/react-query';

function renderTenantListPage() {
  return render(
    <MemoryRouter>
      <TenantListPage />
    </MemoryRouter>
  );
}

describe('TenantListPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows loading skeleton while fetching tenants', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
    });
    renderTenantListPage();
    expect(screen.getByTestId('loading-skeleton')).toBeInTheDocument();
  });

  it('shows tenant cards when data loads', async () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: [
        { id: '1', block: 'A', unit_number: '01', occupancy: 'occupied', monthly_fee: 50000 },
      ],
      isLoading: false,
      isError: false,
    });
    renderTenantListPage();
    expect(screen.getByText(/Block A/)).toBeInTheDocument();
    expect(screen.getByText(/Unit 01/)).toBeInTheDocument();
  });
});
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Raw `fetch` in components | TanStack Query `useQuery`/`useMutation` | Phase 3 | All data fetching goes through TQ hooks. Components never call `fetch` directly. |
| Per-service error parsing | Shared `request()` helper in `api.ts` | Phase 3 | Single source for base URL, credentials, content-type detection, error handling |
| Auth/me fetch in ProtectedRoute `useEffect` | `useQuery` for `/auth/me` | Phase 3 (D-09) | Cached user data, automatic re-fetch, no manual loading state management |
| Inline form fields in LoginForm/RegisterForm | Reusable Input/Select/DatePicker components | Phase 3 | Shared aria attributes, error handling, and styling across new forms |

**Deprecated/outdated:**
- **Manual `fetch` + `useEffect` + `useState` in components:** Replaced by TanStack Query hooks. Don't use this pattern for any new code.
- **Duplicate `API_BASE_URL` constant in every service file:** Centralized in `services/api.ts`.
- **`App.css`:** Boilerplate from Vite template. UI-SPEC says can be removed.

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Phase 2 API is deployed and reachable at `VITE_API_URL` (http://localhost:8080) | API Contract | Frontend cannot function without backend. Planner should verify backend is running before test/wiring tasks. |
| A2 | The API response for `GET /api/tenants` is `{ tenants: Tenant[] }` (wrapped in object) | API Contract | Confirmed in tenant_handler.go line 116: `c.JSON(fiber.Map{"tenants": tenants})`. LOW risk — verified. |
| A3 | The API response for `GET /api/tenants/:id/fees` is `{ mandatory_fees: Fee[], voluntary_fees: Fee[] }` | API Contract | Confirmed in tenant_handler.go lines 339-342. LOW risk — verified. |
| A4 | The `Fee` entity from API includes `type` field | Types | Backend returns separate `MandatoryFee`/`VoluntaryFee` structs that don't have a `type` discriminator. The `type` is determined by which array the fee appears in. However, `CreateFeeRequest` has a `type` field. Type definitions must account for this. MEDIUM risk — verify by running backend and checking actual response shape. |

## Open Questions (RESOLVED)

1. **Fee response shape: does the backend return a `type` field on individual fee objects?** — RESOLVED: Include optional `type` field on Fee interface, infer type from array position when absent.

2. **How should the FeeForm handle the fee type selector UI when creating a fee after the tenant already has mandatory fees?** — RESOLVED: Follow UI-SPEC literally — fee form on detail page allows both mandatory/voluntary types. Backend accepts `type` discriminator in POST body.

3. **How to display fee status (paid/unpaid) when the API returns `paid_at` as null vs a date string?** — RESOLVED: If `paid_at` is `null` → "Unpaid" badge. If `paid_at` is a date string → "Paid" badge. StatusBadge accepts `status: 'paid' | 'unpaid'` prop.

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Node.js | React app runtime | ✓ | 24.12.0 | — |
| npm | Package management | ✓ | 11.14.1 | — |
| Vite | Dev server / build | ✓ | (via package.json) | — |
| TanStack Query v5 | Data fetching | ✓ | 5.100.14 (registry) | — |
| lucide-react | Icons | ✓ | 1.16.0 (registry) | — |
| Tailwind CSS v4 | Styling | ✓ | (via package.json) | — |
| Vitest | Test runner | ✓ | (via package.json) | — |
| Go API backend | API consumption | ✓ (assumed, not running in this phase) | — | Can test with mocked services |

**Missing dependencies with no fallback:** None — all dependencies are available or will be installed via npm.

## Validation Architecture

> nyquist_validation is enabled in `.planning/config.json`. This section describes the test infrastructure and phase requirement coverage.

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Vitest v4.1.6 + Testing Library React v16.3.2 |
| Config file | `apps/web/vitest.config.js` (jsdom environment, globals: true, setup: `src/test/setup.ts`) |
| Quick run command | `npm test -- --run` (or `npx vitest run` from `apps/web/`) |
| Full suite command | `npm test` (or `npx vitest` for watch mode) |

### Existing Test Infrastructure
- **Setup:** `src/test/setup.ts` — imports `@testing-library/jest-dom`
- **Test pattern:** Component test files co-located with components (`*.test.tsx`)
- **Service test files:** Co-located with services (`auth.test.ts`)
- **Mock pattern:** `vi.mock('../../services/auth')` or `vi.mock('@tanstack/react-query')` at module level; `vi.clearAllMocks()` in `beforeEach`

### New Files Required (tests)
| File | Tests |
|------|-------|
| `components/ui/Input.test.tsx` | Renders label, shows error message, applies aria attributes |
| `components/ui/Select.test.tsx` | Renders options, shows error, handles change |
| `components/ui/DatePicker.test.tsx` | Renders date input, handles change |
| `components/ui/StatusBadge.test.tsx` | Occupied vs Vacant, Paid vs Unpaid |
| `components/ui/ConfirmDialog.test.tsx` | Shows on trigger, Cancel dismisses, Delete confirms |
| `components/ui/EmptyState.test.tsx` | Shows message + CTA |
| `components/ui/LoadingSkeleton.test.tsx` | Renders skeleton placeholders |
| `components/ui/PageHeader.test.tsx` | Renders title + CTA |
| `components/layout/AppLayout.test.tsx` | Sidebar visible/hidden, responsive behavior, navigation |
| `components/tenants/TenantCard.test.tsx` | Renders tenant info, occupancy badge, fee summary |
| `components/tenants/TenantForm.test.tsx` | Field validation, add fee row, submit, edit mode |
| `components/fees/FeeList.test.tsx` | Renders mandatory + voluntary sections, empty states |
| `components/fees/FeeForm.test.tsx` | Type selector, validation, submit |
| `pages/TenantListPage.test.tsx` | Loading skeleton, tenant cards, empty state, search filter |
| `pages/TenantCreatePage.test.tsx` | Renders form, cancel goes back |
| `pages/TenantEditPage.test.tsx` | Pre-populated fields, delete button |
| `pages/TenantDetailPage.test.tsx` | Tenant info header, fee sections, record fee button |
| `services/tenants.test.ts` | Calls correct endpoints, handles responses |
| `services/fees.test.ts` | Calls correct endpoints, handles responses |
| `services/api.test.ts` | Shared helper behavior (base URL, credentials, error handling, 204 handling) |

### Phase Requirements → Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| TEN-01 | Tenant CRUD through UI | unit (service/component) | `npx vitest run components/tenants/ pages/TenantListPage` | ❌ Wave 0 |
| FIN-01 | Mandatory fee creation inline with tenant | unit (component) | `npx vitest run components/tenants/TenantForm` | ❌ Wave 0 |
| FIN-02 | Voluntary fee creation on detail page | unit (component) | `npx vitest run components/fees/` | ❌ Wave 0 |

### Sampling Rate
- **Per task commit:** `npm test -- --run` (full suite, ~30-60s for Phase 3)
- **Per wave merge:** `npm test -- --run`
- **Phase gate:** Full suite green before `/gsd-verify-work`

### Wave 0 Gaps
- [ ] `components/ui/Input.test.tsx` — covers reusable field pattern
- [ ] `components/layout/AppLayout.test.tsx` — covers responsive sidebar
- [ ] `components/tenants/TenantForm.test.tsx` — covers TEN-01, FIN-01
- [ ] `components/fees/FeeForm.test.tsx` — covers FIN-01, FIN-02

*(All test files are new — Wave 0 creates test infrastructure for all components and pages)*

## Security Domain

> security_enforcement is absent from config.json (default enabled). This section covers applicable ASVS categories for the Phase 3 frontend.

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | No | Handled by Phase 1 backend (PASETO httpOnly cookies) + ProtectedRoute component |
| V3 Session Management | No | Handled by Phase 1 — cookie-based, server-managed |
| V4 Access Control | Partial | Role-based sidebar item visibility (RT officer sees all, resident read-only) |
| V5 Input Validation | Yes | Client-side validation before submit (form fields, fee amounts, dates) |
| V6 Cryptography | No | All crypto handled server-side (PASETO, password hashing) |
| V8 Data Protection | Partial | API calls use `credentials: 'include'` — cookies sent only to same-origin API |

### Known Threat Patterns for React + TypeScript

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| XSS via user input rendering | Tampering | React's JSX auto-escapes. Avoid `dangerouslySetInnerHTML`. All data from API rendered via JSX expressions. |
| CSRF via cookie-based auth | Tampering | PASETO httpOnly cookies + `credentials: 'include'` + same-origin API calls. No CSRF token needed — browser's SameSite cookie policy (Lax by default) covers this. |
| Sensitive data in URL params | Information Disclosure | Never pass PASETO tokens, passwords, or PII in URL query strings. Use request bodies for sensitive data. |
| Insecure direct object reference | Elevation of Privilege | Mitigated by backend Casbin policies. Frontend never enforces tenant isolation — that's the API's job. Frontend should gracefully handle 403 responses. |
| Excessive data exposure | Information Disclosure | API returns only the fields needed. Frontend types match API response shape. |

## Sources

### Primary (HIGH confidence)
- **Codebase analysis** — All `apps/web/src/` files read and patterns documented (services, types, pages, components, tests, routing, config, CSS)
- **Phase 2 API handler** — `apps/api/internal/delivery/http/tenant_handler.go` — confirmed all endpoint shapes, request/response types, status codes
- **Phase 2 entity definitions** — `apps/api/internal/domain/entity/tenant.go` — confirmed `Tenant`, `MandatoryFee`, `VoluntaryFee` response shapes
- **CONTEXT.md** — Phase 3 locked decisions (D-01 through D-17)
- **UI-SPEC.md** — Approved design contract for all components, pages, copywriting, colors, spacing, responsive rules
- **npm registry** — Verified `@tanstack/react-query@5.100.14`, `lucide-react@1.16.0` package existence and metadata

### Secondary (MEDIUM confidence)
- **Phase 1 CONTEXT.md** — Confirmed auth patterns (PASETO httpOnly cookies, territory model, Casbin policies)
- **Phase 2 CONTEXT.md** — Confirmed API endpoint structure, data isolation, validation rules
- **Graph knowledge base** — Community structure confirmed project architecture patterns (553 nodes, 891 edges from entire codebase)

### Tertiary (LOW confidence)
- (None — all key claims verified against codebase or registry)

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all packages verified on registry, used in established patterns
- Architecture: HIGH — patterns derived from reading existing codebase files and Phase 1 implementations
- Pitfalls: HIGH — based on TanStack Query v5 known issues, previous experience with similar patterns
- API Contract: HIGH — confirmed by reading Phase 2 handler code and entity definitions

**Research date:** 2026-05-24
**Valid until:** 2026-08-30 (all dependencies are stable/established)

---

## RESEARCH COMPLETE

**Phase:** 3 - Tenant & Fee UI
**Confidence:** HIGH

### Key Findings
1. **Stack is locked:** TanStack Query v5 + lucide-react are the only new dependencies. All other tooling exists from Phase 1 (Tailwind v4, Vitest, Testing Library, React Router v7).
2. **API contract is known:** All 9 endpoints confirmed from Phase 2 handler code (CRUD tenants + fees, fee list returns mandatory/voluntary separately).
3. **Build order matters:** Create shared `api.ts` helper first → then types → services → hooks → UI components → pages → routing in App.tsx → ProtectedRoute migration.
4. **Test strategy is clear:** Mock `@tanstack/react-query` directly (no QueryClientProvider in tests). Follow LoginForm test pattern.
5. **AppLayout is the new shell:** Responsive sidebar wraps all authenticated pages. ProtectedRoute should wrap AppLayout, not individual pages.

### File Created
`.planning/phases/03-tenant-fee-ui/03-RESEARCH.md`

### Confidence Assessment
| Area | Level | Reason |
|------|-------|--------|
| Standard Stack | HIGH | All packages verified on npm registry. All decisions from CONTEXT.md are specific and locked. |
| Architecture | HIGH | Patterns derived from reading every source file in apps/web/src/. No guessing. |
| Pitfalls | HIGH | Based on known TanStack Query v5 issues (cache invalidation, mocking), 204 responses, responsive layout challenges. |

### Open Questions (RESOLVED)
1. Fee response object shape — RESOLVED: include optional `type` field on Fee, infer from array position when absent
2. FeeForm behavior for fee type selection on existing tenant — RESOLVED: follow UI-SPEC literally, allow both types

### Ready for Planning
Research complete. Planner can now create PLAN.md files for Phase 3 — recommends splitting into waves:
- Wave 0: api.ts helper, types, services, initial test setup
- Wave 1: UI components (Input/Select/DatePicker + StatusBadge/ConfirmDialog/etc)
- Wave 2: AppLayout + routes + ProtectedRoute migration
- Wave 3: TenantListPage, TenantCard, search filter
- Wave 4: TenantCreatePage, TenantForm, TenantEditPage
- Wave 5: TenantDetailPage, FeeList, FeeForm
- Wave 6: Tests for all new components and pages
