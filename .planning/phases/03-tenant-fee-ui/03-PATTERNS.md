# Phase 3: Tenant & Fee UI — Pattern Map

**Mapped:** 2026-05-24
**Files analyzed:** 43 (17 new source + 6 modified + 20 test)
**Analogs found:** 27 / 27

## File Classification

| New/Modified File | Role | Data Flow | Closest Analog | Match Quality |
|---|---|---|---|---|
| `services/api.ts` | utility | request-response | `services/auth.ts` (parseError pattern) | role-match |
| `services/tenants.ts` | service | CRUD | `services/auth.ts` | role-match |
| `services/fees.ts` | service | CRUD | `services/auth.ts` | role-match |
| `services/auth.ts` (modified) | service | request-response | `services/auth.ts` (self) | exact |
| `types/tenant.ts` | model | — | `types/auth.ts` | exact |
| `types/fee.ts` | model | — | `types/auth.ts` | exact |
| `components/ui/Input.tsx` | component | request-response | `LoginForm.tsx` field pattern | partial |
| `components/ui/Select.tsx` | component | request-response | `LoginForm.tsx` field pattern | partial |
| `components/ui/DatePicker.tsx` | component | request-response | `LoginForm.tsx` field pattern | partial |
| `components/ui/StatusBadge.tsx` | component | — | No analog | — |
| `components/ui/ConfirmDialog.tsx` | component | event-driven | No analog | — |
| `components/ui/PageHeader.tsx` | component | — | No analog | — |
| `components/ui/LoadingSkeleton.tsx` | component | — | No analog | — |
| `components/ui/EmptyState.tsx` | component | — | No analog | — |
| `components/layout/AppLayout.tsx` | component | — | No analog | — |
| `components/tenants/TenantCard.tsx` | component | — | No analog | — |
| `components/tenants/TenantForm.tsx` | component | CRUD | `LoginForm.tsx` | role-match |
| `components/fees/FeeList.tsx` | component | CRUD | No analog | — |
| `components/fees/FeeForm.tsx` | component | CRUD | `LoginForm.tsx` | role-match |
| `hooks/useTenants.ts` | hook | CRUD | No analog (TanStack Query new) | — |
| `hooks/useTenant.ts` | hook | CRUD | No analog | — |
| `hooks/useFees.ts` | hook | CRUD | No analog | — |
| `hooks/useCreateTenant.ts` | hook | CRUD | No analog | — |
| `hooks/useUpdateTenant.ts` | hook | CRUD | No analog | — |
| `hooks/useDeleteTenant.ts` | hook | CRUD | No analog | — |
| `hooks/useCreateFee.ts` | hook | CRUD | No analog | — |
| `hooks/useUpdateFee.ts` | hook | CRUD | No analog | — |
| `hooks/useDeleteFee.ts` | hook | CRUD | No analog | — |
| `pages/TenantListPage.tsx` | page | CRUD | `LoginPage.tsx` | role-match |
| `pages/TenantCreatePage.tsx` | page | CRUD | `LoginPage.tsx` | role-match |
| `pages/TenantEditPage.tsx` | page | CRUD | `LoginPage.tsx` | role-match |
| `pages/TenantDetailPage.tsx` | page | CRUD | `LoginPage.tsx` | role-match |
| `App.tsx` (modified) | config | — | `App.tsx` (existing) | exact |
| `index.css` (modified) | config | — | `index.css` (existing) | exact |
| `routes/ProtectedRoute.tsx` (modified) | middleware | request-response | `ProtectedRoute.tsx` (existing) | exact |

## Pattern Assignments

### `services/api.ts` (utility, request-response)

**Analog:** `services/auth.ts` — existing `parseError<T>()` + `fetch` wrapper pattern

**Imports pattern** (lines 1):
```typescript
// Source: services/auth.ts line 1 — type-only import pattern
import type { ApiResponse, LoginResponse, RegisterRequest } from '../types/auth';
```

**Core fetch pattern (to refactor into shared helper)** (lines 5-17):
```typescript
// Source: services/auth.ts lines 5-17 — parseError + fetch
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

function parseError<T>(response: Response): Promise<T> {
  if (!response.ok) {
    return response.json().then(
      (data) => {
        throw new Error(data.error || 'An unexpected error occurred');
      },
      () => {
        throw new Error('Connection lost. Check your internet and try again.');
      }
    );
  }
  return response.json();
}
```

**API call pattern (to be built on shared helper)** (lines 19-28):
```typescript
// Source: services/auth.ts lines 19-28 — fetch with credentials: 'include'
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

**Expected new pattern (based on CONTEXT.md D-01, D-02):**
```typescript
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

---

### `services/tenants.ts` (service, CRUD)

**Analog:** `services/auth.ts` — existing service file pattern

