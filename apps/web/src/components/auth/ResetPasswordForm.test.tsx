import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ResetPasswordForm from './ResetPasswordForm';
import * as auth from '../../services/auth';

vi.mock('../../services/auth', () => ({
  requestPasswordReset: vi.fn(),
  confirmPasswordReset: vi.fn(),
}));

function renderResetForm(initialEntries = ['/reset']) {
  return render(
    <MemoryRouter initialEntries={initialEntries}>
      <ResetPasswordForm />
    </MemoryRouter>
  );
}

describe('ResetPasswordForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('request reset mode (no token)', () => {
    it('shows email input initially', () => {
      renderResetForm();

      expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /send reset link/i })).toBeInTheDocument();
    });

    it('submitting email shows "Reset link sent" message', async () => {
      (auth.requestPasswordReset as ReturnType<typeof vi.fn>).mockResolvedValue({ message: 'Reset link sent' });
      renderResetForm();

      fireEvent.change(screen.getByLabelText(/email/i), {
        target: { value: 'test@test.com' },
      });
      fireEvent.click(screen.getByRole('button', { name: /send reset link/i }));

      await waitFor(() => {
        expect(
          screen.getByText(/check your email. we sent a reset link to test@test.com/i)
        ).toBeInTheDocument();
      });
    });

    it('has link back to login page', () => {
      renderResetForm();

      expect(screen.getByRole('link', { name: /back to sign in/i })).toHaveAttribute('href', '/login');
    });
  });

  describe('new password mode (with token)', () => {
    it('shows token + password form when URL contains ?token= parameter', () => {
      renderResetForm(['/reset?token=test-token-123']);

      expect(screen.getByLabelText(/new password/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /update password/i })).toBeInTheDocument();
    });

    it('password validation same as registration form', async () => {
      renderResetForm(['/reset?token=test-token-123']);

      fireEvent.change(screen.getByLabelText(/new password/i), {
        target: { value: 'weak' },
      });
      fireEvent.change(screen.getByLabelText(/confirm password/i), {
        target: { value: 'weak' },
      });
      fireEvent.click(screen.getByRole('button', { name: /update password/i }));

      await waitFor(() => {
        expect(screen.getByText(/password must be at least 8 characters/i)).toBeInTheDocument();
      });
    });

    it('successful reset redirects to /login', async () => {
      (auth.confirmPasswordReset as ReturnType<typeof vi.fn>).mockResolvedValue({ message: 'Password updated' });
      renderResetForm(['/reset?token=test-token-123']);

      fireEvent.change(screen.getByLabelText(/new password/i), {
        target: { value: 'SecurePass123!' },
      });
      fireEvent.change(screen.getByLabelText(/confirm password/i), {
        target: { value: 'SecurePass123!' },
      });
      fireEvent.click(screen.getByRole('button', { name: /update password/i }));

      await waitFor(() => {
        expect(screen.getByText(/password updated. you can now sign in/i)).toBeInTheDocument();
      });
    });

    it('invalid/expired token shows "Invalid or expired reset token" error', async () => {
      (auth.confirmPasswordReset as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Invalid or expired reset token'));
      renderResetForm(['/reset?token=expired-token']);

      fireEvent.change(screen.getByLabelText(/new password/i), {
        target: { value: 'SecurePass123!' },
      });
      fireEvent.change(screen.getByLabelText(/confirm password/i), {
        target: { value: 'SecurePass123!' },
      });

      const form = screen.getByRole('button', { name: /update password/i });
      fireEvent.click(form);

      await waitFor(
        () => {
          expect(screen.getByRole('alert')).toHaveTextContent('Invalid or expired reset token');
        },
        { timeout: 3000 }
      );
    });
  });
});
