import { Pencil, Trash2 } from 'lucide-react';
import StatusBadge from '../ui/StatusBadge';
import EmptyState from '../ui/EmptyState';
import LoadingSkeleton from '../ui/LoadingSkeleton';
import type { Fee } from '../../types/fee';

export function formatIDR(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
}

export function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-GB', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  });
}

interface FeeListProps {
  mandatoryFees: Fee[];
  voluntaryFees: Fee[];
  monthlyFee?: number;
  onEditFee?: (feeId: string) => void;
  onDeleteFee?: (feeId: string) => void;
  onRecordMandatory?: () => void;
  onRecordVoluntary?: () => void;
  isLoading?: boolean;
}

export default function FeeList({
  mandatoryFees = [],
  voluntaryFees = [],
  onEditFee,
  onDeleteFee,
  onRecordMandatory,
  onRecordVoluntary,
  isLoading = false,
}: FeeListProps) {
  if (isLoading) {
    return <LoadingSkeleton variant="list" count={3} />;
  }

  return (
    <div className="space-y-8">
      {/* Mandatory Fees Section */}
      <section>
        <h3 className="mb-4 text-sm font-semibold uppercase tracking-wide text-gray-500">
          Mandatory Fees
        </h3>
        {mandatoryFees.length === 0 ? (
          <EmptyState
            heading="No Mandatory Fees Set"
            body="Every tenant needs at least one mandatory fee. Add one now."
            action={
              onRecordMandatory
                ? { label: 'Record Mandatory Fee', onClick: onRecordMandatory }
                : undefined
            }
          />
        ) : (
          <div className="space-y-3">
            {mandatoryFees.map((fee) => (
              <FeeCard
                key={fee.id}
                fee={fee}
                onEdit={onEditFee}
                onDelete={onDeleteFee}
              />
            ))}
          </div>
        )}
      </section>

      {/* Voluntary Contributions Section */}
      <section>
        <h3 className="mb-4 mt-8 text-sm font-semibold uppercase tracking-wide text-gray-500">
          Voluntary Contributions
        </h3>
        {voluntaryFees.length === 0 ? (
          <EmptyState
            heading="No Voluntary Contributions Yet"
            body="Residents can contribute voluntarily here."
            action={
              onRecordVoluntary
                ? { label: 'Record Contribution', onClick: onRecordVoluntary }
                : undefined
            }
          />
        ) : (
          <div className="space-y-3">
            {voluntaryFees.map((fee) => (
              <FeeCard
                key={fee.id}
                fee={fee}
                onEdit={onEditFee}
                onDelete={onDeleteFee}
              />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}

function FeeCard({
  fee,
  onEdit,
  onDelete,
}: {
  fee: Fee;
  onEdit?: (feeId: string) => void;
  onDelete?: (feeId: string) => void;
}) {
  return (
    <div
      className="rounded-lg border border-gray-200 bg-white p-4"
      data-testid="fee-card"
    >
      <div className="flex justify-between items-start">
        <div className="space-y-1">
          <p className="text-sm font-medium text-gray-900">{fee.description}</p>
          <p className="text-sm font-semibold text-gray-900">
            {formatIDR(fee.amount)}
          </p>
          <p className="text-xs text-gray-500">
            Effective: {formatDate(fee.effective_date)}
          </p>
          <StatusBadge status={fee.paid_at ? 'paid' : 'unpaid'} />
        </div>
        <div className="flex items-start gap-1">
          {onEdit && (
            <button
              onClick={() => onEdit(fee.id)}
              className="flex min-h-[44px] min-w-[44px] items-center justify-center text-gray-400 hover:text-blue-600"
              aria-label="Edit fee"
            >
              <Pencil className="h-4 w-4" />
            </button>
          )}
          {onDelete && (
            <button
              onClick={() => onDelete(fee.id)}
              className="flex min-h-[44px] min-w-[44px] items-center justify-center text-gray-400 hover:text-red-600"
              aria-label="Delete fee"
            >
              <Trash2 className="h-4 w-4" />
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
