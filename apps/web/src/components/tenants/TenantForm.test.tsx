import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import TenantForm from './TenantForm';
import type { Tenant, CreateTenantRequest } from '../../types/tenant';

const mockSubmit = vi.fn();

function renderForm(overrides?: Partial<React.ComponentProps<typeof TenantForm>>) {
  return render(
    <MemoryRouter>
      <TenantForm onSubmit={mockSubmit} {...overrides} />
    </MemoryRouter>
  );
}

const mockTenant: Tenant = {
  id: '1',
  block: 'A',
  unit_number: '01',
  occupancy: 'occupied',
  monthly_fee: 50000,
  territory_id: 'territory-1',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

describe('TenantForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders all form fields', () => {
    renderForm();
    expect(screen.getByLabelText(/block/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/unit number/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/occupancy status/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/monthly fee/i)).toBeInTheDocument();
  });

  it('shows mandatory fee section with initial empty fee row', () => {
    renderForm();
    expect(screen.getByText('Mandatory Fees')).toBeInTheDocument();
    expect(screen.getByText('Fee 1')).toBeInTheDocument();
  });

  it('adds fee row when "+ Add Another Fee" is clicked', () => {
    renderForm();
    fireEvent.click(screen.getByText('Add Another Fee'));
    expect(screen.getByText('Fee 2')).toBeInTheDocument();
  });

  it('removes fee row when X button is clicked', () => {
    renderForm();
    // Add a second fee row first
    fireEvent.click(screen.getByText('Add Another Fee'));
    expect(screen.getByText('Fee 2')).toBeInTheDocument();

    // Remove the first fee row
    const removeButtons = screen.getAllByRole('button', { name: /remove fee/i });
    fireEvent.click(removeButtons[0]);

    // The second fee should now be labeled Fee 1
    expect(screen.getByText('Fee 1')).toBeInTheDocument();
    expect(screen.queryByText('Fee 2')).not.toBeInTheDocument();
  });

  it('shows validation error for empty block on submit', async () => {
    renderForm();
    fireEvent.click(screen.getByRole('button', { name: /save tenant/i }));

    await waitFor(() => {
      expect(screen.getByText('Block is required.')).toBeInTheDocument();
    });
  });

  it('shows validation error when total mandatory fees exceed monthly fee cap', async () => {
    renderForm();

    fireEvent.change(screen.getByLabelText(/block/i), {
      target: { value: 'A' },
    });
    fireEvent.change(screen.getByLabelText(/unit number/i), {
      target: { value: '01' },
    });
    fireEvent.change(screen.getByLabelText(/monthly fee/i), {
      target: { value: '50000' },
    });

    // Add a second mandatory fee
    fireEvent.click(screen.getByText('Add Another Fee'));

    // Fill both fees with amounts that sum > 50000
    const amountInputs = screen.getAllByLabelText(/amount/i);
    fireEvent.change(amountInputs[0], { target: { value: '30000' } });
    fireEvent.change(amountInputs[1], { target: { value: '30000' } });

    // Fill descriptions and dates
    const descriptionInputs = screen.getAllByLabelText(/description/i);
    fireEvent.change(descriptionInputs[0], { target: { value: 'Security' } });
    fireEvent.change(descriptionInputs[1], { target: { value: 'Maintenance' } });

    const dateInputs = screen.getAllByLabelText(/effective date/i);
    fireEvent.change(dateInputs[0], { target: { value: '2026-06-01' } });
    fireEvent.change(dateInputs[1], { target: { value: '2026-06-01' } });

    fireEvent.click(screen.getByRole('button', { name: /save tenant/i }));

    await waitFor(() => {
      expect(
        screen.getByText('Total mandatory fees cannot exceed the monthly fee.')
      ).toBeInTheDocument();
    });
  });

  it('shows validation error for empty fee description on submit', async () => {
    renderForm();

    // Fill required fields but leave fee description empty
    fireEvent.change(screen.getByLabelText(/block/i), {
      target: { value: 'A' },
    });
    fireEvent.change(screen.getByLabelText(/unit number/i), {
      target: { value: '01' },
    });
    fireEvent.change(screen.getByLabelText(/monthly fee/i), {
      target: { value: '50000' },
    });

    // Our fee entry has empty description by default

    fireEvent.click(screen.getByRole('button', { name: /save tenant/i }));

    await waitFor(() => {
      expect(screen.getByText('Description is required.')).toBeInTheDocument();
    });
  });

  it('calls onSubmit with form data when validation passes', async () => {
    mockSubmit.mockResolvedValue(undefined);

    renderForm();

    fireEvent.change(screen.getByLabelText(/block/i), {
      target: { value: 'A' },
    });
    fireEvent.change(screen.getByLabelText(/unit number/i), {
      target: { value: '01' },
    });
    fireEvent.change(screen.getByLabelText(/monthly fee/i), {
      target: { value: '50000' },
    });

    // Fill fee description and effective date
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
      expect(mockSubmit).toHaveBeenCalledTimes(1);
      const data: CreateTenantRequest = mockSubmit.mock.calls[0][0];
      expect(data.block).toBe('A');
      expect(data.unit_number).toBe('01');
      expect(data.occupancy).toBe('occupied');
      expect(data.monthly_fee).toBe(50000);
      expect(data.mandatory_fees).toHaveLength(1);
      expect(data.mandatory_fees[0].description).toBe('Security Fee');
      expect(data.mandatory_fees[0].amount).toBe(25000);
      expect(data.mandatory_fees[0].effective_date).toBe('2026-06-01');
    });
  });

  it('shows submitError alert when provided', async () => {
    mockSubmit.mockRejectedValue(new Error('Server error occurred'));

    renderForm();

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
      expect(screen.getByText('Server error occurred')).toBeInTheDocument();
    });
  });

  it('pre-populates fields when initialData is provided', () => {
    renderForm({ initialData: mockTenant });

    const blockInput = screen.getByLabelText(/block/i) as HTMLInputElement;
    expect(blockInput.value).toBe('A');

    const unitInput = screen.getByLabelText(/unit number/i) as HTMLInputElement;
    expect(unitInput.value).toBe('01');

    const feeInput = screen.getByLabelText(/monthly fee/i) as HTMLInputElement;
    expect(feeInput.value).toBe('50000');

    expect(screen.getByText('Edit Unit A-01')).toBeInTheDocument();
  });

  it('shows spinner and "Saving..." when isLoading is true', () => {
    renderForm({ isLoading: true });
    const submitButton = screen.getByRole('button', { name: /saving/i });
    expect(submitButton).toBeDisabled();
  });
});
