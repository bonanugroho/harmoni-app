import { useState, useEffect, createContext, useContext } from 'react';
import { Navigate, useLocation } from 'react-router-dom';

// Auth context for sharing user data
export const AuthContext = createContext(null);

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

/**
 * ProtectedRoute - redirects to /login if not authenticated.
 * Optionally restricts access by role.
 */
export default function ProtectedRoute({ children, requiredRole }) {
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const [user, setUser] = useState(null);
  const location = useLocation();

  useEffect(() => {
    // Check authentication status via cookie
    // Since cookies are httpOnly, we check via a /me endpoint or similar
    fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:3000'}/auth/me`, {
      credentials: 'include',
    })
      .then((res) => {
        if (!res.ok) {
          setIsAuthenticated(false);
          setUser(null);
          return;
        }
        return res.json();
      })
      .then((data) => {
        if (data) {
          setIsAuthenticated(true);
          setUser(data.user);
        }
      })
      .catch(() => {
        setIsAuthenticated(false);
        setUser(null);
      });
  }, []);

  // Loading state while checking auth status
  if (isAuthenticated === null) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <div className="mx-auto h-8 w-8 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
          <p className="mt-4 text-sm text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  // Redirect to /login if not authenticated
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Optional role-based access control
  if (requiredRole && user?.role !== requiredRole) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="rounded-md bg-red-50 p-6 text-center">
          <h2 className="text-lg font-semibold text-red-700">Access Denied</h2>
          <p className="mt-2 text-sm text-red-600">
            You do not have permission to access this page.
          </p>
        </div>
      </div>
    );
  }

  // Render children with auth context
  return <AuthContext.Provider value={{ user, isAuthenticated }}>{children}</AuthContext.Provider>;
}
