---
status: complete
phase: 03-tenant-fee-ui
source:
  - 03-01-SUMMARY.md
  - 03-02-SUMMARY.md
  - 03-03-SUMMARY.md
  - 03-04-SUMMARY.md
  - 03-05-SUMMARY.md
started: 2026-05-24T17:00:00Z
updated: 2026-05-24T20:07:00Z
---

## Current Test

[testing complete]

## Tests

### 1. Sidebar Navigation — Desktop Layout
expected: On desktop (≥1024px), sidebar is fixed 260px left panel with nav links (Dashboard, Tenants, Settings for officers). Active route highlighted. No hamburger on desktop.
result: issue
reported: "Nav shows 'Reports' instead of 'Dashboard'"
severity: major
fix: Changed "Reports" → "Dashboard" in AppLayout.tsx + updated tests (7/7 pass)

### 2. Sidebar Navigation — Mobile Layout
expected: On mobile (<1024px), sidebar is hidden by default. Hamburger button visible. Tapping hamburger opens sidebar as overlay drawer. Tapping backdrop or nav link closes it.
result: pass

### 3. Sidebar — Role-Based Visibility
expected: Settings link is only visible when logged in as rt_officer or rw_officer. Dashboard and Tenants are visible to all roles.
result: pass

### 4. Tenant List — Loading State
expected: While loading, the page shows a 3-column skeleton grid (LoadingSkeleton with card variant). No flicker.
result: pass

### 5. Tenant List — Empty State
expected: When no tenants exist, the page shows an EmptyState with heading "No Tenants Yet" and body text. No redundant CTA (PageHeader already has "+ Add Tenant").
result: pass

### 6. Tenant List — Search/Filter
expected: A search input filters tenant cards by block or unit number. Filtering is client-side with ~300ms debounce. Results update in real-time as user types.
result: pass

### 7. Tenant List — Responsive Grid
expected: Tenant cards display in a responsive grid: 1 column on mobile, 2 on tablet, 3 on desktop. Each card shows occupancy badge, formatted IDR monthly fee, and fee summary.
result: pass

### 8. Create Tenant — Form Fields
expected: Create Tenant form has: Block (required), Unit Number (required), Occupancy Status (select: Occupied/Vacant), Monthly Fee (required, positive number), and a "Mandatory Fees" section with add/remove rows (description, amount, due date per row). At least 1 mandatory fee required.
result: issue
reported: "Fee amount changes when interacting with date field — 50000 becomes 49996, 100000 becomes 99994"
severity: major
fix: Removed `type="number"` from fee amount inputs — browser locale formatting was corrupting the value on re-render. Replaced with `inputMode="numeric"`.

### 9. Create Tenant — Validation
expected: Submitting with invalid data shows inline error messages per field. Block and unit number required. Monthly fee must be positive. Each mandatory fee row requires description, amount, and date. Total mandatory fees cannot exceed monthly fee cap.
result: issue
reported: "With 2 mandatory fees summing > monthly cap, form still submits"
severity: major
fix: Changed validation from per-fee cap check to sum-based check. Added test for multi-fee cap scenario.

### 10. Create Tenant — Successful Submit
expected: Filling all fields and submitting shows submit button spinner + "Saving..." text. On success, navigates to /tenants list page. New tenant appears in the list.
result: issue
reported: '"Saving..." text not shown during submit'
severity: minor
fix: Added local `isSubmitting` state in TenantForm that activates on submit click, ensuring loading text shows even on fast API responses.

### 11. Edit Tenant — Pre-Populated Form
expected: Navigating to /tenants/:id/edit shows the form pre-populated with existing tenant data including Block, Unit, Occupancy, Monthly Fee. Pencil icon on TenantCard provides quick access to edit page.
result: pass

### 12. Edit Tenant — Delete with Confirmation
expected: Edit page has a Delete button. Clicking it opens ConfirmDialog with warning text. Confirming deletes the tenant and navigates to /tenants. Canceling closes the dialog.
result: pass

### 13. Tenant Detail — Header and Fee Sections
expected: Tenant detail page shows tenant info header (block/unit, occupancy badge, monthly fee). Below that, FeeList with two sections: "Mandatory Fees" and "Voluntary Fees". Each section shows fee cards or empty state.
result: pass

### 14. Record Fee — Modal Form
expected: Clicking "Record Fee" opens a modal overlay with FeeForm: type selector (Mandatory/Voluntary), description, amount, effective date, payment date. Validation: amount > 0 and ≤ monthly fee cap, sum of all mandatory fees ≤ monthly fee. Backend should also enforce this.
result: issue
reported: "Validation must check total mandatory fees sum against monthly fee cap (backend too)"
severity: major
fix: Added `existingMandatoryTotal` and `editingFeeId` props to FeeForm. Validation now checks `existingTotal + newAmount <= monthlyFee` for mandatory fees. Backend enforcement is Phase 2 — needs separate work.

