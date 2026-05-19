const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

/**
 * Auth service - API client for authentication endpoints.
 * Uses fetch with credentials: 'include' for httpOnly cookie handling.
 */

/**
 * Parse error response from API.
 * API returns: { "error": "...", "code": "..." }
 */
function parseError(response) {
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

/**
 * Login with email and password.
 * @param {string} email
 * @param {string} password
 * @returns {Promise<{user: object}>}
 */
export async function login(email, password) {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ email, password }),
  });

  return parseError(response);
}

/**
 * Register a new user.
 * @param {object} userData
 * @param {string} userData.email
 * @param {string} userData.password
 * @param {string} userData.fullName
 * @param {string} userData.territoryId
 * @returns {Promise<object>}
 */
export async function register({ email, password, fullName, territoryId }) {
  const response = await fetch(`${API_BASE_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({
      email,
      password,
      full_name: fullName,
      territory_id: territoryId,
    }),
  });

  return parseError(response);
}

/**
 * Request a password reset link.
 * @param {string} email
 * @returns {Promise<{message: string}>}
 */
export async function requestPasswordReset(email) {
  const response = await fetch(`${API_BASE_URL}/auth/reset`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ email }),
  });

  return parseError(response);
}

/**
 * Confirm password reset with token and new password.
 * @param {string} token
 * @param {string} newPassword
 * @returns {Promise<{message: string}>}
 */
export async function confirmPasswordReset(token, newPassword) {
  const response = await fetch(`${API_BASE_URL}/auth/reset/confirm`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ token, new_password: newPassword }),
  });

  return parseError(response);
}

/**
 * Logout - clear local state.
 * Note: Server manages httpOnly cookie expiry.
 * Client clears any local state (e.g., context).
 */
export function logout() {
  // Server handles cookie removal via Set-Cookie with MaxAge=0
  // Client can clear any local state here if needed
  if (typeof window !== 'undefined') {
    // Clear any client-side user data
    window.location.href = '/login';
  }
}
