import { request } from './api';
import type { Tenant, CreateTenantRequest, UpdateTenantRequest } from '../types/tenant';

interface ListTenantsResponse {
  tenants: Tenant[];
}

export async function listTenants(): Promise<Tenant[]> {
  const data = await request<ListTenantsResponse>('/api/tenants');
  return data.tenants;
}

export async function getTenant(id: string): Promise<Tenant> {
  return request<Tenant>(`/api/tenants/${id}`);
}

export async function createTenant(data: CreateTenantRequest): Promise<Tenant> {
  return request<Tenant>('/api/tenants', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateTenant(id: string, data: UpdateTenantRequest): Promise<Tenant> {
  return request<Tenant>(`/api/tenants/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteTenant(id: string): Promise<void> {
  return request<void>(`/api/tenants/${id}`, { method: 'DELETE' });
}
