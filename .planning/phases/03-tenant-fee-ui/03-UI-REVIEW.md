---
phase: 3
slug: tenant-fee-ui
score: 18/24
top_fixes:
  - "Add focus trap to ConfirmDialog and FeeForm modal — accessibility blocker for keyboard users"
  - "Fix sidebar: replace ChevronLeft with Settings icon for Settings item, add Reports placeholder with BarChart3"
  - "Add prefers-reduced-motion support to animations (skeleton pulse, sidebar slide, card hover scale)"
date: 2026-05-24
---

# UI Audit: Phase 3 — Tenant & Fee UI

**Audited:** 2026-05-24
**Baseline:** 03-UI-SPEC.md (approved design contract)
**Screenshots:** not captured (no dev server running — code-only audit)
**Registry Safety:** not applicable (no shadcn/registry components used)

---

## Summary

| Pillar | Score |
|--------|-------|
| 1. Copywriting | 4/4 |
| 2. Visuals | 3/4 |
| 3. Color | 3/4 |
| 4. Typography | 3/4 |
| 5. Spacing | 3/4 |
| 6. Experience Design | 2/4 |
| **Total** | **18/24** |

---

## Top 3 Priority Fixes

1. **Add focus trap to ConfirmDialog and FeeForm modal** — Both dialogs lack focus trapping. Tab can move focus behind the modal overlay, which breaks keyboard navigation for screen reader users and violates WCAG 2.1.2 (No Keyboard Trap) guidance. *(apps/web/src/components/ui/ConfirmDialog.tsx, apps/web/src/components/fees/FeeForm.tsx)*
2. **Fix sidebar navigation icons and items** — Settings uses `ChevronLeft` instead of `Settings` icon from lucide-react. Reports nav item (`BarChart3` icon, per UI-SPEC sidebar table) is entirely missing. These are visual contract deviations that will compound when Phase 4 needs to add Reports. *(apps/web/src/components/layout/AppLayout.tsx)*
3. **Add `prefers-reduced-motion` support** — The spec explicitly requires respecting reduced motion: disable sidebar slide animation, disable card hover scale, disable skeleton pulse. No `@media (prefers-reduced-motion)` queries found anywhere in the components. This is an accessibility requirement gap. *(all animation-enabled components)*

---

## Detailed Findings

### Pillar 1: Copywriting (4/4)

**Good:**
- All CTAs match the copywriting contract exactly: "+ Add Tenant", "Save Tenant", "Update Tenant", "Record Fee", "Save Fee", "Cancel"
- Back links use "← Back to Tenants" consistently across all pages
- Empty state headings and body copy match spec: "No Tenants Yet" / "Start by adding your first tenant to this RT."
- Form field labels match spec: "Block", "Unit Number", "Occupancy Status", "Monthly Fee (Rp)", "Description", "Amount (Rp)", "Effective Date", "Fee Type", "Payment Date"
- Submit button shows "Saving..." with spinner during API calls
- Delete confirmations use spec-correct text: "Delete this fee? This cannot be undone."
- Form validation error messages match spec verbatim (e.g., "Block is required.", "Monthly fee must be a positive amount.")
- Fee section empty states use spec text: "No Mandatory Fees Set" / "No Voluntary Contributions Yet"
- Search placeholder matches spec: "Search by block or unit number..."

**Issues:**
- No deviations found. Copywriting contract is fully satisfied.

---

### Pillar 2: Visuals (3/4)

**Good:**
- Tenant cards have clear visual hierarchy: block+unit title (semibold), monthly fee (semibold), fee summary (secondary color)
- Status badges use correct semantic colors (blue for occupied, amber for vacant, emerald for paid, red for unpaid)
- Card interactive states: `hover:shadow-md`, `active:scale-[0.98]` as specified
- Action icons (Pencil, Trash2) with hover color transitions (gray → blue/red)
- Loading skeletons with pulse animation for card/list/form variants
- Search bar with Search icon indicator
- Empty state with centered layout + CTA button
- Skip link as first focusable element (visually hidden until focused)
- Tenant detail info summary bar with compact layout

