import { useNavigate } from 'react-router-dom';
import { useCreateTenant } from '../hooks/useCreateTenant';
import TenantForm from '../components/tenants/TenantForm';
import type { CreateTenantRequest } from '../types/tenant';

export default function TenantCreatePage() {
  const navigate = useNavigate();
  const mutation = useCreateTenant();

  async function handleSubmit(data: CreateTenantRequest) {
    try {
      await mutation.mutateAsync(data);
      navigate('/tenants');
    } catch (error) {
      throw error;
    }
  }

  return (
    <TenantForm onSubmit={handleSubmit} isLoading={mutation.isPending} />
  );
}