### 15. Record Fee — Successful Creation
expected: Filling FeeForm and submitting creates the fee. Fee appears immediately in the list under the correct section. Modal closes on success.
result: pass

### 16. Edit Fee — Pre-Populated Form
expected: Clicking edit icon on a fee card opens FeeForm pre-populated with that fee's data. Changes save on submit.
result: pass

### 17. Delete Fee — Confirmation
expected: Clicking delete icon on a fee card opens ConfirmDialog. Confirming deletes the fee and it disappears from the list. Canceling closes the dialog.
result: pass

### 18. Fee Card — Status Badge
expected: Each fee card shows a StatusBadge: "Paid" (green) or "Unpaid" (yellow/red) based on payment status. Amount shown in IDR format (Rp). Effective date shown in "DD MMM YYYY" format.
result: pass

## Summary

total: 18
passed: 13
issues: 5
pending: 0
skipped: 0
blocked: 0

## Gaps

- truth: "Desktop sidebar shows 'Dashboard' as nav link (not 'Reports')"
  status: failed
  reason: "User reported: Nav shows 'Reports' instead of 'Dashboard'"
  severity: major
  test: 1
  root_cause: "AppLayout.tsx navItems had label 'Reports' with path '/reports' — no /reports route exists. Should be 'Dashboard' with path '/dashboard'."
  artifacts:
    - path: "apps/web/src/components/layout/AppLayout.tsx"
      issue: "Nav item label 'Reports' should be 'Dashboard'"
  missing:
    - "Changed label from 'Reports' to 'Dashboard'"
    - "Changed path from '/reports' to '/dashboard'"
    - "Changed icon from BarChart3 to LayoutDashboard"
    - "Updated AppLayout.test.tsx to reference 'Dashboard' instead of 'Reports'"
    - "Added test for Dashboard active highlighting"

- truth: "Fee amount input maintains entered value on re-render"
  status: failed
  reason: "User reported: Fee amount changes from 50000 to 49996 after interacting with date picker"
  severity: major
  test: 8
  root_cause: "`type=\"number\"` with React controlled input causes browser locale formatting issues — some browsers reformat the display value (e.g., with thousands separators) and re-parse incorrectly on re-render."
  artifacts:
    - path: "apps/web/src/components/tenants/TenantForm.tsx"
      issue: "Fee amount input uses `type=\"number\"` on line 276"
  missing:
    - "Changed to `inputMode=\"numeric\"` (no `type=\"number\"`)"
    - "Browser treats as text input, no locale formatting applied"

- truth: "Total mandatory fees cannot exceed monthly fee cap"
  status: failed
  reason: "User reported: 2 fees summing > monthly cap still submits"
  severity: major
  test: 9
  root_cause: "Validation checked each fee amount individually against cap, not the sum of all mandatory fees."
  artifacts:
    - path: "apps/web/src/components/tenants/TenantForm.tsx"
      issue: "Per-fee cap check `f.amount > fee` on line 85"
  missing:
    - "Replaced with sum-based check: `totalMandatoryFees > fee`"
    - "Added test for multi-fee cap scenario"

- truth: "Submit button shows Saving... text during submission"
  status: failed
  reason: "User reported: Saving... text not shown"
  severity: minor
  test: 10
  root_cause: "TanStack Query's isPending flips too fast on local network — React never renders the loading state."
  artifacts:
    - path: "apps/web/src/components/tenants/TenantForm.tsx"
      issue: "Button relied solely on `isLoading` prop"
  missing:
    - "Added local `isSubmitting` state set synchronously before async submit"
    - "Button uses `isLoading || isSubmitting`"

- truth: "FeeForm validates total mandatory fees sum against monthly fee cap"
  status: failed
  reason: "User reported: Validation must check total mandatory fees sum"
  severity: major
  test: 14
  root_cause: "FeeForm only validated single fee amount against cap, not cumulative sum."
  artifacts:
    - path: "apps/web/src/components/fees/FeeForm.tsx"
      issue: "No cumulative sum check for mandatory fees"
    - path: "apps/web/src/pages/TenantDetailPage.tsx"
      issue: "existingMandatoryTotal not passed to FeeForm"
  missing:
    - "Added `existingMandatoryTotal` + `editingFeeId` props to FeeForm"
    - "Validation: `otherTotal + newAmount > monthlyFee` for mandatory fees"
    - "Backend enforcement noted as Phase 2 work"


