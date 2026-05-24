import { type SelectHTMLAttributes } from 'react';

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label: string;
  error?: string;
  options: Array<{ value: string; label: string }>;
  placeholder?: string;
}

export default function Select({ label, error, id, options, placeholder, className = '', ...props }: SelectProps) {
  const inputId = id || props.name;
  const errorId = `${inputId}-error`;

  return (
    <div className="space-y-2">
      <label htmlFor={inputId} className="block text-base font-medium text-gray-700">
        {label}
      </label>
      <select
        id={inputId}
        className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
          error
            ? 'border-red-600 focus:border-red-600'
            : 'border-gray-200 focus:border-blue-600'
        } ${className}`}
        aria-invalid={!!error}
        aria-describedby={error ? errorId : undefined}
        {...props}
      >
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
      {error && (
        <p id={errorId} className="text-sm text-red-600" role="alert">
          {error}
        </p>
      )}
    </div>
  );
}
