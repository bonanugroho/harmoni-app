import { describe, it, expect } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import Select from './Select';

const options = [
  { value: 'occupied', label: 'Occupied' },
  { value: 'vacant', label: 'Vacant' },
];

describe('Select', () => {
  it('renders label and all options', () => {
    render(<Select label="Status" name="status" options={options} />);
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText('Occupied')).toBeInTheDocument();
    expect(screen.getByText('Vacant')).toBeInTheDocument();
  });

  it('shows error message when error prop is provided', () => {
    render(<Select label="Status" name="status" options={options} error="Status is required" />);
    expect(screen.getByText('Status is required')).toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveTextContent('Status is required');
  });

  it('renders placeholder option when provided', () => {
    render(<Select label="Status" name="status" options={options} placeholder="Select status" />);
    expect(screen.getByText('Select status')).toBeInTheDocument();
    expect(screen.getByText('Select status')).toBeDisabled();
  });

  it('selects correct option on change', () => {
    render(<Select label="Status" name="status" options={options} />);
    const select = screen.getByLabelText('Status') as HTMLSelectElement;
    fireEvent.change(select, { target: { value: 'vacant' } });
    expect(select.value).toBe('vacant');
  });
});