**Imports pattern:**
```typescript
// Source: services/auth.ts line 1 — type-only import
import type { ApiResponse, LoginResponse, RegisterRequest } from '../types/auth';

// Expected new pattern:
import { request } from './api';
import type { Tenant } from '../types/tenant';
import type { Fee } from '../types/fee';
```

**CRUD endpoint pattern (expected):**
```typescript
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

export async function getTenant(id: string): Promise<Tenant> {
  return request<Tenant>(`/api/tenants/${id}`);
}

export async function createTenant(data: CreateTenantRequest): Promise<Tenant> {
  return request<Tenant>('/api/tenants', {
    method: 'POST',
    body: JSON.stringify(data),
  });
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

---

### `services/fees.ts` (service, CRUD)

**Analog:** `services/auth.ts` — same role (service), same data flow (CRUD)

**Imports pattern (expected):**
```typescript
import { request } from './api';
import type { Fee } from '../types/fee';
```

**Fee CRUD pattern (expected):**
```typescript
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

---

### `services/auth.ts` (modified service, request-response)

**Analog:** Existing `services/auth.ts` — to be refactored to use shared `request()` helper

**Refactored pattern (each function simplified):**
```typescript
// Current (lines 19-28):
export async function login(email: string, password: string): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ email, password }),
  });
  return parseError<LoginResponse>(response);
}

// Expected refactored pattern:
import { request } from './api';

export async function login(email: string, password: string): Promise<LoginResponse> {
  return request<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}
```

---

### `types/tenant.ts` (model)

**Analog:** `types/auth.ts` — exact role match (model)

**Imports pattern** (lines 1-7):
```typescript
// Source: types/auth.ts — interface-only file, no imports
export interface User {
  id: string;
  email: string;
  full_name?: string;
  role?: string;
  territory_id?: string;
}

export interface LoginResponse {
  user: User;
}
```

**Expected new pattern:**
```typescript
// types/tenant.ts — same structure as types/auth.ts
export interface Tenant {
  id: string;
  block: string;
  unit_number: string;
  occupancy: 'occupied' | 'vacant';
  monthly_fee: number;
  territory_id: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTenantRequest {
  block: string;
  unit_number: string;
  occupancy: 'occupied' | 'vacant';
  monthly_fee: number;
  mandatory_fees: CreateFeeRequest[];
}

export interface CreateFeeRequest {
  type: 'mandatory' | 'voluntary';
  amount: number;
  description: string;
  effective_date: string;
  paid_at?: string;
}
```

---

### `types/fee.ts` (model)

**Analog:** `types/auth.ts` — exact role match (model)

**Expected pattern:**
```typescript
export interface Fee {
  id: string;
  tenant_id: string;
  type: 'mandatory' | 'voluntary';
  amount: number;
  description: string;
  effective_date: string;
  paid_at: string | null;
  created_at: string;
}
```

---

### `components/ui/Input.tsx` (component, request-response)

**Analog:** `components/auth/LoginForm.tsx` lines 60-83 — field markup pattern

**Imports pattern** (line 1):
```typescript
// Source: LoginForm.tsx line 1
import { useState, type FormEvent, type ChangeEvent } from 'react';
```

**Expected pattern (derived from LoginForm field markup):**
```typescript
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

**Key styling tokens from LoginForm** (lines 69-73, 95-99):
```typescript
// Source: LoginForm.tsx lines 69-73 — input element styling
className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
  fieldErrors.email
    ? 'border-red-600 focus:border-red-600'
    : 'border-gray-200 focus:border-blue-600'
}`}

// Source: LoginForm.tsx lines 61-62, 79-83 — label + error message
<label htmlFor="email" className="block text-base font-medium text-gray-700">Email</label>
<p id="email-error" className="text-sm text-red-600" role="alert">{fieldErrors.email}</p>
```

---

### `components/ui/Select.tsx` (component, request-response)

**Analog:** `LoginForm.tsx` lines 60-83 — field markup pattern (same role)

**Expected pattern (extends Input field markup):**
```typescript
import { type SelectHTMLAttributes } from 'react';

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label: string;
  error?: string;
  options: Array<{ value: string; label: string }>;
}

export default function Select({ label, error, id, options, className = '', ...props }: SelectProps) {
  const inputId = id || props.name;
  const errorId = `${inputId}-error`;

  return (
    <div className="space-y-2">
      <label htmlFor={inputId} className="block text-base font-medium text-gray-700">
        {label}
      </label>
      <select
        id={inputId}
        className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
          error
            ? 'border-red-600 focus:border-red-600'
            : 'border-gray-200 focus:border-blue-600'
        } ${className}`}
        aria-invalid={!!error}
        aria-describedby={error ? errorId : undefined}
        {...props}
      >
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
      {error && (
        <p id={errorId} className="text-sm text-red-600" role="alert">
          {error}
        </p>
      )}
    </div>
  );
}
```

---

### `components/ui/DatePicker.tsx` (component, request-response)

**Analog:** `LoginForm.tsx` lines 60-83 — field markup pattern (same role)

**Expected pattern (native `<input type="date">` wrapper):**
```typescript
import { type InputHTMLAttributes } from 'react';

