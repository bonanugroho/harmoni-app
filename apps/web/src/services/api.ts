export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface ApiError {
  error: string;
  code?: string;
}

export type RequestOptions = RequestInit;

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
