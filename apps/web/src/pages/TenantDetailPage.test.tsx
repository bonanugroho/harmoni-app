import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import TenantDetailPage from './TenantDetailPage';

vi.mock('../hooks/useTenant');
vi.mock('../hooks/useFees');
vi.mock('../hooks/useCreateFee');
vi.mock('../hooks/useUpdateFee');
vi.mock('../hooks/useDeleteFee');

import { useTenant } from '../hooks/useTenant';
import { useFees } from '../hooks/useFees';
import { useCreateFee } from '../hooks/useCreateFee';
import { useUpdateFee } from '../hooks/useUpdateFee';
import { useDeleteFee } from '../hooks/useDeleteFee';

const mockTenant = {
  id: '1',
  block: 'A',
  unit_number: '01',
  occupancy: 'occupied' as const,
  monthly_fee: 50000,
  territory_id: 't1',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

const mockFees = {
  mandatory_fees: [
    {
      id: 'fee-1',
      tenant_id: '1',
      type: 'mandatory',
      amount: 25000,
      description: 'Security Fee',
      effective_date: '2026-06-01',
      paid_at: '2026-06-02T00:00:00Z',
      created_at: '2026-05-01T00:00:00Z',
    },
  ],
  voluntary_fees: [
    {
      id: 'fee-2',
      tenant_id: '1',
      type: 'voluntary',
      amount: 50000,
      description: 'Holiday Bonus',
      effective_date: '2026-12-01',
      paid_at: null,
      created_at: '2026-05-01T00:00:00Z',
    },
  ],
};

const mockMutateAsync = vi.fn().mockResolvedValue(undefined);

function renderPage() {
  return render(
    <MemoryRouter initialEntries={['/tenants/1']}>
      <TenantDetailPage />
    </MemoryRouter>
  );
}

describe('TenantDetailPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();

    (useTenant as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockTenant,
      isLoading: false,
      isError: false,
    });

    (useFees as ReturnType<typeof vi.fn>).mockReturnValue({
      data: mockFees,
      isLoading: false,
      isError: false,
    });

    (useCreateFee as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
    });

    (useUpdateFee as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
    });

    (useDeleteFee as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
    });
  });

  it('shows LoadingSkeleton while tenant is loading', () => {
    (useTenant as ReturnType<typeof vi.fn>).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
    });

    renderPage();
    const skeletons = screen.getAllByTestId('loading-skeleton');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders tenant info header with block, unit_number, occupancy badge, monthly_fee', () => {
    renderPage();
    expect(screen.getByText(/Block A/)).toBeInTheDocument();
    expect(screen.getByText(/Unit 01/)).toBeInTheDocument();
    expect(screen.getByText('Occupied')).toBeInTheDocument();
    // Rp appears in header and fee cards, so use getAllByText
    expect(screen.getAllByText(/Rp/).length).toBeGreaterThan(0);
  });

  it('renders FeeList with mandatory and voluntary fees', () => {
    renderPage();
    expect(screen.getByText('Mandatory Fees')).toBeInTheDocument();
    expect(screen.getByText('Voluntary Contributions')).toBeInTheDocument();
    expect(screen.getByText('Security Fee')).toBeInTheDocument();
    expect(screen.getByText('Holiday Bonus')).toBeInTheDocument();
  });

  it('shows error when tenant fetch fails', () => {
    (useTenant as ReturnType<typeof vi.fn>).mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
    });

    renderPage();
    expect(screen.getByText('Failed to load tenant. Please try again.')).toBeInTheDocument();
  });

  it('opens FeeForm modal when Record Fee button clicked', () => {
    renderPage();
    fireEvent.click(screen.getByText('Record Fee'));
    // Modal renders FeeForm with form fields inside
    expect(screen.getByLabelText('Description')).toBeInTheDocument();
    expect(screen.getByLabelText('Fee Type')).toBeInTheDocument();
  });

  it('creates fee on FeeForm submit', async () => {
    const createMutateAsync = vi.fn().mockResolvedValue(undefined);
    (useCreateFee as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: createMutateAsync,
      isPending: false,
    });

    renderPage();
    fireEvent.click(screen.getByText('Record Fee'));

    // Fill FeeForm fields
    const descInput = screen.getByLabelText('Description');
    const amountInput = screen.getByLabelText('Amount (Rp)');
    const dateInput = screen.getByLabelText('Effective Date');

    fireEvent.change(descInput, { target: { value: 'New Fee' } });
    fireEvent.change(amountInput, { target: { value: '25000' } });
    fireEvent.change(dateInput, { target: { value: '2026-07-01' } });

    fireEvent.click(screen.getByRole('button', { name: /save fee/i }));

    await waitFor(() => {
      expect(createMutateAsync).toHaveBeenCalled();
    });
  });

  it('opens ConfirmDialog when delete icon clicked', () => {
    renderPage();
    // Click delete icon on first fee card
    const deleteButtons = screen.getAllByLabelText('Delete fee');
    fireEvent.click(deleteButtons[0]);

    expect(screen.getByText('Delete this fee?')).toBeInTheDocument();
    expect(screen.getByText('This cannot be undone.')).toBeInTheDocument();
  });

  it('deletes fee on ConfirmDialog confirm', async () => {
    const deleteMutateAsync = vi.fn().mockResolvedValue(undefined);
    (useDeleteFee as ReturnType<typeof vi.fn>).mockReturnValue({
      mutateAsync: deleteMutateAsync,
      isPending: false,
    });

    renderPage();
    // Click delete icon
    const deleteButtons = screen.getAllByLabelText('Delete fee');
    fireEvent.click(deleteButtons[0]);
    // Confirm delete
    fireEvent.click(screen.getByText('Delete'));
    await waitFor(() => {
      expect(deleteMutateAsync).toHaveBeenCalled();
    });
  });
});
