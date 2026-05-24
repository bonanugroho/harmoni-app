import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import StatusBadge from './StatusBadge';

describe('StatusBadge', () => {
  it('renders "Occupied" for occupied status', () => {
    render(<StatusBadge status="occupied" />);
    expect(screen.getByText('Occupied')).toBeInTheDocument();
  });

  it('renders "Vacant" for vacant status', () => {
    render(<StatusBadge status="vacant" />);
    expect(screen.getByText('Vacant')).toBeInTheDocument();
  });

  it('renders "Paid" for paid status', () => {
    render(<StatusBadge status="paid" />);
    expect(screen.getByText('Paid')).toBeInTheDocument();
  });

  it('renders "Unpaid" for unpaid status', () => {
    render(<StatusBadge status="unpaid" />);
    expect(screen.getByText('Unpaid')).toBeInTheDocument();
  });
});
