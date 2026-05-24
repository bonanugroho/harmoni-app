import { useState, type ReactNode } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../routes/ProtectedRoute';
import { Building2, SwatchBook, LayoutDashboard, Menu, X, ChevronLeft } from 'lucide-react';

interface AppLayoutProps {
  children: ReactNode;
}

interface NavItem {
  label: string;
  path: string;
  icon: typeof Building2;
  roles?: string[];
}

const navItems: NavItem[] = [
  { label: 'Tenants', path: '/tenants', icon: Building2 },
  { label: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { label: 'Settings', path: '/settings', icon: ChevronLeft, roles: ['rt_officer', 'rw_officer'] },
];

export default function AppLayout({ children }: AppLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();

  const visibleItems = navItems.filter((item) => {
    if (!item.roles) return true;
    return user?.role && item.roles.includes(user.role);
  });

  const isActive = (path: string) => {
    if (path === '/tenants') return location.pathname.startsWith('/tenants');
    return location.pathname.startsWith(path);
  };

  const handleNavigate = (path: string) => {
    setSidebarOpen(false);
    navigate(path);
  };

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* Skip link */}
      <a
        href="#main-content"
        className="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-60 focus:rounded-md focus:bg-white focus:px-4 focus:py-2 focus:text-sm focus:text-blue-600 focus:shadow-lg"
      >
        Skip to content
      </a>

      {/* Mobile backdrop overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-gray-900/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
          aria-hidden="true"
        />
      )}

      {/* Sidebar - fixed on all sizes, hidden on mobile via translate, always visible on lg */}
      <aside
        className={`fixed inset-y-0 left-0 z-50 flex w-64 flex-col bg-gray-100 shadow-sm transition-transform duration-200 ease-in-out lg:translate-x-0 ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        {/* Brand header */}
        <div className="flex h-16 items-center justify-between border-b border-gray-200 px-4">
          <button
            onClick={() => handleNavigate('/dashboard')}
            className="flex items-center gap-2 text-[28px] font-semibold text-gray-900"
          >
            <SwatchBook className="h-7 w-7" />
            Harmoni
          </button>
          <button
            onClick={() => setSidebarOpen(false)}
            className="flex min-h-[44px] min-w-[44px] items-center justify-center lg:hidden"
            aria-label="Close sidebar"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Navigation */}
        <nav className="mt-4 flex-1 space-y-1 px-3">
          {visibleItems.map((item) => {
            const Icon = item.icon;
            const active = isActive(item.path);
            return (
              <button
                key={item.path}
                onClick={() => handleNavigate(item.path)}
                className={`flex min-h-[44px] w-full items-center gap-3 rounded-md px-4 py-2.5 text-sm font-medium transition-colors ${
                  active
                    ? 'border-l-4 border-blue-600 bg-blue-100 text-blue-700'
                    : 'border-l-4 border-transparent text-gray-700 hover:bg-gray-200'
                }`}
              >
                <Icon className="h-5 w-5 flex-shrink-0" />
                <span>{item.label}</span>
              </button>
            );
          })}
        </nav>
      </aside>

      {/* Main content area - left margin on desktop matches sidebar width */}
      <div className="flex flex-1 flex-col lg:ml-64">
        {/* Mobile header with hamburger */}
        <header className="flex h-16 items-center gap-4 border-b border-gray-200 bg-white px-4 lg:hidden">
          <button
            onClick={() => setSidebarOpen(true)}
            className="flex min-h-[44px] min-w-[44px] items-center justify-center"
            aria-label="Toggle sidebar"
          >
            <Menu className="h-5 w-5" />
          </button>
        </header>

        {/* Content */}
        <main id="main-content" className="flex-1 p-4 lg:p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
