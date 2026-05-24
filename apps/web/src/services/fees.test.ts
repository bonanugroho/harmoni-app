import { describe, it, expect, vi, beforeEach } from 'vitest';
import { request } from './api';

vi.mock('./api', () => ({
  request: vi.fn(),
}));

import { listFees, createFee, updateFee, deleteFee } from './fees';
import type { CreateFeeRequest, UpdateFeeRequest } from '../types/fee';

describe('fees service', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('listFees calls request /api/tenants/:id/fees', async () => {
    const mockResponse = {
      mandatory_fees: [{ id: '1', amount: 25000, description: 'Security Fee' }],
      voluntary_fees: [],
    };
    (request as ReturnType<typeof vi.fn>).mockResolvedValue(mockResponse);

    const result = await listFees('1');

    expect(request).toHaveBeenCalledWith('/api/tenants/1/fees');
    expect(result).toEqual(mockResponse);
  });

  it('createFee calls request /api/tenants/:id/fees with POST and body', async () => {
    const data: CreateFeeRequest = {
      type: 'mandatory',
      amount: 25000,
      description: 'Security Fee',
      effective_date: '2026-06-01',
    };
    (request as ReturnType<typeof vi.fn>).mockResolvedValue({ id: '1', ...data });

    const result = await createFee('1', data);

    expect(request).toHaveBeenCalledWith('/api/tenants/1/fees', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    expect(result).toEqual({ id: '1', ...data });
  });

  it('updateFee calls request /api/tenants/:id/fees/:feeId with PUT and body', async () => {
    const data: UpdateFeeRequest = { amount: 30000 };
    (request as ReturnType<typeof vi.fn>).mockResolvedValue(undefined);

    const result = await updateFee('1', 'f1', data);

    expect(request).toHaveBeenCalledWith('/api/tenants/1/fees/f1', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    expect(result).toBeUndefined();
  });

  it('deleteFee calls request /api/tenants/:id/fees/:feeId with DELETE', async () => {
    (request as ReturnType<typeof vi.fn>).mockResolvedValue(undefined);

    const result = await deleteFee('1', 'f1');

    expect(request).toHaveBeenCalledWith('/api/tenants/1/fees/f1', { method: 'DELETE' });
    expect(result).toBeUndefined();
  });
});
