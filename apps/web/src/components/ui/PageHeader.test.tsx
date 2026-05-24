import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import PageHeader from './PageHeader';

describe('PageHeader', () => {
  it('renders title', () => {
    render(<PageHeader title="Tenants" />);
    expect(screen.getByText('Tenants')).toBeInTheDocument();
  });

  it('renders action button when action prop provided', () => {
    render(
      <PageHeader
        title="Tenants"
        action={{ label: '+ Add Tenant', onClick: vi.fn() }}
      />
    );
    expect(screen.getByText('+ Add Tenant')).toBeInTheDocument();
  });

  it('does not render action button when action prop not provided', () => {
    render(<PageHeader title="Tenants" />);
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('calls action.onClick when button clicked', () => {
    const onClick = vi.fn();
    render(
      <PageHeader
        title="Tenants"
        action={{ label: '+ Add Tenant', onClick }}
      />
    );
    fireEvent.click(screen.getByText('+ Add Tenant'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
