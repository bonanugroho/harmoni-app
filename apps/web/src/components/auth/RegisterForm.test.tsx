import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import RegisterForm from './RegisterForm';
import * as auth from '../../services/auth';

vi.mock('../../services/auth', () => ({
  register: vi.fn(),
}));

function renderRegisterForm() {
  return render(
    <MemoryRouter>
      <RegisterForm />
    </MemoryRouter>
  );
}

describe('RegisterForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders with all required fields', () => {
    renderRegisterForm();

    expect(screen.getByLabelText(/full name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^password$/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/territory/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /create account/i })).toBeInTheDocument();
  });

  it('populates territory dropdown with available territories', () => {
    renderRegisterForm();

    const select = screen.getByLabelText(/territory/i);
    expect(select).toBeInTheDocument();

    fireEvent.mouseDown(select);
    expect(screen.getByText('RT 01')).toBeInTheDocument();
    expect(screen.getByText('RT 02')).toBeInTheDocument();
    expect(screen.getByText('RW 01')).toBeInTheDocument();
  });

  it('shows password strength indicator updates in real-time', () => {
    renderRegisterForm();

    const passwordInput = screen.getByLabelText(/^password$/i);

    fireEvent.change(passwordInput, { target: { value: 'password' } });
    expect(screen.getByText(/password strength: weak/i)).toBeInTheDocument();
    expect(screen.getByText(/password strength: weak/i)).toHaveClass('text-red-600');

    fireEvent.change(passwordInput, { target: { value: 'SecurePass123!' } });
    expect(screen.getByText(/password strength: strong/i)).toBeInTheDocument();
    expect(screen.getByText(/password strength: strong/i)).toHaveClass('text-green-600');
  });

  it('weak password shows strength "weak" in red', () => {
    renderRegisterForm();

    fireEvent.change(screen.getByLabelText(/^password$/i), {
      target: { value: 'password' },
    });

    const strengthText = screen.getByText(/password strength: weak/i);
    expect(strengthText).toBeInTheDocument();
    expect(strengthText).toHaveClass('text-red-600');
  });

  it('strong password shows strength "strong" in green', () => {
    renderRegisterForm();

    fireEvent.change(screen.getByLabelText(/^password$/i), {
      target: { value: 'SecurePass123!' },
    });

    const strengthText = screen.getByText(/password strength: strong/i);
    expect(strengthText).toBeInTheDocument();
    expect(strengthText).toHaveClass('text-green-600');
  });

  it('redirects to /login on successful registration', async () => {
    (auth.register as ReturnType<typeof vi.fn>).mockResolvedValue({ id: '1' });
    renderRegisterForm();

    fireEvent.change(screen.getByLabelText(/full name/i), {
      target: { value: 'Test User' },
    });
    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: 'test@test.com' },
    });
    fireEvent.change(screen.getByLabelText(/^password$/i), {
      target: { value: 'SecurePass123!' },
    });
    fireEvent.change(screen.getByLabelText(/confirm password/i), {
      target: { value: 'SecurePass123!' },
    });
    fireEvent.change(screen.getByLabelText(/territory/i), {
      target: { value: 'rt-01' },
    });
    fireEvent.click(screen.getByRole('button', { name: /create account/i }));

    await waitFor(() => {
      expect(auth.register).toHaveBeenCalledWith({
        email: 'test@test.com',
        password: 'SecurePass123!',
        fullName: 'Test User',
        territoryId: 'rt-01',
      });
    });
  });

  it('shows "Email already registered" error for duplicate email', async () => {
    (auth.register as ReturnType<typeof vi.fn>).mockRejectedValue(new Error('Email already registered'));
    renderRegisterForm();

    fireEvent.change(screen.getByLabelText(/full name/i), {
      target: { value: 'Test User' },
    });
    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: 'test@test.com' },
    });
    fireEvent.change(screen.getByLabelText(/^password$/i), {
      target: { value: 'SecurePass123!' },
    });
    fireEvent.change(screen.getByLabelText(/confirm password/i), {
      target: { value: 'SecurePass123!' },
    });
    fireEvent.change(screen.getByLabelText(/territory/i), {
      target: { value: 'rt-01' },
    });
    fireEvent.click(screen.getByRole('button', { name: /create account/i }));

    await waitFor(() => {
      expect(
        screen.getByText('An account with this email already exists. Try signing in instead.')
      ).toBeInTheDocument();
    });
  });

  it('has link to login page', () => {
    renderRegisterForm();

    expect(screen.getByRole('link', { name: /sign in/i })).toHaveAttribute('href', '/login');
  });
});
