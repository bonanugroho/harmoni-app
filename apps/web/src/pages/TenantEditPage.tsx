import { useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useTenant } from '../hooks/useTenant';
import { useUpdateTenant } from '../hooks/useUpdateTenant';
import { useDeleteTenant } from '../hooks/useDeleteTenant';
import TenantForm from '../components/tenants/TenantForm';
import ConfirmDialog from '../components/ui/ConfirmDialog';
import LoadingSkeleton from '../components/ui/LoadingSkeleton';
import type { CreateTenantRequest } from '../types/tenant';

export default function TenantEditPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data: tenant, isLoading, isError } = useTenant(id || '');
  const updateMutation = useUpdateTenant(id || '');
  const deleteMutation = useDeleteTenant();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  async function handleUpdate(data: CreateTenantRequest) {
    try {
      await updateMutation.mutateAsync({
        block: data.block,
        unit_number: data.unit_number,
        occupancy: data.occupancy,
        monthly_fee: data.monthly_fee,
      });
      navigate('/tenants');
    } catch (error) {
      throw error;
    }
  }

  async function handleDelete() {
    try {
      await deleteMutation.mutateAsync(id || '');
      navigate('/tenants');
    } catch (error) {
      setShowDeleteDialog(false);
    }
  }

  if (isLoading) {
    return <LoadingSkeleton variant="form" />;
  }

  if (isError || !tenant) {
    return (
      <div className="space-y-4">
        <div
          role="alert"
          className="rounded-md bg-red-50 p-4 text-sm text-red-700"
        >
          Failed to load tenant. It may have been removed or you may not have
          access.
        </div>
        <Link
          to="/tenants"
          className="inline-flex items-center text-sm text-blue-600 hover:text-blue-700"
        >
          ← Back to Tenants
        </Link>
      </div>
    );
  }

  return (
    <div>
      <TenantForm
        initialData={tenant}
        onSubmit={handleUpdate}
        isLoading={updateMutation.isPending}
      />

      {/* Danger Zone */}
      <div className="mx-auto mt-8 max-w-lg border-t border-gray-200 pt-6">
        <button
          type="button"
          onClick={() => setShowDeleteDialog(true)}
          className="flex min-h-[44px] items-center rounded-md border border-red-300 bg-white px-4 py-2 text-sm font-medium text-red-600 hover:bg-red-50"
        >
          Delete Tenant
        </button>
      </div>

      <ConfirmDialog
        isOpen={showDeleteDialog}
        title={`Delete Unit ${tenant.block}-${tenant.unit_number}?`}
        message="This cannot be undone."
        confirmLabel="Delete"
        cancelLabel="Cancel"
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteDialog(false)}
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
}
