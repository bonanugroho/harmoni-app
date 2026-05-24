import { describe, it, expect, vi, beforeEach } from 'vitest';
import { request } from './api';

vi.mock('./api', () => ({
  request: vi.fn(),
}));

import { listTenants, getTenant, createTenant, updateTenant, deleteTenant } from './tenants';
import type { CreateTenantRequest, UpdateTenantRequest } from '../types/tenant';

describe('tenants service', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('listTenants calls request /api/tenants and returns tenants array', async () => {
    const mockTenants = [{ id: '1', block: 'A', unit_number: '01' }];
    (request as ReturnType<typeof vi.fn>).mockResolvedValue({ tenants: mockTenants });

    const result = await listTenants();

    expect(request).toHaveBeenCalledWith('/api/tenants');
    expect(result).toEqual(mockTenants);
  });

  it('getTenant calls request /api/tenants/:id', async () => {
    const mockTenant = { id: '1', block: 'A', unit_number: '01' };
    (request as ReturnType<typeof vi.fn>).mockResolvedValue(mockTenant);

    const result = await getTenant('1');

    expect(request).toHaveBeenCalledWith('/api/tenants/1');
    expect(result).toEqual(mockTenant);
  });

  it('createTenant calls request /api/tenants with POST and body', async () => {
    const data: CreateTenantRequest = {
      block: 'A',
      unit_number: '01',
      occupancy: 'occupied',
      monthly_fee: 50000,
      mandatory_fees: [],
    };
    (request as ReturnType<typeof vi.fn>).mockResolvedValue({ id: '1', ...data });

    const result = await createTenant(data);

    expect(request).toHaveBeenCalledWith('/api/tenants', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    expect(result).toEqual({ id: '1', ...data });
  });

  it('updateTenant calls request /api/tenants/:id with PUT and body', async () => {
    const data: UpdateTenantRequest = { monthly_fee: 60000 };
    (request as ReturnType<typeof vi.fn>).mockResolvedValue({ id: '1', ...data });

    const result = await updateTenant('1', data);

    expect(request).toHaveBeenCalledWith('/api/tenants/1', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    expect(result).toEqual({ id: '1', ...data });
  });

  it('deleteTenant calls request /api/tenants/:id with DELETE', async () => {
    (request as ReturnType<typeof vi.fn>).mockResolvedValue(undefined);

    const result = await deleteTenant('1');

    expect(request).toHaveBeenCalledWith('/api/tenants/1', { method: 'DELETE' });
    expect(result).toBeUndefined();
  });
});