interface DatePickerProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  label: string;
  error?: string;
}

export default function DatePicker({ label, error, id, className = '', ...props }: DatePickerProps) {
  const inputId = id || props.name;
  const errorId = `${inputId}-error`;

  return (
    <div className="space-y-2">
      <label htmlFor={inputId} className="block text-base font-medium text-gray-700">
        {label}
      </label>
      <input
        id={inputId}
        type="date"
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

---

### `components/ui/StatusBadge.tsx` (component)

**No close analog — new component type.** See shared patterns below for creation guidance.

**Expected pattern:**
```typescript
interface StatusBadgeProps {
  variant: 'occupied' | 'vacant' | 'paid' | 'unpaid';
}

export default function StatusBadge({ variant }: StatusBadgeProps) {
  const styles: Record<string, string> = {
    occupied: 'bg-blue-100 text-blue-700',
    vacant: 'bg-amber-100 text-amber-800',
    paid: 'bg-emerald-100 text-emerald-800',
    unpaid: 'bg-red-100 text-red-800',
  };

  const labels: Record<string, string> = {
    occupied: 'Occupied',
    vacant: 'Vacant',
    paid: 'Paid',
    unpaid: 'Unpaid',
  };

  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-sm font-medium ${styles[variant]}`}>
      {labels[variant]}
    </span>
  );
}
```

---

### `components/ui/ConfirmDialog.tsx` (component, event-driven)

**No close analog — new component type.**

**Expected pattern (modal with dialog role, focus trap, Escape to close):**
```typescript
import { useEffect, useRef } from 'react';

interface ConfirmDialogProps {
  open: boolean;
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
  destructive?: boolean;
}

export default function ConfirmDialog({
  open, title, message, confirmLabel = 'Delete', cancelLabel = 'Cancel',
  onConfirm, onCancel, loading = false, destructive = true,
}: ConfirmDialogProps) {
  const cancelRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (open) cancelRef.current?.focus();
  }, [open]);

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape' && open) onCancel();
    }
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [open, onCancel]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-gray-900/50" onClick={onCancel}>
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="confirm-title"
        className="w-full max-w-sm rounded-lg bg-white p-6 shadow-lg mx-4"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 id="confirm-title" className="text-lg font-semibold text-gray-900">{title}</h2>
        <p className="mt-2 text-sm text-gray-600">{message}</p>
        <div className="mt-6 flex justify-end gap-3">
          <button
            ref={cancelRef}
            onClick={onCancel}
            disabled={loading}
            className="min-h-[44px] min-w-[96px] rounded-md border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            {cancelLabel}
          </button>
          <button
            onClick={onConfirm}
            disabled={loading}
            className={`min-h-[44px] min-w-[96px] rounded-md px-4 py-2 text-sm font-medium text-white ${
              destructive ? 'bg-red-600 hover:bg-red-700' : 'bg-blue-600 hover:bg-blue-700'
            } disabled:cursor-not-allowed disabled:opacity-50`}
          >
            {loading ? 'Deleting...' : confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
```

---

### `components/ui/PageHeader.tsx` (component)

**No close analog — new component type.**

**Expected pattern:**
```typescript
import type { ReactNode } from 'react';

interface PageHeaderProps {
  title: string;
  action?: ReactNode;
  backLink?: { to: string; label: string };
}

export default function PageHeader({ title, action, backLink }: PageHeaderProps) {
  return (
    <div className="mb-6">
      {backLink && (
        <a href={backLink.to} className="mb-2 inline-flex items-center text-sm text-blue-600 hover:text-blue-700">
          ← {backLink.label}
        </a>
      )}
      <div className="flex items-center justify-between">
        <h1 className="text-[24px] font-semibold text-gray-900">{title}</h1>
        {action && <div>{action}</div>}
      </div>
    </div>
  );
}
```

---

### `components/ui/LoadingSkeleton.tsx` (component)

**No close analog — new component type.**

**Expected pattern:**
```typescript
interface LoadingSkeletonProps {
  variant?: 'card' | 'list' | 'form';
  count?: number;
}

export default function LoadingSkeleton({ variant = 'card', count = 3 }: LoadingSkeletonProps) {
  if (variant === 'card') {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" data-testid="loading-skeleton">
        {Array.from({ length: count }).map((_, i) => (
          <div key={i} className="animate-pulse rounded-lg border border-gray-200 p-4">
            <div className="h-4 w-3/4 rounded bg-gray-200" />
            <div className="mt-3 h-3 w-1/2 rounded bg-gray-200" />
            <div className="mt-4 h-3 w-full rounded bg-gray-200" />
          </div>
        ))}
      </div>
    );
  }

  if (variant === 'list') {
    return (
      <div className="space-y-3" data-testid="loading-skeleton">
        {Array.from({ length: count }).map((_, i) => (
          <div key={i} className="animate-pulse rounded-lg border border-gray-200 p-4">
            <div className="h-4 w-1/2 rounded bg-gray-200" />
            <div className="mt-2 h-3 w-3/4 rounded bg-gray-200" />
          </div>
        ))}
      </div>
    );
  }

  return null;
}
```

---

### `components/ui/EmptyState.tsx` (component)

**No close analog — new component type.**

**Expected pattern:**
```typescript
import type { ReactNode } from 'react';

