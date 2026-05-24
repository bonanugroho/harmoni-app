import { useState, useEffect } from 'react';
import Select from '../ui/Select';
import Input from '../ui/Input';
import DatePicker from '../ui/DatePicker';
import type { Fee, CreateFeeRequest } from '../../types/fee';

interface FeeFormProps {
  tenantId: string;
  monthlyFee: number;
  initialData?: Fee;
  onSubmit: (data: CreateFeeRequest) => Promise<void>;
  isLoading?: boolean;
  onCancel: () => void;
}

export default function FeeForm({
  tenantId: _tenantId,
  monthlyFee,
  initialData,
  onSubmit,
  isLoading = false,
  onCancel,
}: FeeFormProps) {
  const [feeType, setFeeType] = useState<'mandatory' | 'voluntary'>('mandatory');
  const [description, setDescription] = useState('');
  const [amount, setAmount] = useState('');
  const [effectiveDate, setEffectiveDate] = useState('');
  const [paidAt, setPaidAt] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitError, setSubmitError] = useState('');

  useEffect(() => {
    if (initialData) {
      setFeeType(initialData.type || 'mandatory');
      setDescription(initialData.description);
      setAmount(String(initialData.amount));
      setEffectiveDate(initialData.effective_date);
      setPaidAt(initialData.paid_at || '');
    }
  }, [initialData]);

  function validate(): Record<string, string> {
    const errs: Record<string, string> = {};

    if (!description.trim()) {
      errs.description = 'Description is required.';
    }

    if (!amount.trim() || isNaN(parseFloat(amount))) {
      errs.amount = 'Amount is required.';
    } else if (parseFloat(amount) <= 0) {
      errs.amount = 'Amount must be a positive number.';
    } else if (parseFloat(amount) > monthlyFee) {
      errs.amount = "Fee amount cannot exceed the tenant's monthly fee.";
    }

    if (!effectiveDate.trim()) {
      errs.effective_date = 'Effective date is required.';
    } else {
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      const effDate = new Date(effectiveDate);
      if (effDate < today) {
        errs.effective_date = 'Effective date cannot be in the past.';
      }
    }

    if (paidAt.trim() && effectiveDate.trim()) {
      const paidDate = new Date(paidAt);
      const effDate = new Date(effectiveDate);
      if (paidDate < effDate) {
        errs.paid_at = 'Payment date must be after the effective date.';
      }
    }

    return errs;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitError('');

    const validationErrors = validate();
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length > 0) {
      return;
    }

    try {
      const data: CreateFeeRequest = {
        type: feeType,
        description: description.trim(),
        amount: parseFloat(amount),
        effective_date: effectiveDate,
      };
      if (paidAt.trim()) {
        data.paid_at = paidAt;
      }
      await onSubmit(data);
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'An unexpected error occurred');
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {submitError && (
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          {submitError}
        </div>
      )}

      <Select
        name="fee_type"
        label="Fee Type"
        options={[
          { value: 'mandatory', label: 'Mandatory Fee' },
          { value: 'voluntary', label: 'Voluntary Contribution' },
        ]}
        value={feeType}
        onChange={(e) => setFeeType(e.target.value as 'mandatory' | 'voluntary')}
        error={errors.fee_type}
      />

      <Input
        name="description"
        label="Description"
        placeholder="Security Fee"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        error={errors.description}
      />

      <Input
        name="amount"
        label="Amount (Rp)"
        type="number"
        placeholder="25000"
        value={amount}
        onChange={(e) => setAmount(e.target.value)}
        error={errors.amount}
      />

      <DatePicker
        name="effective_date"
        label="Effective Date"
        value={effectiveDate}
        onChange={(e) => setEffectiveDate(e.target.value)}
        error={errors.effective_date}
      />

      <DatePicker
        name="paid_at"
        label="Payment Date"
        value={paidAt}
        onChange={(e) => setPaidAt(e.target.value)}
        error={errors.paid_at}
      />
      <p className="-mt-2 text-xs text-gray-500">(Optional)</p>

      <div className="flex justify-end gap-3 pt-2">
        <button
          type="button"
          onClick={onCancel}
          className="min-h-[44px] rounded-md border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isLoading}
          className="min-h-[44px] rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {isLoading ? 'Saving...' : 'Save Fee'}
        </button>
      </div>
    </form>
  );
}
