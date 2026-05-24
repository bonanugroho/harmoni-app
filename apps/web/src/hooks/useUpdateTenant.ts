import { useMutation, useQueryClient } from '@tanstack/react-query';
import { updateTenant } from '../services/tenants';
import type { Tenant, UpdateTenantRequest } from '../types/tenant';

export function useUpdateTenant(id: string) {
  const queryClient = useQueryClient();
  return useMutation<Tenant, Error, UpdateTenantRequest>({
    mutationFn: (data) => updateTenant(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });
}
