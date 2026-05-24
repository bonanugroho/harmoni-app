import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import EmptyState from './EmptyState';

describe('EmptyState', () => {
  it('renders heading and body text', () => {
    render(
      <EmptyState
        heading="No Tenants Yet"
        body="Start by adding your first tenant."
      />
    );
    expect(screen.getByText('No Tenants Yet')).toBeInTheDocument();
    expect(screen.getByText('Start by adding your first tenant.')).toBeInTheDocument();
  });

  it('renders CTA button when action provided', () => {
    render(
      <EmptyState
        heading="No Tenants Yet"
        body="Start by adding your first tenant."
        action={{ label: '+ Add Tenant', onClick: vi.fn() }}
      />
    );
    expect(screen.getByText('+ Add Tenant')).toBeInTheDocument();
  });

  it('does not render CTA when action not provided', () => {
    render(
      <EmptyState
        heading="No Tenants Yet"
        body="Start by adding your first tenant."
      />
    );
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('calls action.onClick when CTA clicked', () => {
    const onClick = vi.fn();
    render(
      <EmptyState
        heading="No Tenants Yet"
        body="Start by adding your first tenant."
        action={{ label: '+ Add Tenant', onClick }}
      />
    );
    fireEvent.click(screen.getByText('+ Add Tenant'));
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