interface EmptyStateProps {
  heading: string;
  body: string;
  action?: ReactNode;
}

export default function EmptyState({ heading, body, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-gray-200 py-12 px-4 text-center">
      <h3 className="text-base font-semibold text-gray-900">{heading}</h3>
      <p className="mt-1 text-sm text-gray-500">{body}</p>
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}
```

---

### `components/layout/AppLayout.tsx` (component)

**No close analog — new shell component.** No existing layout component exists.

**Expected pattern (responsive sidebar structure):**
```typescript
import { useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { Building2, BarChart3, Settings, Menu, X } from 'lucide-react';
import { useAuth } from '../../routes/ProtectedRoute';

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { user } = useAuth();
  const isOfficer = user?.role === 'rt_officer' || user?.role === 'rw_officer';

  const navItems = [
    { to: '/tenants', label: 'Tenants', icon: Building2, visible: true },
    { to: '/reports', label: 'Reports', icon: BarChart3, visible: true },
    { to: '/settings', label: 'Settings', icon: Settings, visible: isOfficer },
  ];

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Mobile overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-40 bg-gray-900/50 lg:hidden" onClick={() => setSidebarOpen(false)} />
      )}

      {/* Sidebar */}
      <aside className={`fixed inset-y-0 left-0 z-50 w-[260px] transform bg-gray-100 transition-transform duration-200 lg:static lg:translate-x-0 ${
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      }`}>
        <div className="flex h-16 items-center justify-between px-4 border-b border-gray-200">
          <NavLink to="/dashboard" className="text-[28px] font-semibold text-gray-900">Harmoni</NavLink>
          <button onClick={() => setSidebarOpen(false)} className="min-h-[44px] min-w-[44px] lg:hidden" aria-label="Close sidebar">
            <X className="h-5 w-5" />
          </button>
        </div>
        <nav className="mt-4 space-y-1 px-3">
          {navItems.filter(i => i.visible).map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              onClick={() => setSidebarOpen(false)}
              className={({ isActive }) =>
                `flex min-h-[44px] items-center gap-3 rounded-md px-4 py-2.5 text-sm font-medium transition-colors ${
                  isActive ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-200'
                }`
              }
            >
              <item.icon className="h-5 w-5" />
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      {/* Main content */}
      <div className="flex flex-1 flex-col">
        <header className="flex h-16 items-center gap-4 border-b border-gray-200 bg-white px-4 lg:px-6">
          <button onClick={() => setSidebarOpen(true)} className="min-h-[44px] min-w-[44px] lg:hidden" aria-label="Open sidebar">
            <Menu className="h-5 w-5" />
          </button>
          {/* Page title can be passed via context or rendered inline in pages */}
        </header>
        <main className="flex-1 p-4 lg:p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
```

---

### `components/tenants/TenantCard.tsx` (component)

**No close analog — new component type.**

**Expected pattern:**
```typescript
import { Link } from 'react-router-dom';
import StatusBadge from '../ui/StatusBadge';
import type { Tenant } from '../../types/tenant';

interface TenantCardProps {
  tenant: Tenant;
  mandatoryFeeCount: number;
  voluntaryFeeCount: number;
}

export default function TenantCard({ tenant, mandatoryFeeCount, voluntaryFeeCount }: TenantCardProps) {
  const hasFees = mandatoryFeeCount > 0 || voluntaryFeeCount > 0;
  const monthlyFee = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(tenant.monthly_fee);

  return (
    <Link
      to={`/tenants/${tenant.id}`}
      className="block rounded-lg border border-gray-200 bg-white p-4 shadow-sm transition-shadow hover:shadow-md active:scale-[0.98]"
    >
      <div className="flex items-start justify-between">
        <div>
          <h3 className="text-base font-semibold text-gray-900">
            Block {tenant.block} · Unit {tenant.unit_number}
          </h3>
          <p className="mt-1 text-base font-semibold text-gray-900">{monthlyFee} / month</p>
        </div>
        <StatusBadge variant={tenant.occupancy} />
      </div>
      <p className="mt-2 text-sm text-gray-500">
        {hasFees
          ? `${mandatoryFeeCount} mandatory fees · ${voluntaryFeeCount} contributions`
          : 'No fees configured'}
      </p>
    </Link>
  );
}
```

---

### `components/tenants/TenantForm.tsx` (component, CRUD)

**Analog:** `LoginForm.tsx` — form component with validation, loading, error states

**Form field pattern** (lines 60-83):
```typescript
// Source: LoginForm.tsx lines 60-83 — field pattern with label, input, error
<div className="space-y-2">
  <label htmlFor="email" className="block text-base font-medium text-gray-700">
    Email
  </label>
  <input
    id="email"
    type="email"
    value={email}
    onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
    className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
      fieldErrors.email
        ? 'border-red-600 focus:border-red-600'
        : 'border-gray-200 focus:border-blue-600'
    }`}
    placeholder="you@example.com"
    autoComplete="email"
    aria-describedby={fieldErrors.email ? 'email-error' : undefined}
    aria-invalid={!!fieldErrors.email}
  />
  {fieldErrors.email && (
    <p id="email-error" className="text-sm text-red-600" role="alert">
      {fieldErrors.email}
    </p>
  )}
</div>
```

