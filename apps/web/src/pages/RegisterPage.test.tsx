import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import RegisterPage from './RegisterPage';

describe('RegisterPage', () => {
  it('renders RegisterForm with Harmoni branding', () => {
    render(
      <MemoryRouter>
        <RegisterPage />
      </MemoryRouter>
    );

    expect(screen.getByText('Harmoni')).toBeInTheDocument();
    expect(screen.getByText('Community Financial Management')).toBeInTheDocument();
    expect(screen.getByText('Create Your Account')).toBeInTheDocument();
  });

  it('has links to other auth pages', () => {
    render(
      <MemoryRouter>
        <RegisterPage />
      </MemoryRouter>
    );

    expect(screen.getByRole('link', { name: /sign in/i })).toHaveAttribute('href', '/login');
  });

  it('displays copyright footer', () => {
    render(
      <MemoryRouter>
        <RegisterPage />
      </MemoryRouter>
    );

    expect(screen.getByText(/harmoni\. all rights reserved/i)).toBeInTheDocument();
  });
});
