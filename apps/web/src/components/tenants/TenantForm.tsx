import { useState, useEffect, type FormEvent, type ChangeEvent } from 'react';
import { Link } from 'react-router-dom';
import { X, Plus } from 'lucide-react';
import Input from '../ui/Input';
import Select from '../ui/Select';
import type { Tenant, CreateTenantRequest } from '../../types/tenant';

interface MandatoryFeeEntry {
  description: string;
  amount: number;
  effective_date: string;
}

interface TenantFormProps {
  initialData?: Tenant;
  onSubmit: (data: CreateTenantRequest) => Promise<void>;
  isLoading?: boolean;
}

interface FormErrors {
  block?: string;
  unit_number?: string;
  monthly_fee?: string;
  mandatoryFees?: string;
  feeDescription?: Record<number, string>;
  feeAmount?: Record<number, string>;
  feeEffectiveDate?: Record<number, string>;
}

export default function TenantForm({
  initialData,
  onSubmit,
  isLoading = false,
}: TenantFormProps) {
  const [block, setBlock] = useState('');
  const [unitNumber, setUnitNumber] = useState('');
  const [occupancy, setOccupancy] = useState<'occupied' | 'vacant'>('occupied');
  const [monthlyFee, setMonthlyFee] = useState('');
  const [mandatoryFees, setMandatoryFees] = useState<MandatoryFeeEntry[]>([
    { description: '', amount: 0, effective_date: '' },
  ]);
  const [errors, setErrors] = useState<FormErrors>({});
  const [submitError, setSubmitError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (initialData) {
      setBlock(initialData.block);
      setUnitNumber(initialData.unit_number);
      setOccupancy(initialData.occupancy);
      setMonthlyFee(String(initialData.monthly_fee));
    }
  }, [initialData]);

  function validate(): FormErrors {
    const newErrors: FormErrors = {};

    if (!block.trim()) {
      newErrors.block = 'Block is required.';
    }

    if (!unitNumber.trim()) {
      newErrors.unit_number = 'Unit number is required.';
    }

    const fee = Number(monthlyFee);
    if (!monthlyFee || fee <= 0) {
      newErrors.monthly_fee = 'Monthly fee must be a positive amount.';
    }

    if (!initialData) {
      if (mandatoryFees.length === 0) {
        newErrors.mandatoryFees = 'At least one mandatory fee is required.';
      }

      const feeDescriptionErrors: Record<number, string> = {};
      const feeAmountErrors: Record<number, string> = {};
      const feeEffectiveDateErrors: Record<number, string> = {};

      const totalMandatoryFees = mandatoryFees.reduce((sum, f) => sum + Number(f.amount), 0);
      if (fee > 0 && totalMandatoryFees > fee) {
        newErrors.mandatoryFees = 'Total mandatory fees cannot exceed the monthly fee.';
      }

      mandatoryFees.forEach((f, i) => {
        if (!f.description.trim()) {
          feeDescriptionErrors[i] = 'Description is required.';
        }
        if (f.amount <= 0) {
          feeAmountErrors[i] = 'Amount is required.';
        }
        if (!f.effective_date.trim()) {
          feeEffectiveDateErrors[i] = 'Effective date is required.';
        }
      });

      if (Object.keys(feeDescriptionErrors).length > 0) {
        newErrors.feeDescription = feeDescriptionErrors;
      }
      if (Object.keys(feeAmountErrors).length > 0) {
        newErrors.feeAmount = feeAmountErrors;
      }
      if (Object.keys(feeEffectiveDateErrors).length > 0) {
        newErrors.feeEffectiveDate = feeEffectiveDateErrors;
      }
    }

    return newErrors;
  }

  function handleAddFee() {
    setMandatoryFees([
      ...mandatoryFees,
      { description: '', amount: 0, effective_date: '' },
    ]);
  }

  function handleRemoveFee(index: number) {
    setMandatoryFees(mandatoryFees.filter((_, i) => i !== index));
  }

  function handleFeeChange(
    index: number,
    field: keyof MandatoryFeeEntry,
    value: string
  ) {
    const updated = [...mandatoryFees];
    if (field === 'amount') {
      updated[index] = { ...updated[index], amount: value ? Number(value) : 0 };
    } else {
      updated[index] = { ...updated[index], [field]: value };
    }
    setMandatoryFees(updated);
  }

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setSubmitError('');
    setErrors({});

    const validationErrors = validate();
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    const data: CreateTenantRequest = {
      block: block.trim(),
      unit_number: unitNumber.trim(),
      occupancy,
      monthly_fee: Number(monthlyFee),
      ...(initialData
        ? {}
        : {
            mandatory_fees: mandatoryFees.map((f) => ({
              type: 'mandatory' as const,
              description: f.description.trim(),
              amount: f.amount,
              effective_date: f.effective_date,
            })),
          }),
    };

    setIsSubmitting(true);
    try {
      await onSubmit(data);
    } catch (err) {
      setSubmitError(
        err instanceof Error ? err.message : 'Failed to save. Check your connection and try again.'
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  const pageTitle = initialData
    ? `Edit Unit ${initialData.block}-${initialData.unit_number}`
    : 'Add New Tenant';

  return (
    <div>
      <Link
        to="/tenants"
        className="text-sm text-blue-600 hover:text-blue-700"
      >
        ← Back to Tenants
      </Link>
      <h1 className="mt-4 text-2xl font-semibold text-gray-900">{pageTitle}</h1>

      <form onSubmit={handleSubmit} noValidate>
        <div className="mx-auto mt-6 max-w-lg space-y-6 rounded-lg border border-gray-200 bg-white p-6">
          {submitError && (
            <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
              {submitError}
            </div>
          )}

          <Input
            name="block"
            label="Block"
            placeholder="A"
            value={block}
            onChange={(e: ChangeEvent<HTMLInputElement>) => setBlock(e.target.value)}
            error={errors.block}
          />

          <Input
            name="unit_number"
            label="Unit Number"
            placeholder="01"
            value={unitNumber}
            onChange={(e: ChangeEvent<HTMLInputElement>) => setUnitNumber(e.target.value)}
            error={errors.unit_number}
          />

          <Select
            name="occupancy"
            label="Occupancy Status"
            options={[
              { value: 'occupied', label: 'Occupied' },
              { value: 'vacant', label: 'Vacant' },
            ]}
            value={occupancy}
            onChange={(e: ChangeEvent<HTMLSelectElement>) =>
              setOccupancy(e.target.value as 'occupied' | 'vacant')
            }
          />

          <Input
            name="monthly_fee"
            label="Monthly Fee (Rp)"
            placeholder="50000"
            type="number"
            value={monthlyFee}
            onChange={(e: ChangeEvent<HTMLInputElement>) => setMonthlyFee(e.target.value)}
            error={errors.monthly_fee}
          />

          {/* Mandatory Fees Section (create mode only — fees managed separately on detail page) */}
          {!initialData && (
            <div className="border-t border-gray-200 pt-6">
              <h2 className="mb-4 text-sm font-semibold text-gray-700">
                Mandatory Fees
              </h2>

              {errors.mandatoryFees && (
                <p className="mb-3 text-sm text-red-600" role="alert">
                  {errors.mandatoryFees}
                </p>
              )}

              {mandatoryFees.map((fee, index) => (
                <div
                  key={index}
                  className="mb-4 rounded-md border border-gray-100 bg-gray-50 p-4"
                >
                  <div className="flex items-start justify-between">
                    <span className="text-xs font-medium text-gray-500">
                      Fee {index + 1}
                    </span>
                    {mandatoryFees.length > 1 && (
                      <button
                        type="button"
                        onClick={() => handleRemoveFee(index)}
                        className="flex min-h-[44px] min-w-[44px] items-center justify-center text-red-500 hover:text-red-700"
                        aria-label={`Remove fee ${index + 1}`}
                      >
                        <X className="h-4 w-4" />
                      </button>
                    )}
                  </div>

                  <div className="mt-2 space-y-3">
                    <Input
                      name={`fee_description_${index}`}
                      label="Description"
                      placeholder="Security Fee"
                      value={fee.description}
                      onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        handleFeeChange(index, 'description', e.target.value)
                      }
                      error={errors.feeDescription?.[index]}
                    />

                    <Input
                      name={`fee_amount_${index}`}
                      label="Amount (Rp)"
                      placeholder="25000"
                      inputMode="numeric"
                      value={fee.amount || ''}
                      onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        handleFeeChange(index, 'amount', e.target.value)
                      }
                      error={errors.feeAmount?.[index]}
                    />

                    <Input
                      name={`fee_effective_date_${index}`}
                      label="Effective Date"
                      type="date"
                      value={fee.effective_date}
                      onChange={(e: ChangeEvent<HTMLInputElement>) =>
                        handleFeeChange(index, 'effective_date', e.target.value)
                      }
                      error={errors.feeEffectiveDate?.[index]}
                    />
                  </div>
                </div>
              ))}

              <button
                type="button"
                onClick={handleAddFee}
                className="mt-3 flex items-center gap-1 text-sm font-medium text-blue-600 hover:text-blue-700"
              >
                <Plus className="h-4 w-4" />
                Add Another Fee
              </button>
            </div>
          )}

          {/* Bottom Buttons */}
          <div className="flex items-center justify-end gap-3 pt-4">
            <Link
              to="/tenants"
              className="flex min-h-[44px] items-center rounded-md border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </Link>
            <button
              type="submit"
              disabled={isLoading || isSubmitting}
              className="flex min-h-[44px] items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {isLoading || isSubmitting ? (
                <>
                  <svg
                    className="h-4 w-4 animate-spin"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    aria-hidden="true"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                    />
                  </svg>
                  Saving...
                </>
              ) : initialData ? (
                'Update Tenant'
              ) : (
                'Save Tenant'
              )}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
}
