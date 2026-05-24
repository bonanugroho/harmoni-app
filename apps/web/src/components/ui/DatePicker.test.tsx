import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import DatePicker from './DatePicker';

describe('DatePicker', () => {
  it('renders label and date input', () => {
    render(<DatePicker label="Effective Date" name="effectiveDate" />);
    expect(screen.getByText('Effective Date')).toBeInTheDocument();
    expect(screen.getByLabelText('Effective Date')).toBeInTheDocument();
  });

  it('shows error message when error prop is provided', () => {
    render(<DatePicker label="Effective Date" name="effectiveDate" error="Date is required" />);
    expect(screen.getByText('Date is required')).toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveTextContent('Date is required');
  });

  it('input type is date', () => {
    render(<DatePicker label="Effective Date" name="effectiveDate" />);
    expect(screen.getByLabelText('Effective Date')).toHaveAttribute('type', 'date');
  });
});
