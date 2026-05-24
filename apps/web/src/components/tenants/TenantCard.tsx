import type { Tenant } from '../../types/tenant';
import StatusBadge from '../ui/StatusBadge';

function formatIDR(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
}

interface TenantCardProps {
  tenant: Tenant;
  onClick?: () => void;
  mandatoryFeeCount?: number;
  voluntaryFeeCount?: number;
}

export default function TenantCard({
  tenant,
  onClick,
  mandatoryFeeCount = 0,
  voluntaryFeeCount = 0,
}: TenantCardProps) {
  const hasFees = mandatoryFeeCount > 0 || voluntaryFeeCount > 0;

  return (
    <div
      className="min-h-[100px] cursor-pointer rounded-lg border border-gray-200 bg-white p-4 shadow-sm transition-shadow hover:shadow-md active:scale-[0.98]"
      onClick={onClick}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => e.key === 'Enter' && onClick?.()}
    >
      <div className="flex items-start justify-between">
        <span className="text-base font-semibold text-gray-900">
          Block {tenant.block} · Unit {tenant.unit_number}
        </span>
        <StatusBadge status={tenant.occupancy} />
      </div>
      <p className="mt-2 text-sm font-semibold text-gray-900">
        {formatIDR(tenant.monthly_fee)} / month
      </p>
      <p className="mt-1 text-sm text-gray-500">
        {hasFees
          ? `${mandatoryFeeCount} mandatory fees · ${voluntaryFeeCount} contributions`
          : 'No fees configured'}
      </p>
    </div>
  );
}

export { formatIDR };
