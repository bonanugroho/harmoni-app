import { useQuery } from '@tanstack/react-query';
import { listTenants } from '../services/tenants';
import type { Tenant } from '../types/tenant';

export function useTenants() {
  return useQuery<Tenant[]>({
    queryKey: ['tenants'],
    queryFn: listTenants,
  });
}
