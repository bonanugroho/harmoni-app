import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import ProtectedRoute, { useAuth } from './ProtectedRoute';

vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual('@tanstack/react-query');
  return {
    ...actual,
    useQuery: vi.fn(),
  };
});

vi.mock('../services/api', () => ({
  request: vi.fn(),
}));

import { useQuery } from '@tanstack/react-query';
import type { Mock } from 'vitest';

function renderWithRouter(ui: React.ReactNode, { initialEntries = ['/protected'] }: { initialEntries?: string[] } = {}) {
  return render(
    <MemoryRouter initialEntries={initialEntries}>
      <Routes>
        <Route path="/protected" element={ui} />
        <Route path="/login" element={<div>Login Page</div>} />
      </Routes>
    </MemoryRouter>
  );
}

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('redirects to /login when not authenticated', async () => {
    (useQuery as Mock).mockReturnValue({ data: undefined, isLoading: false, isError: true });

    renderWithRouter(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByText('Login Page')).toBeInTheDocument();
    });
  });

  it('renders children when authenticated', async () => {
    (useQuery as Mock).mockReturnValue({
      data: { user: { id: '1', email: 'test@test.com', role: 'resident' } },
      isLoading: false,
      isError: false,
    });

    renderWithRouter(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
  });

  it('passes user data to children via context', async () => {
    (useQuery as Mock).mockReturnValue({
      data: { user: { id: '1', email: 'test@test.com', role: 'resident' } },
      isLoading: false,
      isError: false,
    });

    function TestConsumer() {
      const { user } = useAuth();
      return <div>User: {user?.email}</div>;
    }

    renderWithRouter(
      <ProtectedRoute>
        <TestConsumer />
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByText('User: test@test.com')).toBeInTheDocument();
    });
  });

  it('restricts access by role when requiredRole prop is set', async () => {
    (useQuery as Mock).mockReturnValue({
      data: { user: { id: '1', email: 'test@test.com', role: 'resident' } },
      isLoading: false,
      isError: false,
    });

    renderWithRouter(
      <ProtectedRoute requiredRole="rt_officer">
        <div>Protected Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByText('Access Denied')).toBeInTheDocument();
    });
  });

  it('allows access when user has required role', async () => {
    (useQuery as Mock).mockReturnValue({
      data: { user: { id: '1', email: 'test@test.com', role: 'rt_officer' } },
      isLoading: false,
      isError: false,
    });

    renderWithRouter(
      <ProtectedRoute requiredRole="rt_officer">
        <div>Protected Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
  });

  it('shows loading spinner while checking auth status', () => {
    (useQuery as Mock).mockReturnValue({ data: undefined, isLoading: true, isError: false });

    renderWithRouter(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>
    );

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });
});