**Submit button pattern** (lines 118-124):
```typescript
// Source: LoginForm.tsx lines 118-124 — loading state submit button
<button
  type="submit"
  disabled={loading}
  className="min-h-[44px] w-full rounded-md bg-blue-600 px-4 py-2 text-base font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
>
  {loading ? 'Signing in...' : 'Sign In'}
</button>
```

**Error alert pattern** (lines 54-58):
```typescript
// Source: LoginForm.tsx lines 54-58 — error alert
{error && (
  <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
    {error}
  </div>
)}
```

---

### `components/fees/FeeList.tsx` (component, CRUD)

**No close analog — new component type.**

**Expected pattern:**
```typescript
import FeeCard from './FeeCard'; // or inline
import StatusBadge from '../ui/StatusBadge';
import EmptyState from '../ui/EmptyState';
import type { Fee } from '../../types/fee';

interface FeeListProps {
  mandatoryFees: Fee[];
  voluntaryFees: Fee[];
  onEdit: (fee: Fee) => void;
  onDelete: (fee: Fee) => void;
}

export default function FeeList({ mandatoryFees, voluntaryFees, onEdit, onDelete }: FeeListProps) {
  return (
    <div className="space-y-8">
      <section>
        <h2 className="mb-4 text-base font-semibold text-gray-900">Mandatory Fees</h2>
        {mandatoryFees.length === 0 ? (
          <EmptyState heading="No Mandatory Fees Set" body="Every tenant needs at least one mandatory fee. Add one now." />
        ) : (
          <div className="space-y-3">
            {mandatoryFees.map((fee) => (
              <FeeListItem key={fee.id} fee={fee} onEdit={onEdit} onDelete={onDelete} />
            ))}
          </div>
        )}
      </section>

      <section>
        <h2 className="mb-4 text-base font-semibold text-gray-900">Voluntary Contributions</h2>
        {voluntaryFees.length === 0 ? (
          <EmptyState heading="No Voluntary Contributions Yet" body="Residents can contribute voluntarily here." />
        ) : (
          <div className="space-y-3">
            {voluntaryFees.map((fee) => (
              <FeeListItem key={fee.id} fee={fee} onEdit={onEdit} onDelete={onDelete} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
```

---

### `components/fees/FeeForm.tsx` (component, CRUD)

**Analog:** `LoginForm.tsx` — form component pattern

Same patterns as TenantForm (Input/Select/DatePicker components for fields, submit button pattern, error alert pattern, loading state pattern).

---

### `hooks/useTenants.ts` (hook, CRUD)

**No close analog — TanStack Query is new to the project.** Pattern from RESEARCH.md.

**Expected pattern:**
```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { listTenants, getTenant, createTenant, updateTenant, deleteTenant } from '../services/tenants';
import type { Tenant, CreateTenantRequest } from '../services/tenants';

export function useTenants() {
  return useQuery<Tenant[]>({
    queryKey: ['tenants'],
    queryFn: listTenants,
  });
}

export function useTenant(id: string) {
  return useQuery<Tenant>({
    queryKey: ['tenants', id],
    queryFn: () => getTenant(id),
    enabled: !!id,
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

export function useUpdateTenant(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Partial<CreateTenantRequest>) => updateTenant(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });
}

export function useDeleteTenant() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteTenant(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });
}
```

---

### `hooks/useFees.ts` (hook, CRUD)

