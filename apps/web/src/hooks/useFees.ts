import { useQuery } from '@tanstack/react-query';
import { listFees } from '../services/fees';
import type { ListFeesResponse } from '../types/fee';

export function useFees(tenantId: string) {
  return useQuery<ListFeesResponse>({
    queryKey: ['fees', tenantId],
    queryFn: () => listFees(tenantId),
    enabled: !!tenantId,
  });
}
