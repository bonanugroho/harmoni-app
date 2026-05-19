import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import LoginPage from './LoginPage';

describe('LoginPage', () => {
  it('renders LoginForm with Harmoni branding', () => {
    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    expect(screen.getByText('Harmoni')).toBeInTheDocument();
    expect(screen.getByText('Community Financial Management')).toBeInTheDocument();
    expect(screen.getByText('Sign In to Harmoni')).toBeInTheDocument();
  });

  it('has links to other auth pages', () => {
    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    expect(screen.getByRole('link', { name: /create one/i })).toHaveAttribute('href', '/register');
    expect(screen.getByRole('link', { name: /forgot password/i })).toHaveAttribute('href', '/reset');
  });

  it('displays copyright footer', () => {
    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );

    expect(screen.getByText(/harmoni\. all rights reserved/i)).toBeInTheDocument();
  });
});
