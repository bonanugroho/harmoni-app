import { useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { Pencil } from 'lucide-react';
import { useTenant } from '../hooks/useTenant';
import { useFees } from '../hooks/useFees';
import { useCreateFee } from '../hooks/useCreateFee';
import { useUpdateFee } from '../hooks/useUpdateFee';
import { useDeleteFee } from '../hooks/useDeleteFee';
import PageHeader from '../components/ui/PageHeader';
import StatusBadge from '../components/ui/StatusBadge';
import LoadingSkeleton from '../components/ui/LoadingSkeleton';
import FeeList from '../components/fees/FeeList';
import FeeForm from '../components/fees/FeeForm';
import ConfirmDialog from '../components/ui/ConfirmDialog';
import type { Fee, CreateFeeRequest } from '../types/fee';
import { formatIDR } from '../components/fees/FeeList';

export default function TenantDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: tenant, isLoading: tenantLoading, isError: tenantError } = useTenant(id!);
  const { data: fees, isLoading: feesLoading } = useFees(id!);

  const createFeeMutation = useCreateFee(id!);
  const updateFeeMutation = useUpdateFee(id!);
  const deleteFeeMutation = useDeleteFee(id!);

  const [showFeeForm, setShowFeeForm] = useState(false);
  const [editingFee, setEditingFee] = useState<Fee | null>(null);
  const [deleteFeeId, setDeleteFeeId] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState('');

  function handleCloseFeeForm() {
    setShowFeeForm(false);
    setEditingFee(null);
    setSubmitError('');
  }

  async function handleFeeSubmit(data: CreateFeeRequest) {
    try {
      if (editingFee) {
        await updateFeeMutation.mutateAsync({ feeId: editingFee.id, data });
      } else {
        await createFeeMutation.mutateAsync(data);
      }
      handleCloseFeeForm();
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'An unexpected error occurred');
    }
  }

  async function handleDeleteFee() {
    if (!deleteFeeId) return;
    try {
      await deleteFeeMutation.mutateAsync(deleteFeeId);
      setDeleteFeeId(null);
    } catch {
      // Error handled by mutation
    }
  }

  const existingMandatoryTotal = (fees?.mandatory_fees || []).reduce(
    (sum, f) => sum + Number(f.amount), 0
  );

  function handleEditFee(feeId: string) {
    const allFees = [
      ...(fees?.mandatory_fees || []),
      ...(fees?.voluntary_fees || []),
    ];
    const fee = allFees.find((f) => f.id === feeId);
    if (fee) {
      setEditingFee(fee);
      setShowFeeForm(true);
    }
  }

  if (tenantLoading) {
    return <LoadingSkeleton variant="form" />;
  }

  if (tenantError || !tenant) {
    return (
      <div className="space-y-4">
        <Link to="/tenants" className="text-sm text-blue-600 hover:text-blue-700">
          ← Back to Tenants
        </Link>
        <div role="alert" className="rounded-md bg-red-50 p-4 text-sm text-red-700">
          Failed to load tenant. Please try again.
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Link to="/tenants" className="text-sm text-blue-600 hover:text-blue-700">
        ← Back to Tenants
      </Link>

      <PageHeader
        title={`Unit ${tenant.block}-${tenant.unit_number} — Fees`}
        action={{
          label: 'Record Fee',
          onClick: () => setShowFeeForm(true),
        }}
      />

      {/* Tenant info summary bar */}
      <div className="flex items-center gap-4 rounded-lg bg-white p-4 shadow-sm">
        <span className="text-sm text-gray-700">
          Block {tenant.block} · Unit {tenant.unit_number}
        </span>
        <StatusBadge status={tenant.occupancy} />
        <span className="text-sm font-semibold text-gray-900">
          {formatIDR(tenant.monthly_fee)}/mo
        </span>
        <button
          onClick={() => navigate(`/tenants/${tenant.id}/edit`)}
          className="ml-auto flex min-h-[44px] min-w-[44px] items-center justify-center gap-1.5 rounded-md px-3 text-sm font-medium text-blue-600 hover:bg-blue-50"
          aria-label={`Edit ${tenant.block}-${tenant.unit_number}`}
        >
          <Pencil className="h-4 w-4" />
          Edit
        </button>
      </div>

      <FeeList
        mandatoryFees={fees?.mandatory_fees || []}
        voluntaryFees={fees?.voluntary_fees || []}
        monthlyFee={tenant.monthly_fee}
        onEditFee={handleEditFee}
        onDeleteFee={(feeId) => setDeleteFeeId(feeId)}
        onRecordMandatory={() => setShowFeeForm(true)}
        onRecordVoluntary={() => setShowFeeForm(true)}
        isLoading={feesLoading}
      />

      {/* FeeForm Modal */}
      {showFeeForm && (
        <div className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto pt-10 sm:items-center">
          <div
            className="fixed inset-0 bg-gray-900/50"
            onClick={handleCloseFeeForm}
          />
          <div className="relative z-10 w-full max-w-md rounded-lg bg-white p-6 shadow-xl mx-4">
            <h2 className="mb-4 text-lg font-semibold text-gray-900">
              {editingFee ? 'Edit Fee' : 'Record Fee'}
            </h2>
            <FeeForm
              tenantId={id!}
              monthlyFee={tenant.monthly_fee}
              existingMandatoryTotal={existingMandatoryTotal}
              editingFeeId={editingFee?.id}
              initialData={editingFee || undefined}
              onSubmit={handleFeeSubmit}
              isLoading={createFeeMutation.isPending || updateFeeMutation.isPending}
              onCancel={handleCloseFeeForm}
            />
          </div>
        </div>
      )}

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={!!deleteFeeId}
        title="Delete this fee?"
        message="This cannot be undone."
        onConfirm={handleDeleteFee}
        onCancel={() => setDeleteFeeId(null)}
        isLoading={deleteFeeMutation.isPending}
      />
    </div>
  );
}