**Expected pattern (same pattern as useTenants):**
```typescript
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { listFees, createFee, updateFee, deleteFee } from '../services/fees';
import type { Fee, CreateFeeRequest } from '../services/fees';

export function useFees(tenantId: string) {
  return useQuery({
    queryKey: ['fees', tenantId],
    queryFn: () => listFees(tenantId),
    enabled: !!tenantId,
  });
}

export function useCreateFee(tenantId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateFeeRequest) => createFee(tenantId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fees', tenantId] });
    },
  });
}

export function useUpdateFee(tenantId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ feeId, data }: { feeId: string; data: Partial<CreateFeeRequest> }) =>
      updateFee(tenantId, feeId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fees', tenantId] });
    },
  });
}

export function useDeleteFee(tenantId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (feeId: string) => deleteFee(tenantId, feeId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fees', tenantId] });
    },
  });
}
```

---

### `pages/TenantListPage.tsx` (page, CRUD)

**Analog:** `pages/LoginPage.tsx` — page shell that renders components

**Expected pattern:**
```typescript
// Source: LoginPage.tsx — page wrapper pattern
export default function TenantListPage() {
  const { data: tenants, isLoading, isError } = useTenants();
  const [search, setSearch] = useState('');

  const filtered = useMemo(() => {
    if (!tenants) return [];
    if (!search) return tenants;
    const q = search.toLowerCase();
    return tenants.filter(
      (t) => t.block.toLowerCase().includes(q) || t.unit_number.toLowerCase().includes(q)
    );
  }, [tenants, search]);

  return (
    <div>
      <PageHeader
        title="Tenants"
        action={
          <Link to="/tenants/new" className="min-h-[44px] inline-flex items-center rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700">
            + Add Tenant
          </Link>
        }
      />
      <input
        type="search"
        placeholder="Search by block or unit number..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="mb-4 min-h-[44px] w-full rounded-md border border-gray-200 px-3 py-2 text-base outline-none focus:border-blue-600"
      />
      {isLoading ? (
        <LoadingSkeleton variant="card" count={3} />
      ) : isError ? (
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          Failed to load tenants. Pull down to refresh or try again.
        </div>
      ) : filtered.length === 0 ? (
        <EmptyState heading="No Tenants Yet" body="Start by adding your first tenant to this RT." action={...} />
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {filtered.map((tenant) => (
            <TenantCard key={tenant.id} tenant={tenant} mandatoryFeeCount={0} voluntaryFeeCount={0} />
          ))}
        </div>
      )}
    </div>
  );
}
```

---

### `App.tsx` (modified config)

**Analog:** Existing `App.tsx` — exact match

**Current pattern** (lines 1-31):
```typescript
// Source: App.tsx lines 1-31 — browser router + routes + ProtectedRoute
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ResetPasswordPage from './pages/ResetPasswordPage';
import ProtectedRoute from './routes/ProtectedRoute';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/reset" element={<ResetPasswordPage />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <div>...</div>
            </ProtectedRoute>
          }
        />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
```

**Expected update pattern (add QueryClientProvider + tenant routes):**
```typescript
// Source: CONTEXT.md D-07, UI-SPEC routing
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
      staleTime: 30_000,
      retry: 1,
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
          <Route path="/tenants" element={<ProtectedRoute><AppLayout><TenantListPage /></AppLayout></ProtectedRoute>} />
          <Route path="/tenants/new" element={<ProtectedRoute><AppLayout><TenantCreatePage /></AppLayout></ProtectedRoute>} />
          <Route path="/tenants/:id" element={<ProtectedRoute><AppLayout><TenantDetailPage /></AppLayout></ProtectedRoute>} />
          <Route path="/tenants/:id/edit" element={<ProtectedRoute><AppLayout><TenantEditPage /></AppLayout></ProtectedRoute>} />
          <Route path="/dashboard" element={<ProtectedRoute><AppLayout>...</AppLayout></ProtectedRoute>} />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
```

---

### `index.css` (modified config)

**Analog:** Existing `index.css` — exact match

**Current pattern** (lines 1-30):
```css
@import "tailwindcss";

@theme {
  --color-dominant: #FFFFFF;
  --color-surface: #F9FAFB;
  --color-secondary: #F3F4F6;
  --color-border: #E5E7EB;
  --color-accent: #2563EB;
  --color-accent-hover: #1D4ED8;
  --color-destructive: #DC2626;
  --color-text-primary: #111827;
  --color-text-secondary: #6B7280;
  --color-text-muted: #9CA3AF;
  --color-strength-weak: #DC2626;
  --color-strength-medium: #F59E0B;
  --color-strength-strong: #10B981;
}
```

