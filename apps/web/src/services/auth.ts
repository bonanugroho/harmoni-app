import type { ApiResponse, LoginResponse, RegisterRequest } from '../types/auth';

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

export async function login(email: string, password: string): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ email, password }),
  });

  return parseError<LoginResponse>(response);
}

export async function register(userData: RegisterRequest): Promise<ApiResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({
      email: userData.email,
      password: userData.password,
      full_name: userData.fullName,
      territory_id: userData.territoryId,
    }),
  });

  return parseError<ApiResponse>(response);
}

export async function requestPasswordReset(email: string): Promise<ApiResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/reset`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ email }),
  });

  return parseError<ApiResponse>(response);
}

export async function confirmPasswordReset(token: string, newPassword: string): Promise<ApiResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/reset/confirm`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ token, new_password: newPassword }),
  });

  return parseError<ApiResponse>(response);
}

export function logout(): void {
  if (typeof window !== 'undefined') {
    window.location.href = '/login';
  }
}
