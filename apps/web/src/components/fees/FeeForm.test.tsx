import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import FeeForm from './FeeForm';
import type { Fee } from '../../types/fee';

const mockFee: Fee = {
  id: 'fee-1',
  tenant_id: 'tenant-1',
  type: 'mandatory',
  amount: 25000,
  description: 'Security Fee',
  effective_date: '2026-07-01',
  paid_at: null,
  created_at: '2026-05-01T00:00:00Z',
};

function renderFeeForm(props: Record<string, unknown> = {}) {
  const defaultProps = {
    tenantId: 'tenant-1',
    monthlyFee: 50000,
    onSubmit: vi.fn().mockResolvedValue(undefined),
    onCancel: vi.fn(),
    isLoading: false,
  };

  return render(
    <MemoryRouter>
      <FeeForm {...defaultProps} {...props} />
    </MemoryRouter>
  );
}

describe('FeeForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders all form fields', () => {
    renderFeeForm();
    expect(screen.getByLabelText('Fee Type')).toBeInTheDocument();
    expect(screen.getByLabelText('Description')).toBeInTheDocument();
    expect(screen.getByLabelText('Amount (Rp)')).toBeInTheDocument();
    expect(screen.getByLabelText('Effective Date')).toBeInTheDocument();
    expect(screen.getByLabelText('Payment Date')).toBeInTheDocument();
  });

  it('type selector shows Mandatory Fee and Voluntary Contribution options', () => {
    renderFeeForm();
    expect(screen.getByText('Mandatory Fee')).toBeInTheDocument();
    expect(screen.getByText('Voluntary Contribution')).toBeInTheDocument();
  });

  it('shows validation error for empty description', async () => {
    renderFeeForm();
    fireEvent.change(screen.getByLabelText('Amount (Rp)'), {
      target: { value: '25000' },
    });
    fireEvent.change(screen.getByLabelText('Effective Date'), {
      target: { value: '2026-07-01' },
    });
    fireEvent.click(screen.getByRole('button', { name: /save fee/i }));
    await waitFor(() => {
      expect(screen.getByText('Description is required.')).toBeInTheDocument();
    });
  });

  it('shows validation error for empty amount', async () => {
    renderFeeForm();
    fireEvent.change(screen.getByLabelText('Description'), {
      target: { value: 'Security Fee' },
    });
    fireEvent.change(screen.getByLabelText('Effective Date'), {
      target: { value: '2026-07-01' },
    });
    fireEvent.click(screen.getByRole('button', { name: /save fee/i }));
    await waitFor(() => {
      expect(screen.getByText('Amount is required.')).toBeInTheDocument();
    });
  });

  it('shows validation error for amount > monthlyFee', async () => {
    renderFeeForm({ monthlyFee: 50000 });
    fireEvent.change(screen.getByLabelText('Description'), {
      target: { value: 'Security Fee' },
    });
    fireEvent.change(screen.getByLabelText('Amount (Rp)'), {
      target: { value: '60000' },
    });
    fireEvent.change(screen.getByLabelText('Effective Date'), {
      target: { value: '2026-07-01' },
    });
    fireEvent.click(screen.getByRole('button', { name: /save fee/i }));
    await waitFor(() => {
      expect(
        screen.getByText("Fee amount cannot exceed the tenant's monthly fee.")
      ).toBeInTheDocument();
    });
  });

  it('shows error for past effective_date', async () => {
    renderFeeForm();
    fireEvent.change(screen.getByLabelText('Description'), {
      target: { value: 'Security Fee' },
    });
    fireEvent.change(screen.getByLabelText('Amount (Rp)'), {
      target: { value: '25000' },
    });
    fireEvent.change(screen.getByLabelText('Effective Date'), {
      target: { value: '2000-01-01' },
    });
    fireEvent.click(screen.getByRole('button', { name: /save fee/i }));
    await waitFor(() => {
      expect(
        screen.getByText('Effective date cannot be in the past.')
      ).toBeInTheDocument();
    });
  });

  it('calls onSubmit with form data when validation passes', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);
    renderFeeForm({ onSubmit });
    fireEvent.change(screen.getByLabelText('Description'), {
      target: { value: 'Security Fee' },
    });
    fireEvent.change(screen.getByLabelText('Amount (Rp)'), {
      target: { value: '25000' },
    });
    fireEvent.change(screen.getByLabelText('Effective Date'), {
      target: { value: '2026-07-01' },
    });
    fireEvent.click(screen.getByRole('button', { name: /save fee/i }));
    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        type: 'mandatory',
        description: 'Security Fee',
        amount: 25000,
        effective_date: '2026-07-01',
      });
    });
  });

  it('pre-populates fields when initialData is provided (edit mode)', () => {
    renderFeeForm({ initialData: mockFee });
    const descriptionInput = screen.getByLabelText('Description') as HTMLInputElement;
    const amountInput = screen.getByLabelText('Amount (Rp)') as HTMLInputElement;
    expect(descriptionInput.value).toBe('Security Fee');
    expect(amountInput.value).toBe('25000');
  });

  it('shows error alert when onSubmit fails', async () => {
    const onSubmit = vi.fn().mockRejectedValue(new Error('Something went wrong'));
    renderFeeForm({ onSubmit });
    fireEvent.change(screen.getByLabelText('Description'), {
      target: { value: 'Security Fee' },
    });
    fireEvent.change(screen.getByLabelText('Amount (Rp)'), {
      target: { value: '25000' },
    });
    fireEvent.change(screen.getByLabelText('Effective Date'), {
      target: { value: '2026-07-01' },
    });
    fireEvent.click(screen.getByRole('button', { name: /save fee/i }));
    await waitFor(() => {
      expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    });
  });

  it('calls onCancel when Cancel button clicked', () => {
    const onCancel = vi.fn();
    renderFeeForm({ onCancel });
    fireEvent.click(screen.getByText('Cancel'));
    expect(onCancel).toHaveBeenCalled();
  });
});
