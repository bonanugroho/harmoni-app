---
wave: 4
depends_on:
  - 002
files_modified:
  - apps/web/src/pages/LoginPage.tsx
  - apps/web/src/pages/RegisterPage.tsx
  - apps/web/src/pages/ResetPasswordPage.tsx
  - apps/web/src/components/auth/LoginForm.tsx
  - apps/web/src/components/auth/RegisterForm.tsx
  - apps/web/src/components/auth/ResetPasswordForm.tsx
  - apps/web/src/services/auth.ts
  - apps/web/src/routes/ProtectedRoute.tsx
autonomous: true
requirements:
  - AUTH-01
---

# Plan 4: Frontend Auth Pages

## Objective
Create mobile-first authentication pages (login, register, reset password) with React + Tailwind CSS.

## Tasks

### Task 1: Auth Service (API Client)
<read_first>
- apps/web/src/services/auth.ts (create)
- apps/api/internal/delivery/http/auth_handler.go (API contract reference)
</read_first>

<action>
Implement auth API client:
- login(email, password) → Promise<User>
- register(email, password, fullName, territoryId) → Promise<User>
- requestPasswordReset(email) → Promise<void>
- confirmPasswordReset(token, newPassword) → Promise<void>
- logout() → clear local state
- Use fetch with credentials: 'include' for httpOnly cookies
- Handle error responses: extract error message from {"error": "...", "code": "..."}
</action>

<acceptance_criteria>
- apps/web/src/services/auth.ts implements all auth API methods
- login sends POST /auth/login with credentials: 'include'
- register sends POST /auth/register with user data
- requestPasswordReset sends POST /auth/reset
- confirmPasswordReset sends POST /auth/reset/confirm
- Error responses are parsed and thrown as Error objects
- `npm test -- auth.test.ts` passes with mocked fetch
</acceptance_criteria>

---

### Task 2: Login Form Component
<read_first>
- apps/web/src/components/auth/LoginForm.tsx (create)
- apps/web/src/services/auth.ts
- .planning/PROJECT.md (mobile-first constraint)
</read_first>

<action>
Create login form component:
- Fields: email, password
- Validation: email format, password required
- Submit: call login() service, redirect to dashboard on success
- Error display: show error message from API
- Loading state: disable submit button during request
- Mobile-first: full-width inputs, large touch targets (min 44px)
- Link to register page and reset password page
</action>

<acceptance_criteria>
- LoginForm renders with email and password fields
- Submit button disabled during loading state
- Invalid email format shows inline error
- Empty password shows "Password is required" error
- Successful login redirects to /dashboard
- Failed login shows error message from API
- Component is responsive on 320px width screens
- Touch targets are minimum 44px height
- `npm test -- LoginForm.test.tsx` passes
</acceptance_criteria>

---

### Task 3: Register Form Component
<read_first>
- apps/web/src/components/auth/RegisterForm.tsx (create)
- apps/web/src/services/auth.ts
- .planning/phases/01-core-authentication-rbac/01-CONTEXT.md (password policy)
</read_first>

<action>
Create registration form component:
- Fields: full_name, email, password, territory_id (dropdown)
- Validation: email format, password complexity (8+ chars, uppercase, lowercase, number, symbol), territory required
- Password strength indicator: weak/medium/strong
- Submit: call register() service, redirect to login on success
- Error display: show error message from API
- Mobile-first: stacked layout, full-width inputs
- Link to login page
</action>

<acceptance_criteria>
- RegisterForm renders with all required fields
- Territory dropdown populated with available territories
- Password validation enforces: 8+ chars, uppercase, lowercase, number, symbol
- Password strength indicator updates in real-time
- Weak password ("password") shows strength "weak" in red
- Strong password ("SecurePass123!") shows strength "strong" in green
- Successful registration redirects to /login
- Duplicate email shows "Email already registered" error
- Component is responsive on 320px width screens
- `npm test -- RegisterForm.test.tsx` passes
</acceptance_criteria>

---

### Task 4: Reset Password Form Component
<read_first>
- apps/web/src/components/auth/ResetPasswordForm.tsx (create)
- apps/web/src/services/auth.ts
</read_first>

<action>
Create reset password flow:
- Step 1: Email input → request reset link
- Step 2: Token + new password input (from email link)
- Validation: email format, password complexity, token required
- Success message: "Password reset link sent to your email"
- Error display: show error message from API
- Mobile-first: simple stacked layout
- Link back to login page
</action>

<acceptance_criteria>
- ResetPasswordForm shows email input initially
- Submitting email shows "Reset link sent" message
- Token + password form shown when URL contains ?token= parameter
- Password validation same as registration form
- Successful reset redirects to /login
- Invalid/expired token shows "Invalid or expired reset token" error
- Component is responsive on 320px width screens
- `npm test -- ResetPasswordForm.test.tsx` passes
</acceptance_criteria>

---

### Task 5: Protected Route Component
<read_first>
- apps/web/src/routes/ProtectedRoute.tsx (create)
- apps/web/src/services/auth.ts
</read_first>

<action>
Create protected route wrapper:
- Check authentication status (cookie-based)
- Redirect to /login if not authenticated
- Pass user data to child components
- Optional role-based access control (redirect if wrong role)
- Loading state while checking auth status
</action>

<acceptance_criteria>
- ProtectedRoute redirects to /login when not authenticated
- ProtectedRoute renders children when authenticated
- ProtectedRoute passes user data to children via context
- Optional role prop restricts access by role
- Loading spinner shown while checking auth status
- `npm test -- ProtectedRoute.test.tsx` passes
</acceptance_criteria>

---

### Task 6: Auth Pages
<read_first>
- apps/web/src/pages/LoginPage.tsx (create)
- apps/web/src/pages/RegisterPage.tsx (create)
- apps/web/src/pages/ResetPasswordPage.tsx (create)
- apps/web/src/components/auth/LoginForm.tsx
- apps/web/src/components/auth/RegisterForm.tsx
- apps/web/src/components/auth/ResetPasswordForm.tsx
</read_first>

<action>
Create page wrappers for auth forms:
- LoginPage: LoginForm + branding + links
- RegisterPage: RegisterForm + branding + links
- ResetPasswordPage: ResetPasswordForm + branding + links
- Consistent layout: centered card, mobile-first
- Add Harmoni branding/logo
- Add footer with copyright
</action>

<acceptance_criteria>
- LoginPage renders LoginForm with Harmoni branding
- RegisterPage renders RegisterForm with Harmoni branding
- ResetPasswordPage renders ResetPasswordForm with Harmoni branding
- All pages have links to other auth pages
- Pages are centered and responsive on mobile
- Footer displays copyright notice
- Visual consistency across all auth pages
- `npm test -- *Page.test.tsx` passes
</acceptance_criteria>

---

## Verification

1. **Login Page:**
   ```bash
   npm run dev
   # Open http://localhost:5173/login
   ```
   Expected: Login form with email/password fields, links to register and reset

2. **Registration Flow:**
   ```bash
   # Fill form with valid data
   # Submit
   ```
   Expected: Redirect to /login with success message

3. **Password Reset Flow:**
   ```bash
   # Visit /reset
   # Enter email
   # Submit
   ```
   Expected: "Reset link sent" message

4. **Protected Route:**
   ```bash
   # Navigate to /dashboard without login
   ```
   Expected: Redirect to /login

5. **Mobile Responsiveness:**
   ```bash
   # Chrome DevTools → Device Mode → 320px width
   ```
   Expected: All forms usable, touch targets >= 44px