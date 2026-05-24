import { type InputHTMLAttributes } from 'react';

interface DatePickerProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  label: string;
  error?: string;
}

export default function DatePicker({ label, error, id, className = '', ...props }: DatePickerProps) {
  const inputId = id || props.name;
  const errorId = `${inputId}-error`;

  return (
    <div className="space-y-2">
      <label htmlFor={inputId} className="block text-base font-medium text-gray-700">
        {label}
      </label>
      <input
        id={inputId}
        type="date"
        className={`min-h-[44px] w-full rounded-md border px-3 py-2 text-base outline-none transition-colors ${
          error
            ? 'border-red-600 focus:border-red-600'
            : 'border-gray-200 focus:border-blue-600'
        } ${className}`}
        aria-invalid={!!error}
        aria-describedby={error ? errorId : undefined}
        {...props}
      />
      {error && (
        <p id={errorId} className="text-sm text-red-600" role="alert">
          {error}
        </p>
      )}
    </div>
  );
}
