export interface Fee {
  id: string;
  tenant_id: string;
  type?: 'mandatory' | 'voluntary';
  amount: number;
  description: string;
  effective_date: string;
  paid_at: string | null;
  created_at: string;
  updated_at?: string;
}

export interface CreateFeeRequest {
  type: 'mandatory' | 'voluntary';
  amount: number;
  description: string;
  effective_date: string;
  paid_at?: string;
}

export interface UpdateFeeRequest {
  type?: 'mandatory' | 'voluntary';
  amount?: number;
  description?: string;
  effective_date?: string;
  paid_at?: string;
}

export interface ListFeesResponse {
  mandatory_fees: Fee[];
  voluntary_fees: Fee[];
}
