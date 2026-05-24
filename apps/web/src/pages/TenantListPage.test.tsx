import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import TenantListPage from './TenantListPage';

vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual('@tanstack/react-query');
  return {
    ...actual,
    useQuery: vi.fn(),
    useMutation: vi.fn(),
  };
});

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

import { useQuery } from '@tanstack/react-query';

const mockTenants = [
  {
    id: '1',
    block: 'A',
    unit_number: '01',
    occupancy: 'occupied' as const,
    monthly_fee: 50000,
    territory_id: 't1',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
  },
  {
    id: '2',
    block: 'B',
    unit_number: '10',
    occupancy: 'vacant' as const,
    monthly_fee: 75000,
    territory_id: 't1',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
  },
];

function renderPage() {
  return render(
    <MemoryRouter>
      <TenantListPage />
    </MemoryRouter>
  );
}

describe('TenantListPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows LoadingSkeleton when isLoading is true', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
      error: null,
    });

    renderPage();
    const skeletons = screen.getAllByTestId('loading-skeleton');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders TenantCards when data loads', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenants,
      isLoading: false,
      isError: false,
      error: null,
    });

    renderPage();
    expect(screen.getByText(/Block A/i)).toBeInTheDocument();
    expect(screen.getByText(/Block B/i)).toBeInTheDocument();
    expect(screen.getByText(/Unit 01/i)).toBeInTheDocument();
    expect(screen.getByText(/Unit 10/i)).toBeInTheDocument();
  });

  it('shows EmptyState when data is null', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: null,
      isLoading: false,
      isError: false,
      error: null,
    });

    renderPage();
    expect(screen.getByText('No Tenants Yet')).toBeInTheDocument();
  });

  it('shows EmptyState when data is empty array', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: [],
      isLoading: false,
      isError: false,
      error: null,
    });

    renderPage();
    expect(screen.getByText('No Tenants Yet')).toBeInTheDocument();
    expect(
      screen.getByText('Start by adding your first tenant to this RT.')
    ).toBeInTheDocument();
  });

  it('error state shows error message', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
      error: new Error('Failed to load tenants'),
    });

    renderPage();
    expect(
      screen.getByText('Failed to load tenants')
    ).toBeInTheDocument();
  });

  it('search filters by block name', async () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenants,
      isLoading: false,
      isError: false,
      error: null,
    });

    renderPage();

    const searchInput = screen.getByPlaceholderText(/search by block/i);
    fireEvent.change(searchInput, { target: { value: 'A' } });

    // Wait for debounce (300ms)
    await waitFor(
      () => {
        expect(screen.getByText(/Block A/i)).toBeInTheDocument();
        expect(screen.queryByText(/Block B/i)).not.toBeInTheDocument();
      },
      { timeout: 500 }
    );
  });

  it('search filters by unit_number', async () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenants,
      isLoading: false,
      isError: false,
      error: null,
    });

    renderPage();

    const searchInput = screen.getByPlaceholderText(/search by block/i);
    fireEvent.change(searchInput, { target: { value: '10' } });

    await waitFor(
      () => {
        expect(screen.getByText(/Unit 10/i)).toBeInTheDocument();
        expect(screen.queryByText(/Unit 01/i)).not.toBeInTheDocument();
      },
      { timeout: 500 }
    );
  });

  it('"+ Add Tenant" button navigates to /tenants/new', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenants,
      isLoading: false,
      isError: false,
      error: null,
    });

    renderPage();
    fireEvent.click(screen.getByText('+ Add Tenant'));
    expect(mockNavigate).toHaveBeenCalledWith('/tenants/new');
  });
});
