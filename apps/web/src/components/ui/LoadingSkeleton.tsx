interface LoadingSkeletonProps {
  variant?: 'card' | 'list' | 'form';
  count?: number;
}

export default function LoadingSkeleton({
  variant = 'card',
  count = 3,
}: LoadingSkeletonProps) {
  if (variant === 'card') {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" data-testid="loading-skeleton">
        {Array.from({ length: count }).map((_, i) => (
          <div
            key={i}
            className="animate-pulse rounded-lg border border-gray-200 p-4"
            data-testid="loading-skeleton"
          >
            <div className="h-4 w-3/4 rounded bg-gray-200" />
            <div className="mt-3 h-3 w-1/2 rounded bg-gray-200" />
            <div className="mt-4 h-3 w-full rounded bg-gray-200" />
          </div>
        ))}
      </div>
    );
  }

  if (variant === 'list') {
    return (
      <div className="space-y-3" data-testid="loading-skeleton">
        {Array.from({ length: count }).map((_, i) => (
          <div
            key={i}
            className="animate-pulse rounded-lg border border-gray-200 p-4"
            data-testid="loading-skeleton"
          >
            <div className="h-4 w-1/2 rounded bg-gray-200" />
            <div className="mt-2 h-3 w-3/4 rounded bg-gray-200" />
          </div>
        ))}
      </div>
    );
  }

  if (variant === 'form') {
    return (
      <div className="space-y-4" data-testid="loading-skeleton">
        {Array.from({ length: count }).map((_, i) => (
          <div key={i} className="space-y-2" data-testid="loading-skeleton">
            <div className="h-3 w-1/4 rounded bg-gray-200" />
            <div className="h-12 w-full rounded-md bg-gray-200" />
          </div>
        ))}
      </div>
    );
  }

  return null;
}
