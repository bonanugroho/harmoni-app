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
  monthlyFee,
  onEditFee,
  onDeleteFee,
  onRecordMandatory,
  onRecordVoluntary,
  isLoading = false,
}: FeeListProps) {
  return (
    <div>
      <div>FeeList Placeholder</div>
    </div>
  );
}
