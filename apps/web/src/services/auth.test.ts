import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { login, register, requestPasswordReset, confirmPasswordReset } from '../services/auth';

function mockJsonResponse(overrides: Record<string, unknown> = {}) {
  return {
    ok: true,
    headers: new Headers({ 'content-type': 'application/json' }),
    json: () => Promise.resolve({}),
    ...overrides,
  };
}

describe('auth service', () => {
  beforeEach(() => {
    global.fetch = vi.fn();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('login', () => {
    it('sends POST to /auth/login with credentials: include', async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockJsonResponse({
        json: () => Promise.resolve({ user: { id: '1', email: 'test@test.com' } }),
      }));

      await login('test@test.com', 'password123');

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/login'),
        expect.objectContaining({
          method: 'POST',
          credentials: 'include',
          body: JSON.stringify({ email: 'test@test.com', password: 'password123' }),
        })
      );
    });

    it('returns user data on success', async () => {
      const mockUser = { id: '1', email: 'test@test.com' };
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockJsonResponse({
        json: () => Promise.resolve({ user: mockUser }),
      }));

      const result = await login('test@test.com', 'password123');

      expect(result).toEqual({ user: mockUser });
    });

    it('throws error on invalid credentials', async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockJsonResponse({
        ok: false,
        json: () => Promise.resolve({ error: 'Invalid email or password', code: 'INVALID_CREDENTIALS' }),
      }));

      await expect(login('test@test.com', 'wrong')).rejects.toThrow('Invalid email or password');
    });
  });

  describe('register', () => {
    it('sends POST to /auth/register with user data', async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockJsonResponse({
        json: () => Promise.resolve({ id: '1' }),
      }));

      await register({
        email: 'test@test.com',
        password: 'SecurePass123!',
        fullName: 'Test User',
        territoryId: 'rt-01',
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/register'),
        expect.objectContaining({
          method: 'POST',
          credentials: 'include',
          body: JSON.stringify({
            email: 'test@test.com',
            password: 'SecurePass123!',
            full_name: 'Test User',
            territory_id: 'rt-01',
          }),
        })
      );
    });

    it('throws error on duplicate email', async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockJsonResponse({
        ok: false,
        json: () => Promise.resolve({ error: 'Email already registered', code: 'DUPLICATE_EMAIL' }),
      }));

      await expect(
        register({ email: 'test@test.com', password: 'SecurePass123!', fullName: 'Test', territoryId: 'rt-01' })
      ).rejects.toThrow('Email already registered');
    });
  });

  describe('requestPasswordReset', () => {
    it('sends POST to /auth/reset', async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockJsonResponse({
        json: () => Promise.resolve({ message: 'Reset link sent' }),
      }));

      await requestPasswordReset('test@test.com');

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/reset'),
        expect.objectContaining({
          method: 'POST',
          credentials: 'include',
          body: JSON.stringify({ email: 'test@test.com' }),
        })
      );
    });
  });

  describe('confirmPasswordReset', () => {
    it('sends POST to /auth/reset/confirm', async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockJsonResponse({
        json: () => Promise.resolve({ message: 'Password updated' }),
      }));

      await confirmPasswordReset('token123', 'NewPass123!');

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/reset/confirm'),
        expect.objectContaining({
          method: 'POST',
          credentials: 'include',
          body: JSON.stringify({ token: 'token123', new_password: 'NewPass123!' }),
        })
      );
    });

    it('throws error on invalid token', async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValue(mockJsonResponse({
        ok: false,
        json: () => Promise.resolve({ error: 'Invalid or expired reset token', code: 'INVALID_TOKEN' }),
      }));

      await expect(confirmPasswordReset('expired', 'NewPass123!')).rejects.toThrow(
        'Invalid or expired reset token'
      );
    });
  });
});
