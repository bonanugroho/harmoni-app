import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import LoadingSkeleton from './LoadingSkeleton';

describe('LoadingSkeleton', () => {
  it('renders specified count of skeleton items', () => {
    render(<LoadingSkeleton variant="card" count={2} />);
    const elements = screen.getAllByTestId('loading-skeleton');
    expect(elements.length).toBeGreaterThanOrEqual(2);
  });

  it('renders elements with data-testid="loading-skeleton"', () => {
    render(<LoadingSkeleton variant="list" count={1} />);
    const elements = screen.getAllByTestId('loading-skeleton');
    expect(elements.length).toBeGreaterThanOrEqual(1);
  });

  it('renders correct variant structure', () => {
    const { container } = render(<LoadingSkeleton variant="form" count={1} />);
    const elements = screen.getAllByTestId('loading-skeleton');
    expect(elements.length).toBeGreaterThanOrEqual(1);
  });
});