**Issues:**
- **WARNING**: Settings nav item uses `ChevronLeft` icon (line 4, AppLayout.tsx) instead of `Settings` icon from lucide-react as specified in UI-SPEC sidebar table
- **WARNING**: Reports nav item (`BarChart3` icon, `/reports` route) is missing from sidebar entirely — the implementation added a Dashboard item instead which is not in the spec's sidebar table
- **MINOR**: Fee cards (FeeCard internal component in FeeList.tsx) lack hover state — all other interactive cards in the system implement `hover:shadow-md` but FeeCard does not
- **MINOR**: EmptyState decorative element is a generic gray circle (`rounded-full bg-gray-200`) rather than a meaningful visual illustration or icon matching the context

---

### Pillar 3: Color (3/4)

**Good:**
- All semantic badge colors match spec exactly (blue-100/blue-700 for occupied, amber-100/amber-800 for vacant, emerald-100/emerald-800 for paid, red-100/red-800 for unpaid)
- Primary action buttons use `bg-blue-600` / `hover:bg-blue-700` (accent) as specified
- Destructive buttons use `bg-red-600` / `hover:bg-red-700` as specified
- Form input borders use `border-gray-200` default / `border-blue-600` focus / `border-red-600` error as specified
- Text colors follow hierarchy: `text-gray-900` (primary), `text-gray-600/700` (secondary), `text-gray-500` (muted)
- Page backgrounds use `bg-gray-50` (surface) as specified
- Sidebar uses `bg-gray-100` as specified
- Error alerts use `bg-red-50` / `text-red-700` as specified

**Issues:**
- **MINOR**: CSS custom properties defined in `index.css` (`--color-accent`, `--color-destructive`, `--color-sidebar`, etc.) are never referenced by any component — all styling uses hardcoded Tailwind classes like `bg-blue-600`, `text-gray-900`, `bg-gray-100`. This means the design tokens in `@theme` block provide no actual abstraction or theming benefit.
- **MINOR**: Sidebar active item uses `border-l-4 border-blue-600 bg-blue-100` directly rather than the `--color-sidebar-active` custom property defined for this purpose
- 60/30/10 distribution is maintained naturally (white/gray-50 surfaces dominate, gray-100 sidebar is secondary, blue-600 accent is sparse and reserved for CTAs and active states)

---

### Pillar 4: Typography (3/4)

**Good:**
- Body text uses `text-base` (16px) — matches spec
- Small text uses `text-sm` (14px) for labels, helper text, secondary info — matches spec
- Page headings use `text-2xl` (24px) — matches spec's Heading role
- Sidebar brand uses `text-[28px]` — matches spec's Display role
- Font weight system follows spec: 400 (regular) for body text, 600 (semibold) for headings and labels
- System font stack is correctly implemented in `index.css`

**Issues:**
- **MINOR**: StatusBadge uses `text-xs` (12px) instead of spec's Small (14px). Spec says "Status badges" belong to the "Small" typography role (14px).
- **MINOR**: EmptyState heading uses `text-lg` (18px) — not in the spec's 4-size system (14px, 16px, 24px, 28px)
- **MINOR**: Dialog titles in ConfirmDialog use `text-lg` (18px) — likewise outside the spec's size palette
- **MINOR**: Fee timestamps and "Fee N" labels use `text-xs` (12px) — not in the 4-size system
- Usage tally: `text-xs` (12px), `text-sm` (14px), `text-base` (16px), `text-lg` (18px), `text-2xl` (24px) — 5 sizes where spec defines 4

---

### Pillar 5: Spacing (3/4)

**Good:**
- Card padding uses `p-4` (16px) — matches spec's md token
- Section gaps use `space-y-6` (24px) and `space-y-8` (32px+) — matches spec's lg/xl tokens
- Form field stacks use `space-y-2` (8px) — matches spec's sm
- Touch targets: all buttons, sidebar links, and interactive elements use `min-h-[44px]` — WCAG 2.5.5 compliant
- Icon-only buttons use `min-w-[44px]` alongside min-height — compliant
- Sidebar links use `py-2.5 px-4` — matches spec's sidebar link sizing formula (10px padding + 24px line-height = 44px)
- Confirm dialog buttons use `min-w-[96px]` — matches spec
- Search bar uses `pl-10` for icon spacing — reasonable
- Responsive grid: 1 column mobile, 2 columns sm, 3 columns lg — matches spec

