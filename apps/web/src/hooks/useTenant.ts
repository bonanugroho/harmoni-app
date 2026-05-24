import { useQuery } from '@tanstack/react-query';
import { getTenant } from '../services/tenants';
import type { Tenant } from '../types/tenant';

export function useTenant(id: string) {
  return useQuery<Tenant>({
    queryKey: ['tenants', id],
    queryFn: () => getTenant(id),
    enabled: !!id,
  });
}
