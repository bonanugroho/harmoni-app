import { useState, useEffect, useMemo, type ChangeEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { Search } from 'lucide-react';
import { useTenants } from '../hooks/useTenants';
import TenantCard from '../components/tenants/TenantCard';
import PageHeader from '../components/ui/PageHeader';
import LoadingSkeleton from '../components/ui/LoadingSkeleton';
import EmptyState from '../components/ui/EmptyState';

export default function TenantListPage() {
  const navigate = useNavigate();
  const { data: tenants, isLoading, isError, error } = useTenants();
  const [searchQuery, setSearchQuery] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');

  // 300ms debounce for search
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(searchQuery);
    }, 300);
    return () => clearTimeout(timer);
  }, [searchQuery]);

  const filteredTenants = useMemo(() => {
    if (!tenants) return [];
    if (!debouncedSearch.trim()) return tenants;
    const q = debouncedSearch.toLowerCase();
    return tenants.filter(
      (t) =>
        t.block.toLowerCase().includes(q) ||
        t.unit_number.toLowerCase().includes(q)
    );
  }, [tenants, debouncedSearch]);

  const sortedTenants = useMemo(() => {
    return [...filteredTenants].sort((a, b) => {
      if (a.block !== b.block) return a.block.localeCompare(b.block);
      return a.unit_number.localeCompare(b.unit_number);
    });
  }, [filteredTenants]);

  function handleSearchChange(e: ChangeEvent<HTMLInputElement>) {
    setSearchQuery(e.target.value);
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Tenants"
        action={{
          label: '+ Add Tenant',
          onClick: () => navigate('/tenants/new'),
        }}
      />

      {/* Search Bar */}
      <div className="relative max-w-md">
        <Search className="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400" />
        <input
          type="search"
          name="search"
          placeholder="Search by block or unit number..."
          value={searchQuery}
          onChange={handleSearchChange}
          className="min-h-[44px] w-full rounded-md border border-gray-200 py-2 pl-10 pr-3 text-base outline-none transition-colors focus:border-blue-600"
          aria-label="Search tenants"
        />
      </div>

      {/* Loading State */}
      {isLoading && <LoadingSkeleton variant="card" count={3} />}

      {/* Error State */}
      {isError && (
        <div className="rounded-md bg-red-50 p-4 text-sm text-red-700" role="alert">
          {error?.message || 'Failed to load tenants. Pull down to refresh or try again.'}
        </div>
      )}

      {/* Empty State (no tenants at all) */}
      {!isLoading && !isError && (!tenants || tenants.length === 0) && (
        <EmptyState
          heading="No Tenants Yet"
          body="Start by adding your first tenant to this RT."
        />
      )}

      {/* Empty State (search no results) */}
      {!isLoading && !isError && tenants && tenants.length > 0 && sortedTenants.length === 0 && (
        <EmptyState
          heading="No Results Found"
          body={`No tenants match "${searchQuery}". Try a different search.`}
        />
      )}

      {/* Normal State */}
      {!isLoading && !isError && sortedTenants.length > 0 && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {sortedTenants.map((tenant) => (
            <TenantCard
              key={tenant.id}
              tenant={tenant}
              onClick={() => navigate(`/tenants/${tenant.id}`)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
