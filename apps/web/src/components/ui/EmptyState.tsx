interface EmptyStateProps {
  heading: string;
  body: string;
  action?: {
    label: string;
    onClick: () => void;
  };
}

export default function EmptyState({ heading, body, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-gray-200 py-12 px-4 text-center">
      <div className="mx-auto h-12 w-12 rounded-full bg-gray-200" aria-hidden="true" />
      <h3 className="mt-4 text-lg font-semibold text-gray-900">{heading}</h3>
      <p className="mt-2 text-sm text-gray-500">{body}</p>
      {action && (
        <button
          onClick={action.onClick}
          className="mt-4 min-h-[44px] rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          {action.label}
        </button>
      )}
    </div>
  );
}
