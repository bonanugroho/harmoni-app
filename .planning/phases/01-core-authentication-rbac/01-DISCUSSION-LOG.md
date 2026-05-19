# Phase 1 Discussion Log

**Date:** 2026-05-19

## Areas Discussed

1. **Session & Token Storage**
   - Decision: httpOnly cookies (server‑managed, XSS‑protected)

2. **Password Reset Flow**
   - Decision: Email‑based one‑time‑token links

3. **Territory Model**
   - Decision: One user = one role + one territory
   - Clarified that users cannot hold multiple roles across territories

4. **Role Assignment**
   - Decision: Admin dashboard for creating RT/RW Officer accounts

5. **Password Policy**
   - Decision: Complexity rules (uppercase, lowercase, numbers, symbols)

6. **Role Names / Display Language**
   - Decision: English only for display names

7. **RW Officer Access**
   - Decision: Automatic access to all RT data under their jurisdiction

8. **Read‑Only Roles**
   - Decision: No separate auditor/observer roles needed

9. **Casbin Policy Structure**
   - Presented three options (per‑role, per‑resource, hybrid)
   - User selected **Hybrid** approach with resource‑based policies and `{{territory_id}}` placeholders for RT officers and `*` wildcard for RW officers

10. **Database Migration Strategy**
    - Compared Goose vs golang‑migrate
    - User chose **golang‑migrate CLI** (SQL‑only, versioned, up/down files)

## Summary
All gray areas for Phase 1 have been resolved. The decisions are captured in `01-CONTEXT.md`. No further open items remain for this phase.