**Expected update (add sidebar tokens):**
```css
@theme {
  /* existing tokens ... */
  --color-sidebar: #F3F4F6;
  --color-sidebar-hover: #E5E7EB;
  --color-sidebar-active: #DBEAFE;
}
```

---

### `routes/ProtectedRoute.tsx` (modified middleware, request-response)

**Analog:** Existing `ProtectedRoute.tsx` — exact match

**Current `/auth/me` fetch pattern** (lines 30-52):
```typescript
// Source: ProtectedRoute.tsx lines 30-52 — useEffect + fetch for auth check
useEffect(() => {
  fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:3000'}/auth/me`, {
    credentials: 'include',
  })
    .then((res) => {
      if (!res.ok) {
        setIsAuthenticated(false);
        setUser(null);
        return;
      }
      return res.json();
    })
    .then((data) => {
      if (data) {
        setIsAuthenticated(true);
        setUser(data.user);
      }
    })
    .catch(() => {
      setIsAuthenticated(false);
      setUser(null);
    });
}, []);
```

**Expected migration pattern (use TanStack Query instead of useEffect):**
```typescript
import { useQuery } from '@tanstack/react-query';
import { request } from '../services/api';

export default function ProtectedRoute({ children, requiredRole }: ProtectedRouteProps) {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['auth', 'me'],
    queryFn: () => request<{ user: User }>('/auth/me'),
    retry: false,
  });

  const user = data?.user ?? null;
  const isAuthenticated = !isLoading && !isError && !!user;

  // Loading state (unchanged)
  if (isLoading) { /* spinner */ }

  // Redirect (unchanged)
  if (isError || !user) { return <Navigate to="/login" ... />; }

  // Role check (unchanged)
  if (requiredRole && user.role !== requiredRole) { /* access denied */ }

  return <AuthContext.Provider value={{ user, isAuthenticated }}>{children}</AuthContext.Provider>;
}
```

---

## Shared Patterns

### Form Field Components (Input, Select, DatePicker)
**Source:** `LoginForm.tsx` lines 60-110
**Apply to:** `Input.tsx`, `Select.tsx`, `DatePicker.tsx`

Key rules from existing LoginForm field markup:
```typescript
// 1. Container: space-y-2 for label-to-field gap
<div className="space-y-2">

// 2. Label: block, text-base, font-medium, text-gray-700
<label htmlFor={id} className="block text-base font-medium text-gray-700">
  {label}
</label>

// 3. Input: min-h-[44px], w-full, rounded-md, border, px-3, py-2, text-base, outline-none, transition-colors
className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
  error ? 'border-red-600 focus:border-red-600' : 'border-gray-200 focus:border-blue-600'
}`}

// 4. Error state: aria-invalid, aria-describedby
aria-invalid={!!error}
aria-describedby={error ? errorId : undefined}

// 5. Error message: p tag, text-sm, text-red-600, role="alert"
{error && (
  <p id={errorId} className="text-sm text-red-600" role="alert">{error}</p>
)}
```

### Submit Button Pattern
**Source:** `LoginForm.tsx` lines 118-124
**Apply to:** `TenantForm.tsx`, `FeeForm.tsx`, all page CTAs

```typescript
<button
  type="submit"
  disabled={loading}
  className="min-h-[44px] w-full rounded-md bg-blue-600 px-4 py-2 text-base font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
>
  {loading ? 'Saving...' : 'Save Tenant'}
</button>
```

### Error Alert Pattern
**Source:** `LoginForm.tsx` lines 54-58
**Apply to:** All forms and pages with error state

```typescript
{error && (
  <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
    {error}
  </div>
)}
```

### Page Shell Pattern
**Source:** `LoginPage.tsx` lines 3-27
**Apply to:** All page files (`TenantListPage.tsx`, `TenantCreatePage.tsx`, etc.)

```typescript
export default function SomePage() {
  return (
    <div>
      {/* Page content here. No full-screen layout wrapper — AppLayout handles that. */}
    </div>
  );
}
```

### Service Test Pattern
**Source:** `services/auth.test.ts` lines 1-145
**Apply to:** `services/tenants.test.ts`, `services/fees.test.ts`, `services/api.test.ts`

Key patterns:
```typescript
// Line 6-7: mock fetch
beforeEach(() => { global.fetch = vi.fn(); });
afterEach(() => { vi.restoreAllMocks(); });

// Line 15-18: mock successful response
(fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
  ok: true,
  json: () => Promise.resolve({ tenants: [...] }),
});

// Line 45-48: mock error response
(fetch as ReturnType<typeof vi.fn>).mockResolvedValue({
  ok: false,
  json: () => Promise.resolve({ error: 'Not found', code: 'NOT_FOUND' }),
});

// Line 50: assert error thrown
await expect(func()).rejects.toThrow('Not found');

// Line 22-29: assert fetch called with correct URL and options
expect(fetch).toHaveBeenCalledWith(
  expect.stringContaining('/api/tenants'),
  expect.objectContaining({ method: 'GET', credentials: 'include' })
);
```

