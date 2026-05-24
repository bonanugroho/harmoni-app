import { useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteFee } from '../services/fees';

export function useDeleteFee(tenantId: string) {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: (feeId) => deleteFee(tenantId, feeId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fees', tenantId] });
    },
  });
}
