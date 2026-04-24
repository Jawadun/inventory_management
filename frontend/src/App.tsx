import { BrowserRouter, Routes, Route, Link, Navigate, useNavigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { PublicStats } from './types';

import { StatsPage } from './pages/public/StatsPage';
import { NoticesPage } from './pages/public/NoticesPage';
import { LoginPage } from './pages/auth/LoginPage';
import { RegisterPage } from './pages/auth/RegisterPage';

import { DashboardPage } from './pages/user/DashboardPage';
import { InventoryPage } from './pages/user/InventoryPage';
import { RequestFormPage } from './pages/user/RequestFormPage';
import { MyRequestsPage } from './pages/user/MyRequestsPage';
import { MyIssuesPage } from './pages/user/MyIssuesPage';

import { AdminDashboardPage } from './pages/admin/AdminDashboardPage';
import { AdminUsersPage } from './pages/admin/AdminUsersPage';
import { AdminItemsPage } from './pages/admin/AdminItemsPage';
import { AdminItemFormPage } from './pages/admin/AdminItemFormPage';
import { AdminItemDetailsPage } from './pages/admin/AdminItemDetailsPage';
import { AdminSuppliersPage } from './pages/admin/AdminSuppliersPage';
import { AdminRequestsPage } from './pages/admin/AdminRequestsPage';
import { AdminIssueFormPage } from './pages/admin/AdminIssueFormPage';
import { AdminIssuesPage } from './pages/admin/AdminIssuesPage';
import { AdminNoticesPage } from './pages/admin/AdminNoticesPage';
import { AdminPendingUsersPage } from './pages/admin/AdminPendingUsersPage';

import { ProtectedRoute } from './components/auth/ProtectedRoute';
import { AdminLayout } from './components/layout/AdminLayout';

import './index.css';

function PublicLayout({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, logout, user } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link to="/" className="text-xl font-bold text-blue-600">
                Inventory System
              </Link>
            </div>
            <div className="flex items-center space-x-4">
              {isAuthenticated ? (
                <>
                  <span className="text-sm text-gray-500">{user?.username}</span>
                  <button
                    onClick={handleLogout}
                    className="text-sm text-red-600 hover:text-red-700"
                  >
                    Logout
                  </button>
                </>
              ) : (
                <>
                  <Link to="/register" className="text-sm text-gray-500 hover:text-gray-700">
                    Register
                  </Link>
                  <Link to="/login" className="text-sm text-gray-500 hover:text-gray-700">
                    Login
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
    </div>
  );
}

function AppLayout({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, logout, user } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link to="/dashboard" className="text-xl font-bold text-blue-600">
                Inventory System
              </Link>
            </div>
            <div className="flex items-center space-x-4">
              <Link to="/dashboard" className="text-sm text-gray-500 hover:text-gray-700">Dashboard</Link>
              <Link to="/inventory" className="text-sm text-gray-500 hover:text-gray-700">Inventory</Link>
              <Link to="/requests" className="text-sm text-gray-500 hover:text-gray-700">Requests</Link>
              <Link to="/issues" className="text-sm text-gray-500 hover:text-gray-700">Issues</Link>
              {isAuthenticated && (
                <>
                  <span className="text-sm text-gray-500">{user?.username}</span>
                  <button
                    onClick={handleLogout}
                    className="text-sm text-red-600 hover:text-red-700"
                  >
                    Logout
                  </button>
                </>
              )}
            </div>
          </div>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
    </div>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<PublicLayout><StatsPage /></PublicLayout>} />
          <Route path="/notices" element={<PublicLayout><NoticesPage /></PublicLayout>} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          <Route path="/dashboard" element={<ProtectedRoute><AppLayout><DashboardPage /></AppLayout></ProtectedRoute>} />
          <Route path="/inventory" element={<ProtectedRoute><AppLayout><InventoryPage /></AppLayout></ProtectedRoute>} />
          <Route path="/inventory/:id" element={<ProtectedRoute><AppLayout><InventoryPage /></AppLayout></ProtectedRoute>} />
          <Route path="/requests/new" element={<ProtectedRoute><AppLayout><RequestFormPage /></AppLayout></ProtectedRoute>} />
          <Route path="/requests" element={<ProtectedRoute><AppLayout><MyRequestsPage /></AppLayout></ProtectedRoute>} />
          <Route path="/issues" element={<ProtectedRoute><AppLayout><MyIssuesPage /></AppLayout></ProtectedRoute>} />

          <Route path="/admin" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminDashboardPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/users" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminUsersPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/items" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminItemsPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/items/new" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminItemFormPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/items/:id" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminItemDetailsPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/items/:id/edit" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminItemFormPage isEdit={true} /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/suppliers" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminSuppliersPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/requests" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminRequestsPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/issues/new" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminIssueFormPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/issues" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminIssuesPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/notices" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminNoticesPage /></AdminLayout></ProtectedRoute>} />
          <Route path="/admin/pending-users" element={<ProtectedRoute requiredRole={1 as any}><AdminLayout><AdminPendingUsersPage /></AdminLayout></ProtectedRoute>} />

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}