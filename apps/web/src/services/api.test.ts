import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { request, API_BASE_URL } from './api';

describe('request helper', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('sends credentials: include on every request', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      status: 200,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ data: 'test' }),
    } as Response);

    await request('/test');

    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/test'),
      expect.objectContaining({ credentials: 'include' })
    );
  });

  it('prepends API_BASE_URL to relative URL', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      status: 200,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ data: 'test' }),
    } as Response);

    await request('/api/tenants');

    expect(fetch).toHaveBeenCalledWith(
      `${API_BASE_URL}/api/tenants`,
      expect.any(Object)
    );
  });

  it('returns undefined for 204 No Content responses', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      status: 204,
      headers: new Headers({}),
    } as Response);

    const result = await request('/api/tenants/1', { method: 'DELETE' });

    expect(result).toBeUndefined();
  });

  it('throws Error with API error message on 400+ JSON response', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: false,
      status: 400,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ error: 'Bad request', code: 'INVALID' }),
    } as Response);

    await expect(request('/api/tenants')).rejects.toThrow('Bad request');
  });

  it('throws generic error on 400+ non-JSON response', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: false,
      status: 500,
      headers: new Headers({ 'content-type': 'text/plain' }),
      text: () => Promise.resolve('Internal Server Error'),
    } as Response);

    await expect(request('/api/tenants')).rejects.toThrow(
      'Connection lost. Check your internet and try again.'
    );
  });

  it('returns parsed JSON for successful 200 responses', async () => {
    const mockData = { tenants: [{ id: '1', block: 'A' }] };
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      status: 200,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve(mockData),
    } as Response);

    const result = await request('/api/tenants');

    expect(result).toEqual(mockData);
  });

  it('merges custom options (method, body) correctly into request', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      status: 200,
      headers: new Headers({ 'content-type': 'application/json' }),
      json: () => Promise.resolve({ id: '1' }),
    } as Response);

    await request('/api/tenants', {
      method: 'POST',
      body: JSON.stringify({ block: 'A', unit_number: '01' }),
    });

    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/tenants'),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ block: 'A', unit_number: '01' }),
      })
    );
  });
});
