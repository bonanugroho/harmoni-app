import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ResetPasswordPage from './pages/ResetPasswordPage';
import TenantListPage from './pages/TenantListPage';
import TenantCreatePage from './pages/TenantCreatePage';
import TenantEditPage from './pages/TenantEditPage';
import TenantDetailPage from './pages/TenantDetailPage';
import ProtectedRoute from './routes/ProtectedRoute';
import AppLayout from './components/layout/AppLayout';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
      refetchOnWindowFocus: true,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/reset" element={<ResetPasswordPage />} />
          <Route
            path="/tenants"
            element={
              <ProtectedRoute>
                <AppLayout><TenantListPage /></AppLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/tenants/new"
            element={
              <ProtectedRoute>
                <AppLayout><TenantCreatePage /></AppLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/tenants/:id"
            element={
              <ProtectedRoute>
                <AppLayout><TenantDetailPage /></AppLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/tenants/:id/edit"
            element={
              <ProtectedRoute>
                <AppLayout><TenantEditPage /></AppLayout>
              </ProtectedRoute>
            }
          />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <AppLayout>
                  <div className="flex min-h-screen items-center justify-center">
                    <div className="text-center">
                      <h1 className="text-2xl font-semibold text-gray-900">Welcome to Harmoni</h1>
                      <p className="mt-2 text-gray-600">Your dashboard is coming soon.</p>
                    </div>
                  </div>
                </AppLayout>
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}
