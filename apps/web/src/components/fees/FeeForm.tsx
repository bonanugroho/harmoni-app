import type { Fee, CreateFeeRequest } from '../../types/fee';

interface FeeFormProps {
  tenantId: string;
  monthlyFee: number;
  initialData?: Fee;
  onSubmit: (data: CreateFeeRequest) => Promise<void>;
  isLoading?: boolean;
  onCancel: () => void;
}

export default function FeeForm(_props: FeeFormProps) {
  return <div>FeeForm Placeholder</div>;
}
