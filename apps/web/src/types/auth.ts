export interface User {
  id: string;
  email: string;
  full_name?: string;
  role?: string;
  territory_id?: string;
}

export interface LoginResponse {
  user: User;
}

export interface RegisterRequest {
  email: string;
  password: string;
  fullName: string;
  territoryId: string;
}

export interface ResetRequest {
  email: string;
}

export interface ConfirmResetRequest {
  token: string;
  new_password: string;
}

export interface ApiResponse<T = unknown> {
  user?: T;
  message?: string;
  error?: string;
  code?: string;
}
