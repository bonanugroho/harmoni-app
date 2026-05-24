import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import TenantEditPage from './TenantEditPage';

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
    useParams: () => ({ id: '1' }),
  };
});

import { useQuery, useMutation } from '@tanstack/react-query';

const mockTenant = {
  id: '1',
  block: 'A',
  unit_number: '01',
  occupancy: 'occupied' as const,
  monthly_fee: 50000,
  territory_id: 'territory-1',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

function renderWithProviders() {
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <TenantEditPage />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('TenantEditPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows LoadingSkeleton when tenant is loading', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
    });

    renderWithProviders();

    const skeletons = screen.getAllByTestId('loading-skeleton');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders TenantForm with initialData when tenant loads', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenant,
      isLoading: false,
      isError: false,
    });

    (useMutation as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: vi.fn().mockResolvedValue({}),
      isPending: false,
    });

    renderWithProviders();

    expect(screen.getByText('Edit Unit A-01')).toBeInTheDocument();
    const blockInput = screen.getByLabelText(/block/i) as HTMLInputElement;
    expect(blockInput.value).toBe('A');
  });

  it('shows error when tenant fetch fails', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
    });

    renderWithProviders();

    expect(
      screen.getByText(/failed to load tenant/i)
    ).toBeInTheDocument();
  });

  it('shows delete danger button', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenant,
      isLoading: false,
      isError: false,
    });

    (useMutation as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: vi.fn().mockResolvedValue({}),
      isPending: false,
    });

    renderWithProviders();

    expect(
      screen.getByRole('button', { name: /delete tenant/i })
    ).toBeInTheDocument();
  });

  it('opens ConfirmDialog when delete button clicked', () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenant,
      isLoading: false,
      isError: false,
    });

    (useMutation as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: vi.fn().mockResolvedValue({}),
      isPending: false,
    });

    renderWithProviders();

    fireEvent.click(screen.getByRole('button', { name: /delete tenant/i }));
    expect(screen.getByText('Delete Unit A-01?')).toBeInTheDocument();
    expect(screen.getByText('This cannot be undone.')).toBeInTheDocument();
  });

  it('navigates to /tenants on successful delete', async () => {
    (useQuery as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenant,
      isLoading: false,
      isError: false,
    });

    const deleteMutateAsync = vi.fn().mockResolvedValue(undefined);
    (useMutation as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: deleteMutateAsync,
      isPending: false,
    });

    renderWithProviders();

    // Open the delete dialog
    fireEvent.click(screen.getByRole('button', { name: /delete tenant/i }));

    // Click confirm delete in the dialog
    fireEvent.click(screen.getByRole('button', { name: /^delete$/i }));

    await waitFor(() => {
      expect(deleteMutateAsync).toHaveBeenCalled();
      expect(mockNavigate).toHaveBeenCalledWith('/tenants');
    });
  });
});
