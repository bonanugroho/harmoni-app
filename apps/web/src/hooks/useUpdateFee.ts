import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateFee } from '../services/fees';
import type { UpdateFeeRequest } from '../types/fee';

export function useUpdateFee(tenantId: string) {
  const queryClient = useQueryClient();
  return useMutation<void, Error, { feeId: string; data: UpdateFeeRequest }>({
    mutationFn: ({ feeId, data }) => updateFee(tenantId, feeId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fees', tenantId] });
    },
  });
}
