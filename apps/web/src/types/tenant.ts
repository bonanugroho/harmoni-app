export interface Tenant {
  id: string;
  block: string;
  unit_number: string;
  occupancy: 'occupied' | 'vacant';
  monthly_fee: number;
  territory_id: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTenantRequest {
  block: string;
  unit_number: string;
  occupancy: 'occupied' | 'vacant';
  monthly_fee: number;
  mandatory_fees: import('./fee').CreateFeeRequest[];
}

export interface UpdateTenantRequest {
  block?: string;
  unit_number?: string;
  occupancy?: 'occupied' | 'vacant';
  monthly_fee?: number;
}
