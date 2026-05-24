interface StatusBadgeProps {
  status: 'occupied' | 'vacant' | 'paid' | 'unpaid';
}

const styles: Record<string, string> = {
  occupied: 'bg-blue-100 text-blue-700',
  vacant: 'bg-amber-100 text-amber-800',
  paid: 'bg-emerald-100 text-emerald-800',
  unpaid: 'bg-red-100 text-red-800',
};

const labels: Record<string, string> = {
  occupied: 'Occupied',
  vacant: 'Vacant',
  paid: 'Paid',
  unpaid: 'Unpaid',
};

export default function StatusBadge({ status }: StatusBadgeProps) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${styles[status]}`}
    >
      {labels[status]}
    </span>
  );
}
