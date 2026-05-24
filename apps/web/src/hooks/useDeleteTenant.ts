import { useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteTenant } from '../services/tenants';

export function useDeleteTenant() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: (id) => deleteTenant(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });
}
