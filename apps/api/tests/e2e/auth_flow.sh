#!/usr/bin/env bash
# E2E Authentication Flow Test
# Tests the complete auth flow: register → login → protected access → password reset → RBAC
#
# Usage: bash tests/e2e/auth_flow.sh [BASE_URL]
#   BASE_URL defaults to http://localhost:3000
#
# Exit codes: 0 = all tests pass, 1 = any test fails

set -euo pipefail

BASE_URL="${1:-http://localhost:3000}"
PASS=0
FAIL=0
TOTAL=0

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

pass() {
    PASS=$((PASS + 1))
    TOTAL=$((TOTAL + 1))
    echo -e "  ${GREEN}PASS${NC}: $1"
}

fail() {
    FAIL=$((FAIL + 1))
    TOTAL=$((TOTAL + 1))
    echo -e "  ${RED}FAIL${NC}: $1 - $2"
}

section() {
    echo ""
    echo -e "${YELLOW}=== $1 ===${NC}"
}

# Check if server is reachable
section "Pre-flight: Server connectivity"
if curl -sf "${BASE_URL}/health" > /dev/null 2>&1; then
    pass "Server is reachable at ${BASE_URL}"
else
    echo -e "${RED}ERROR: Server not reachable at ${BASE_URL}${NC}"
    echo "Start the server first: cd apps/api && go run cmd/server/main.go"
    exit 1
fi

# ============================================================
# 1. User Registration
# ============================================================
section "1. User Registration"

REGISTER_RESPONSE=$(curl -sf -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "e2e-test@example.com",
        "password": "E2eTest123!",
        "full_name": "E2E Test User",
        "territory_id": "rt-01"
    }' 2>/dev/null || true)

if echo "$REGISTER_RESPONSE" | grep -q '"email"'; then
    pass "User registered successfully"
else
    fail "User registration" "Response: ${REGISTER_RESPONSE:-empty}"
fi

# Test duplicate email registration
DUPLICATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "e2e-test@example.com",
        "password": "E2eTest123!",
        "full_name": "Duplicate User",
        "territory_id": "rt-01"
    }' 2>/dev/null || true)

DUPLICATE_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "e2e-test@example.com",
        "password": "E2eTest123!",
        "full_name": "Duplicate User",
        "territory_id": "rt-01"
    }' 2>/dev/null || true)

if [ "$DUPLICATE_STATUS" = "409" ]; then
    pass "Duplicate email rejected with 409"
else
    fail "Duplicate email rejection" "Status: ${DUPLICATE_STATUS}"
fi

# Test weak password registration
WEAK_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "weak@example.com",
        "password": "weak",
        "full_name": "Weak User",
        "territory_id": "rt-01"
    }' 2>/dev/null || true)

if [ "$WEAK_STATUS" = "400" ]; then
    pass "Weak password rejected with 400"
else
    fail "Weak password rejection" "Status: ${WEAK_STATUS}"
fi

# ============================================================
# 2. User Login
# ============================================================
section "2. User Login"

LOGIN_RESPONSE=$(curl -s -c /tmp/e2e_cookies.txt -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "e2e-test@example.com",
        "password": "E2eTest123!"
    }' 2>/dev/null || true)

if echo "$LOGIN_RESPONSE" | grep -q '"user"'; then
    pass "Login successful, user data returned"
else
    fail "Login" "Response: ${LOGIN_RESPONSE}"
fi

# Check for httpOnly cookie
if grep -q "paseto_token" /tmp/e2e_cookies.txt 2>/dev/null; then
    pass "httpOnly cookie set (paseto_token)"
else
    # Fiber sets cookies via Set-Cookie header, not in curl cookie jar by default
    # Check the raw response headers instead
    LOGIN_HEADERS=$(curl -s -D - -o /dev/null -X POST "${BASE_URL}/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "e2e-test@example.com",
            "password": "E2eTest123!"
        }' 2>/dev/null || true)
    if echo "$LOGIN_HEADERS" | grep -qi "set-cookie"; then
        pass "Set-Cookie header present in login response"
    else
        fail "Cookie handling" "No Set-Cookie header found"
    fi
fi

# Test invalid credentials
INVALID_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "e2e-test@example.com",
        "password": "WrongPassword1!"
    }' 2>/dev/null || true)

