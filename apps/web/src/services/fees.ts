import { request } from './api';
import type { Fee, CreateFeeRequest, UpdateFeeRequest, ListFeesResponse } from '../types/fee';

export async function listFees(tenantId: string): Promise<ListFeesResponse> {
  return request<ListFeesResponse>(`/api/tenants/${tenantId}/fees`);
}

export async function createFee(tenantId: string, data: CreateFeeRequest): Promise<Fee> {
  return request<Fee>(`/api/tenants/${tenantId}/fees`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateFee(tenantId: string, feeId: string, data: UpdateFeeRequest): Promise<void> {
  return request<void>(`/api/tenants/${tenantId}/fees/${feeId}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteFee(tenantId: string, feeId: string): Promise<void> {
  return request<void>(`/api/tenants/${tenantId}/fees/${feeId}`, { method: 'DELETE' });
}
