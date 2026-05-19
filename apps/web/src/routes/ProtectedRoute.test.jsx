import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import ProtectedRoute, { AuthContext, useAuth } from './ProtectedRoute';

// Mock fetch
global.fetch = vi.fn();

function renderWithRouter(ui, { initialEntries = ['/protected'] } = {}) {
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
    fetch.mockResolvedValue({ ok: false });

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
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ user: { id: '1', email: 'test@test.com', role: 'resident' } }),
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
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ user: { id: '1', email: 'test@test.com', role: 'resident' } }),
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
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ user: { id: '1', email: 'test@test.com', role: 'resident' } }),
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
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ user: { id: '1', email: 'test@test.com', role: 'rt_officer' } }),
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
    // Return a promise that never resolves to simulate loading
    fetch.mockReturnValue(new Promise(() => {}));

    renderWithRouter(
      <ProtectedRoute>
        <div>Protected Content</div>
      </ProtectedRoute>
    );

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });
});
