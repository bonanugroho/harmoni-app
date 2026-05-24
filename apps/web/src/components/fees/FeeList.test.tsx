import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import FeeList from './FeeList';
import type { Fee } from '../../types/fee';

const mockMandatoryFee: Fee = {
  id: 'fee-1',
  tenant_id: 'tenant-1',
  type: 'mandatory',
  amount: 25000,
  description: 'Security Fee',
  effective_date: '2026-06-01',
  paid_at: '2026-06-02T00:00:00Z',
  created_at: '2026-05-01T00:00:00Z',
};

const mockVoluntaryFee: Fee = {
  id: 'fee-2',
  tenant_id: 'tenant-1',
  type: 'voluntary',
  amount: 50000,
  description: 'Holiday Bonus',
  effective_date: '2026-12-01',
  paid_at: null,
  created_at: '2026-05-01T00:00:00Z',
};

function renderFeeList(props: Record<string, unknown> = {}) {
  return render(
    <MemoryRouter>
      <FeeList
        mandatoryFees={[]}
        voluntaryFees={[]}
        {...props}
      />
    </MemoryRouter>
  );
}

describe('FeeList', () => {
  it('renders mandatory fees section heading and fee cards', () => {
    renderFeeList({ mandatoryFees: [mockMandatoryFee] });
    expect(screen.getByText('Mandatory Fees')).toBeInTheDocument();
    expect(screen.getByText('Security Fee')).toBeInTheDocument();
    expect(screen.getByText(/Rp/)).toBeInTheDocument();
  });

  it('renders voluntary fees section heading and fee cards', () => {
    renderFeeList({ voluntaryFees: [mockVoluntaryFee] });
    expect(screen.getByText('Voluntary Contributions')).toBeInTheDocument();
    expect(screen.getByText('Holiday Bonus')).toBeInTheDocument();
  });

  it('shows EmptyState for mandatory fees when none provided', () => {
    renderFeeList({ mandatoryFees: [], voluntaryFees: [mockVoluntaryFee] });
    expect(screen.getByText('No Mandatory Fees Set')).toBeInTheDocument();
    expect(
      screen.getByText('Every tenant needs at least one mandatory fee. Add one now.')
    ).toBeInTheDocument();
  });

  it('shows EmptyState for voluntary fees when none provided', () => {
    renderFeeList({ voluntaryFees: [], mandatoryFees: [mockMandatoryFee] });
    expect(screen.getByText('No Voluntary Contributions Yet')).toBeInTheDocument();
    expect(
      screen.getByText('Residents can contribute voluntarily here.')
    ).toBeInTheDocument();
  });

  it('calls onDeleteFee when delete icon clicked', () => {
    const onDelete = vi.fn();
    renderFeeList({
      mandatoryFees: [mockMandatoryFee],
      onDeleteFee: onDelete,
    });
    const deleteButton = screen.getByLabelText('Delete fee');
    fireEvent.click(deleteButton);
    expect(onDelete).toHaveBeenCalledWith('fee-1');
  });

  it('calls onEditFee when edit icon clicked', () => {
    const onEdit = vi.fn();
    renderFeeList({
      mandatoryFees: [mockMandatoryFee],
      onEditFee: onEdit,
    });
    const editButton = screen.getByLabelText('Edit fee');
    fireEvent.click(editButton);
    expect(onEdit).toHaveBeenCalledWith('fee-1');
  });

  it('shows LoadingSkeleton when isLoading is true', () => {
    renderFeeList({ isLoading: true });
    // LoadingSkeleton variant="list" renders data-testid="loading-skeleton"
    const skeletons = screen.getAllByTestId('loading-skeleton');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders StatusBadge with "Paid" when paid_at is set', () => {
    renderFeeList({ mandatoryFees: [mockMandatoryFee] });
    expect(screen.getByText('Paid')).toBeInTheDocument();
  });

  it('renders StatusBadge with "Unpaid" when paid_at is null', () => {
    renderFeeList({ voluntaryFees: [mockVoluntaryFee] });
    expect(screen.getByText('Unpaid')).toBeInTheDocument();
  });
});
