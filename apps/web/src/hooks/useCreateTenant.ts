import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createTenant } from '../services/tenants';
import type { Tenant, CreateTenantRequest } from '../types/tenant';

export function useCreateTenant() {
  const queryClient = useQueryClient();
  return useMutation<Tenant, Error, CreateTenantRequest>({
    mutationFn: (data) => createTenant(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });
}
