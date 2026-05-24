import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Input from './Input';

describe('Input', () => {
  it('renders label text', () => {
    render(<Input label="Full Name" name="fullName" />);
    expect(screen.getByText('Full Name')).toBeInTheDocument();
  });

  it('renders input with correct id matching htmlFor', () => {
    render(<Input label="Email" id="email" />);
    const input = screen.getByLabelText('Email');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('id', 'email');
  });

  it('shows error message when error prop is provided', () => {
    render(<Input label="Email" name="email" error="Email is required" />);
    expect(screen.getByText('Email is required')).toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveTextContent('Email is required');
  });

  it('sets aria-invalid when error is present', () => {
    render(<Input label="Email" name="email" error="Invalid email" />);
    expect(screen.getByLabelText('Email')).toHaveAttribute('aria-invalid', 'true');
  });

  it('does not show error when error prop is undefined', () => {
    render(<Input label="Email" name="email" />);
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('applies 44px minimum height', () => {
    render(<Input label="Email" name="email" />);
    const input = screen.getByLabelText('Email');
    expect(input.className).toContain('min-h-[44px]');
  });
});