**Issues:**
- **MINOR**: Tenant form card uses `max-w-lg` (512px) — spec says "max-width 600px" for tenant create/edit forms
- **MINOR**: PageHeader component has no `mb-6` bottom margin — pages add their own `space-y-6` on the parent container, but the PageHeader component itself doesn't enforce consistent bottom spacing, leaving it to each consumer
- **MINOR**: Sidebar width is `w-64` (256px) — spec says 260px (minor, 4px off)

---

### Pillar 6: Experience Design (2/4)

**Good:**
- **Loading states**: LoadingSkeleton variants (card, list, form) used on tenant list, fee list, and edit page loader — matches spec
- **Error states**: Red `role="alert"` banners on list load failure, form submit failure, and edit page error — matches spec
- **Empty states**: EmptyState component used for all empty scenarios (no tenants, no search results, no mandatory fees, no voluntary contributions) — matches spec
- **Form validation**: Client-side validation blocks submit with inline error messages beneath fields — matches spec
- **Search**: 300ms debounced client-side filter by block and unit number — matches spec
- **Sorting**: Tenants sorted by block then unit_number ascending — matches spec
- **Delete confirmation**: ConfirmDialog with Escape key handler and backdrop click dismiss — matches spec
- **Keyboard accessibility**: TenantCard has `tabIndex={0}` and Enter key handler, form fields have proper labels and `aria-invalid`/`aria-describedby` — matches spec
- **Skip link**: First focusable element on page in AppLayout — matches spec
- **Sidebar navigation**: Active route highlighted with blue accent, mobile hamburger toggle with overlay backdrop — matches spec
- **Back navigation**: "← Back to Tenants" link present on all sub-pages — matches spec
- **Submit button states**: Disabled + spinner + "Saving..." during submission — matches spec

**Issues:**
- **BLOCKER**: **No focus trap in dialogs** — ConfirmDialog and the FeeForm modal both allow Tab focus to escape behind the modal overlay. This breaks keyboard-only and screen reader navigation. The spec explicitly requires "Focus trap inside dialog" under Accessibility Requirements. *(ConfirmDialog.tsx:44-47, TenantDetailPage.tsx:142-162)*
- **WARNING**: **No `prefers-reduced-motion` support** — All animations (sidebar slide/translate, card hover scale, skeleton pulse) will play even for users who have requested reduced motion. The spec explicitly requires "Respect prefers-reduced-motion — disable sidebar slide animation, disable card hover scale, disable skeleton pulse." No `@media (prefers-reduced-motion)` queries found anywhere. *(AppLayout.tsx:65, TenantCard.tsx:29, LoadingSkeleton.tsx:16, all animation classes)*
- **WARNING**: **No focus return after dialog close** — The spec requires "After dialog close, focus returns to trigger button." ConfirmDialog focuses the Cancel button when opened but doesn't restore focus to the trigger when closed. The dialog component unmounts without focus management. *(ConfirmDialog.tsx:24-31)*
- **MINOR**: **No clear search button** — Spec mentions "Clear search button when input has value" for the tenant list search. The search field has no clear/X button when text is entered. *(TenantListPage.tsx:57-67)*
- **MINOR**: **No optimistic delete for fees** — Spec mentions "On successful delete: remove item from list (optimistic)" for the delete confirmation pattern. The current implementation waits for the server response before updating UI. *(TenantDetailPage.tsx:53-61)*
- **MINOR**: **Missing `aria-labelledby` on ConfirmDialog** — The dialog element has `role="dialog"` and `aria-modal="true"` but no `aria-labelledby` pointing to the title `<h2>`, which would improve screen reader announcement. Not a spec requirement but an a11y best practice gap. *(ConfirmDialog.tsx:44-47)*
- **MINOR**: **formatIDR utility duplicated** — The IDR currency formatting function is defined in both `TenantCard.tsx` (line 4) and `FeeList.tsx` (line 7), exported from both. The `TenantDetailPage` imports from FeeList, creating a hidden dependency. Should be extracted to a shared utility.

---

## Registry Safety

Not applicable. No shadcn initialized. No third-party registry components used. All 10 new UI components are hand-rolled with Tailwind utility classes, consistent with Phase 1's established pattern. UI-SPEC.md §Registry Safety confirms: "(none) — not applicable, no registry components used."

