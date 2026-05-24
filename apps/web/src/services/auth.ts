import { request } from './api';
import type { ApiResponse, LoginResponse, RegisterRequest } from '../types/auth';

export async function login(email: string, password: string): Promise<LoginResponse> {
  return request<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export async function register(userData: RegisterRequest): Promise<ApiResponse> {
  return request<ApiResponse>('/auth/register', {
    method: 'POST',
    body: JSON.stringify({
      email: userData.email,
      password: userData.password,
      full_name: userData.fullName,
      territory_id: userData.territoryId,
    }),
  });
}

export async function requestPasswordReset(email: string): Promise<ApiResponse> {
  return request<ApiResponse>('/auth/reset', {
    method: 'POST',
    body: JSON.stringify({ email }),
  });
}

export async function confirmPasswordReset(token: string, newPassword: string): Promise<ApiResponse> {
  return request<ApiResponse>('/auth/reset/confirm', {
    method: 'POST',
    body: JSON.stringify({ token, new_password: newPassword }),
  });
}

export function logout(): void {
  if (typeof window !== 'undefined') {
    window.location.href = '/login';
  }
}