### Component Test Pattern (mock service)
**Source:** `LoginForm.test.tsx` lines 1-117
**Apply to:** `TenantForm.test.tsx`, `FeeForm.test.tsx`, page tests

Key patterns:
```typescript
// Line 7-9: mock service module
vi.mock('../../services/auth', () => ({
  login: vi.fn(),
}));

// Line 11-17: render helper with MemoryRouter
function renderComponent() {
  return render(
    <MemoryRouter>
      <Component />
    </MemoryRouter>
  );
}

// Line 20-22: reset mocks
beforeEach(() => { vi.clearAllMocks(); });

// Line 33: mock resolved value
(auth.login as ReturnType<typeof vi.fn>).mockResolvedValue({ ... });

// Line 89-90: mock rejected value
(auth.login as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('message'));

// Line 42: wait for async operation
await waitFor(() => { expect(...).toBeInTheDocument(); });
```

### Component Test Pattern (mock TanStack Query)
**Source:** RESEARCH.md lines 637-691, CONTEXT.md D-15
**Apply to:** All page and component tests using TanStack Query hooks

```typescript
vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual('@tanstack/react-query');
  return {
    ...actual,
    useQuery: vi.fn(),
    useMutation: vi.fn(),
  };
});

import { useQuery } from '@tanstack/react-query';

// Mock loading state:
(useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
  data: undefined, isLoading: true, isError: false,
});

// Mock data state:
(useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
  data: [...], isLoading: false, isError: false,
});

// Mock error state:
(useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
  data: undefined, isLoading: false, isError: true,
});
```

### IDR Currency Formatting Utility
**Source:** UI-SPEC copywriting contract, RESEARCH.md Pitfall 5
**Apply to:** TenantCard, FeeList, all amount displays

```typescript
export function formatIDR(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
}
```

### Deferred Client-Side Search Pattern
**Source:** CONTEXT.md agent's discretion — debounced filter (300ms) by block and unit number
**Apply to:** TenantListPage

```typescript
import { useState, useMemo } from 'react';

const [search, setSearch] = useState('');

const filteredTenants = useMemo(() => {
  if (!tenants) return [];
  if (!search.trim()) return tenants;
  const query = search.toLowerCase();
  return tenants.filter(
    (t) => t.block.toLowerCase().includes(query) || t.unit_number.toLowerCase().includes(query)
  );
}, [tenants, search]);
```

---

## No Analog Found

Files with no close match in the codebase (planner should use RESEARCH.md patterns and shared patterns above):

| File | Role | Data Flow | Reason |
|---|---|---|---|
| `components/ui/StatusBadge.tsx` | component | — | No existing badge/status component |
| `components/ui/ConfirmDialog.tsx` | component | event-driven | No existing dialog/modal component |
| `components/ui/PageHeader.tsx` | component | — | No existing page header pattern |
| `components/ui/LoadingSkeleton.tsx` | component | — | No existing skeleton component |
| `components/ui/EmptyState.tsx` | component | — | No existing empty state pattern |
| `components/layout/AppLayout.tsx` | component | — | No existing layout/shell component |
| `components/tenants/TenantCard.tsx` | component | — | No existing card component |
| `components/fees/FeeList.tsx` | component | CRUD | No existing list component |
| `hooks/*.ts` | hook | CRUD | TanStack Query is new to project — no existing hook patterns |

These 9 component files and 9 hook files have no exact analog. The RESEARCH.md provides reference patterns (Pattern 2 for hooks, UI-SPEC provides visual specs). All new components follow the established Tailwind-v4-without-component-library convention from Phase 1.

---

## Metadata

**Analog search scope:** `apps/web/src/` (all existing source files)
**Files scanned:** 25 (all files in apps/web/src/)
**Pattern extraction date:** 2026-05-24

**Existing files evaluated as analogs:**
- `services/auth.ts` — service fetch + error pattern
- `types/auth.ts` — interface-only type pattern
- `components/auth/LoginForm.tsx` — form field markup, validation, submit, error alert
- `components/auth/LoginForm.test.tsx` — mock service + render + fireEvent pattern
- `pages/LoginPage.tsx` — page shell pattern
- `pages/LoginPage.test.tsx` — page test pattern
- `routes/ProtectedRoute.tsx` — auth guard with fetch + role check
- `routes/ProtectedRoute.test.tsx` — route guard test with mock fetch
- `App.tsx` — BrowserRouter + Routes pattern
- `index.css` — Tailwind v4 @theme token pattern
- `services/auth.test.ts` — service test with mock fetch
- `test/setup.ts` — minimal test setup
- `vitest.config.js` — test runner config
