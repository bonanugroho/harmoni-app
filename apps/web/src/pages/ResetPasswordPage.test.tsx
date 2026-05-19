import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ResetPasswordPage from './ResetPasswordPage';

describe('ResetPasswordPage', () => {
  it('renders ResetPasswordForm with Harmoni branding', () => {
    render(
      <MemoryRouter>
        <ResetPasswordPage />
      </MemoryRouter>
    );

    expect(screen.getByText('Harmoni')).toBeInTheDocument();
    expect(screen.getByText('Community Financial Management')).toBeInTheDocument();
    expect(screen.getByText('Reset Your Password')).toBeInTheDocument();
  });

  it('has links to other auth pages', () => {
    render(
      <MemoryRouter>
        <ResetPasswordPage />
      </MemoryRouter>
    );

    expect(screen.getByRole('link', { name: /back to sign in/i })).toHaveAttribute('href', '/login');
  });

  it('displays copyright footer', () => {
    render(
      <MemoryRouter>
        <ResetPasswordPage />
      </MemoryRouter>
    );

    expect(screen.getByText(/harmoni\. all rights reserved/i)).toBeInTheDocument();
  });
});
