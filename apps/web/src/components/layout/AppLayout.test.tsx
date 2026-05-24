import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import AppLayout from './AppLayout';

// Mock useAuth from ProtectedRoute
vi.mock('../../routes/ProtectedRoute', async () => {
  const actual = await vi.importActual('../../routes/ProtectedRoute');
  return {
    ...actual,
    useAuth: vi.fn(),
  };
});

// Mock react-router-dom hooks
const mockNavigate = vi.fn();
let mockPathname = '/tenants';

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useLocation: () => ({ pathname: mockPathname, search: '', hash: '', state: null }),
  };
});

import { useAuth } from '../../routes/ProtectedRoute';

function renderAppLayout() {
  return render(
    <MemoryRouter>
      <AppLayout>
        <div data-testid="children-content">Page Content</div>
      </AppLayout>
    </MemoryRouter>
  );
}

describe('AppLayout', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockPathname = '/tenants';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
      user: { role: 'rt_officer' },
      isAuthenticated: true,
    });
  });

  it('renders children content', () => {
    renderAppLayout();
    expect(screen.getByTestId('children-content')).toBeInTheDocument();
    expect(screen.getByText('Page Content')).toBeInTheDocument();
  });

  it('shows sidebar navigation links (Tenants, Reports)', () => {
    renderAppLayout();
    expect(screen.getByText('Tenants')).toBeInTheDocument();
    expect(screen.getByText('Reports')).toBeInTheDocument();
  });

  it('highlights Tenants link when on /tenants path', () => {
    renderAppLayout();
    const tenantsButton = screen.getByText('Tenants').closest('button');
    expect(tenantsButton?.className).toContain('bg-blue-100');
    expect(tenantsButton?.className).toContain('text-blue-700');
    expect(tenantsButton?.className).toContain('border-blue-600');
  });

  it('shows Settings link for rt_officer role', () => {
    renderAppLayout();
    expect(screen.getByText('Settings')).toBeInTheDocument();
  });

  it('hides Settings link for resident role', () => {
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
      user: { role: 'resident' },
      isAuthenticated: true,
    });
    renderAppLayout();
    expect(screen.queryByText('Settings')).not.toBeInTheDocument();
  });

  it('mobile hamburger toggle opens sidebar', () => {
    renderAppLayout();
    const toggleButton = screen.getByLabelText('Toggle sidebar');
    expect(toggleButton).toBeInTheDocument();
    fireEvent.click(toggleButton);
    // Sidebar becomes visible (translate-x-0 class)
    const sidebar = document.querySelector('aside');
    expect(sidebar?.className).toContain('translate-x-0');
  });
});
