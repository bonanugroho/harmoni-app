import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import TenantCard from './TenantCard';
import type { Tenant } from '../../types/tenant';

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

function renderCard(overrides?: Partial<React.ComponentProps<typeof TenantCard>>) {
  return render(
    <MemoryRouter>
      <TenantCard tenant={mockTenant} {...overrides} />
    </MemoryRouter>
  );
}

describe('TenantCard', () => {
  it('renders block and unit_number', () => {
    renderCard();
    expect(screen.getByText(/Block A/i)).toBeInTheDocument();
    expect(screen.getByText(/Unit 01/i)).toBeInTheDocument();
  });

  it('renders occupancy StatusBadge', () => {
    renderCard();
    expect(screen.getByText('Occupied')).toBeInTheDocument();
  });

  it('renders formatted monthly_fee', () => {
    renderCard();
    expect(screen.getByText(/50.000/)).toBeInTheDocument();
  });

  it('renders fee summary text when no fees', () => {
    renderCard();
    expect(screen.getByText('No fees configured')).toBeInTheDocument();
  });

  it('renders fee summary text with fee counts', () => {
    renderCard({ mandatoryFeeCount: 3, voluntaryFeeCount: 2 });
    expect(screen.getByText('3 mandatory fees · 2 contributions')).toBeInTheDocument();
  });

  it('calls onClick when card is clicked', () => {
    const onClick = vi.fn();
    renderCard({ onClick });
    fireEvent.click(screen.getByRole('button'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('calls onClick when Enter key is pressed', () => {
    const onClick = vi.fn();
    renderCard({ onClick });
    fireEvent.keyDown(screen.getByRole('button'), { key: 'Enter' });
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