---

## Files Audited

### Components (18 files)
- `apps/web/src/components/layout/AppLayout.tsx` — Responsive sidebar + header + content shell
- `apps/web/src/components/ui/Input.tsx` — Reusable input with label/error/aria
- `apps/web/src/components/ui/Select.tsx` — Reusable select with label/error/options
- `apps/web/src/components/ui/DatePicker.tsx` — Native date input wrapper
- `apps/web/src/components/ui/StatusBadge.tsx` — Occupancy/payment status indicator
- `apps/web/src/components/ui/ConfirmDialog.tsx` — Delete confirmation modal dialog
- `apps/web/src/components/ui/PageHeader.tsx` — Page title + action button layout
- `apps/web/src/components/ui/LoadingSkeleton.tsx` — Pulse-animated skeleton placeholders
- `apps/web/src/components/ui/EmptyState.tsx` — Empty state with heading/body/CTA
- `apps/web/src/components/tenants/TenantCard.tsx` — Tenant summary card
- `apps/web/src/components/tenants/TenantForm.tsx` — Tenant create/edit form with dynamic fees
- `apps/web/src/components/fees/FeeList.tsx` — Fee list with mandatory/voluntary sections
- `apps/web/src/components/fees/FeeForm.tsx` — Fee create/edit form (modal)

### Pages (4 files)
- `apps/web/src/pages/TenantListPage.tsx` — Tenant list with search, grid, states
- `apps/web/src/pages/TenantCreatePage.tsx` — Create tenant wrapper
- `apps/web/src/pages/TenantEditPage.tsx` — Edit tenant with delete
- `apps/web/src/pages/TenantDetailPage.tsx` — Tenant detail + fee management

### Services & Hooks (8 files)
- `apps/web/src/services/api.ts` — Shared request() helper
- `apps/web/src/services/tenants.ts` — Tenant CRUD service
- `apps/web/src/services/fees.ts` — Fee CRUD service
- `apps/web/src/services/auth.ts` — Refactored auth service
- `apps/web/src/hooks/useTenants.ts` — Tenant query hooks
- `apps/web/src/hooks/useFees.ts` — Fee query hooks
- `apps/web/src/hooks/useCreateTenant.ts` — Create tenant mutation hook
- `apps/web/src/hooks/useDeleteTenant.ts` — Delete tenant mutation hook

### Config & Routing (4 files)
- `apps/web/src/App.tsx` — Routes + QueryClientProvider
- `apps/web/src/index.css` — @theme tokens + base styles
- `apps/web/src/routes/ProtectedRoute.tsx` — Auth guard with useQuery
- `apps/web/src/types/tenant.ts` / `apps/web/src/types/fee.ts` — TypeScript interfaces

---

## Scoring Rationale

| Pillar | Score | Rationale |
|--------|-------|-----------|
| Copywriting | 4/4 | All CTAs, labels, placeholders, empty states, error messages match spec verbatim. No deviations found. |
| Visuals | 3/4 | Strong component visual design and hierarchy. Deducted for wrong Settings icon (ChevronLeft), missing Reports nav item, and FeeCard lacking hover state. |
| Color | 3/4 | All Tailwind color classes match spec values exactly. Deducted because CSS custom properties in `@theme` are never consumed by components (hardcoded classes throughout), making the token system decorative only. |
| Typography | 3/4 | Body (16px), labels (14px), headings (24px) all correct. Deducted for using `text-xs` (12px) on badges instead of spec's 14px, and `text-lg` (18px) on dialog/empty state headings — these are outside the declared 4-size system. |
| Spacing | 3/4 | 44px touch targets enforced throughout. Spacing scale is consistent. Deducted for form width (512px vs spec's 600px) and missing `mb-6` on PageHeader. |
| Experience Design | 2/4 | State coverage (loading/error/empty) is solid. Major deductions for missing focus trap in dialogs (spec requirement), missing `prefers-reduced-motion` (spec requirement), and missing search clear button. These are accessibility gaps that block users. |

---

## Recommendation Count

- Priority fixes: 3 (listed above)
- Minor recommendations: 8 (documented in Detailed Findings)
