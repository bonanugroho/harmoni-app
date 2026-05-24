import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import TenantCreatePage from './TenantCreatePage';

const mockMutateAsync = vi.fn();

vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual('@tanstack/react-query');
  return {
    ...actual,
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

import { useMutation } from '@tanstack/react-query';

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false } },
});

function renderWithProviders() {
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <TenantCreatePage />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('TenantCreatePage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders TenantForm inside a container', () => {
    (useMutation as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
    });

    renderWithProviders();

    expect(screen.getByText('Add New Tenant')).toBeInTheDocument();
    expect(screen.getByLabelText(/block/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/unit number/i)).toBeInTheDocument();
  });

  it('submit calls mutation.mutateAsync with form data', async () => {
    (useMutation as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
    });

    renderWithProviders();

    // Fill all required fields
    fireEvent.change(screen.getByLabelText(/block/i), {
      target: { value: 'A' },
    });
    fireEvent.change(screen.getByLabelText(/unit number/i), {
      target: { value: '01' },
    });
    fireEvent.change(screen.getByLabelText(/monthly fee/i), {
      target: { value: '50000' },
    });

    // Fill fee fields
    const descriptionInputs = screen.getAllByLabelText(/description/i);
    fireEvent.change(descriptionInputs[0], {
      target: { value: 'Security Fee' },
    });

    const amountInputs = screen.getAllByLabelText(/amount/i);
    fireEvent.change(amountInputs[0], {
      target: { value: '25000' },
    });

    const dateInputs = screen.getAllByLabelText(/effective date/i);
    fireEvent.change(dateInputs[0], {
      target: { value: '2026-06-01' },
    });

    fireEvent.click(screen.getByRole('button', { name: /save tenant/i }));

    await waitFor(() => {
      expect(mockMutateAsync).toHaveBeenCalledTimes(1);
      const data = mockMutateAsync.mock.calls[0][0];
      expect(data.block).toBe('A');
      expect(data.unit_number).toBe('01');
      expect(data.mandatory_fees).toHaveLength(1);
    });
  });

  it('on success navigates to /tenants', async () => {
    mockMutateAsync.mockResolvedValue({ id: '1' });

    (useMutation as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
    });

    renderWithProviders();

    // Fill all required fields
    fireEvent.change(screen.getByLabelText(/block/i), {
      target: { value: 'A' },
    });
    fireEvent.change(screen.getByLabelText(/unit number/i), {
      target: { value: '01' },
    });
    fireEvent.change(screen.getByLabelText(/monthly fee/i), {
      target: { value: '50000' },
    });

    const descriptionInputs = screen.getAllByLabelText(/description/i);
    fireEvent.change(descriptionInputs[0], {
      target: { value: 'Security Fee' },
    });

    const amountInputs = screen.getAllByLabelText(/amount/i);
    fireEvent.change(amountInputs[0], {
      target: { value: '25000' },
    });

    const dateInputs = screen.getAllByLabelText(/effective date/i);
    fireEvent.change(dateInputs[0], {
      target: { value: '2026-06-01' },
    });

    fireEvent.click(screen.getByRole('button', { name: /save tenant/i }));

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/tenants');
    });
  });
});
