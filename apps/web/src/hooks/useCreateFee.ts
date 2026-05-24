import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createFee } from '../services/fees';
import type { Fee, CreateFeeRequest } from '../types/fee';

export function useCreateFee(tenantId: string) {
  const queryClient = useQueryClient();
  return useMutation<Fee, Error, CreateFeeRequest>({
    mutationFn: (data) => createFee(tenantId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fees', tenantId] });
    },
  });
}
