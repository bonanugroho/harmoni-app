# Requirements Document - Harmoni Project

## 1. Introduction
Harmoni app is a community financial management application designed for neighborhood scales (Rukun Tetangga/RT and Rukun Warga/RW). The primary focus of this application is to create transparency and accountability in income and expenditure reporting for all stakeholders (Residents, RT Officers, and RW Officers).

## 2. Functional Requirements
### 2.1 User Management & Access
- The system must support authentication (Login, Registration, Password Reset).
- The system must implement multi-level Role-Based Access Control (RBAC):
    - **Resident:** Can view the dashboard and transaction history for their specific household only.
    - **RT Officer:** Can manage tenant data, income, and expenditures within their specific RT scope.
    - **RW Officer:** Can manage RW-level income/expenditures and monitor the financial health of all RTs under their jurisdiction.

### 2.2 Data & Transaction Management
- **Tenant Management:** Record keeping for houses (block & number), occupancy status, and fixed monthly fees.
- **Income (Inflow):**
    - Mandatory Fees (Fixed monthly fees per unit, e.g., waste management, security).
    - Voluntary Fees (Incidental contributions, e.g., holiday bonuses/THR, social donations).
    - RT to RW Deposits (Transfers from RT cash to RW cash).
- **Expenditures (Outflow):**
    - Recording operational costs (Security salaries, cleaning services, public facility maintenance).

### 2.3 Reporting & Dashboard
- Dashboards must display real-time cash balances.
- Visualizations of historical income vs. expenditure trends.
- **Accounts Receivable Analysis (RT):** Monitoring late fee payments categorized by:
    - On-time (Lancar)
    - Late > 30 days
    - Late > 60 days
    - Late > 90 days

## 3. Non-Functional Requirements
- **Responsiveness:** The web-based application must be optimized for mobile browsers (Mobile-First).
- **Security:** Use modern security token standards (PASETO) and strict API-level authorization.
- **Data Integrity:** Strict data isolation between territories (e.g., RT 01 officers cannot access RT 02 details).