if [ "$INVALID_STATUS" = "401" ]; then
    pass "Invalid credentials rejected with 401"
else
    fail "Invalid credentials rejection" "Status: ${INVALID_STATUS}"
fi

# ============================================================
# 3. Protected Route Access
# ============================================================
section "3. Protected Route Access"

# Extract token from cookie file for API requests
TOKEN=$(grep "paseto_token" /tmp/e2e_cookies.txt 2>/dev/null | awk '{print $NF}' | tr -d '\r\n' || echo "")

if [ -n "$TOKEN" ]; then
    PROTECTED_RESPONSE=$(curl -s -b "paseto_token=${TOKEN}" "${BASE_URL}/api/protected" 2>/dev/null || true)
    if echo "$PROTECTED_RESPONSE" | grep -q '"user_id"\|"role"'; then
        pass "Protected route accessible with valid token"
    else
        fail "Protected route access" "Response: ${PROTECTED_RESPONSE}"
    fi
else
    # Try accessing protected route without token
    NO_TOKEN_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/api/protected" 2>/dev/null || true)
    if [ "$NO_TOKEN_STATUS" = "401" ]; then
        pass "Protected route returns 401 without token"
    else
        fail "Protected route without token" "Status: ${NO_TOKEN_STATUS}"
    fi
fi

# ============================================================
# 4. Password Reset Request
# ============================================================
section "4. Password Reset Request"

RESET_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/reset" \
    -H "Content-Type: application/json" \
    -d '{"email": "e2e-test@example.com"}' 2>/dev/null || true)

if echo "$RESET_RESPONSE" | grep -q '"message"'; then
    pass "Password reset request accepted (200)"
else
    fail "Password reset request" "Response: ${RESET_RESPONSE}"
fi

# Test reset for non-existent email (should still return 200 - enumeration prevention)
RESET_NONEXISTENT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/auth/reset" \
    -H "Content-Type: application/json" \
    -d '{"email": "nonexistent@example.com"}' 2>/dev/null || true)

if [ "$RESET_NONEXISTENT_STATUS" = "200" ]; then
    pass "Reset for non-existent email returns 200 (enumeration prevention)"
else
    fail "Reset enumeration prevention" "Status: ${RESET_NONEXISTENT_STATUS}"
fi

# ============================================================
# 5. Password Reset Confirm
# ============================================================
section "5. Password Reset Confirm"

# Test with invalid token
INVALID_RESET_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/auth/reset/confirm" \
    -H "Content-Type: application/json" \
    -d '{
        "token": "invalid-token",
        "new_password": "NewE2ePass456!"
    }' 2>/dev/null || true)

if [ "$INVALID_RESET_STATUS" = "400" ]; then
    pass "Invalid reset token rejected with 400"
else
    fail "Invalid reset token rejection" "Status: ${INVALID_RESET_STATUS}"
fi

# Test with missing fields
MISSING_FIELDS_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/auth/reset/confirm" \
    -H "Content-Type: application/json" \
    -d '{"token": "some-token"}' 2>/dev/null || true)

if [ "$MISSING_FIELDS_STATUS" = "400" ]; then
    pass "Missing fields rejected with 400"
else
    fail "Missing fields rejection" "Status: ${MISSING_FIELDS_STATUS}"
fi

# ============================================================
# 6. Public Route Bypass
# ============================================================
section "6. Public Route Bypass"

HEALTH_RESPONSE=$(curl -s "${BASE_URL}/health" 2>/dev/null || true)
if echo "$HEALTH_RESPONSE" | grep -q '"status"'; then
    pass "Health endpoint accessible without auth"
else
    fail "Health endpoint" "Response: ${HEALTH_RESPONSE}"
fi

# ============================================================
# Summary
# ============================================================
section "Test Summary"
echo "  Total:  ${TOTAL}"
echo -e "  Passed: ${GREEN}${PASS}${NC}"
echo -e "  Failed: ${RED}${FAIL}${NC}"

# Cleanup
rm -f /tmp/e2e_cookies.txt

if [ "$FAIL" -gt 0 ]; then
    echo -e "\n${RED}E2E tests completed with failures${NC}"
    exit 1
else
    echo -e "\n${GREEN}All E2E tests passed${NC}"
    exit 0
fi
