import React, { ReactNode, useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { adminService } from '../../services';

interface AdminLayoutProps {
  children: ReactNode;
}

const adminNav = [
  { path: '/admin', label: 'Overview', exact: true, badge: '' },
  { path: '/admin/users', label: 'Users', badge: '' },
  { path: '/admin/pending-users', label: 'Pending', badge: 'pending_users' },
  { path: '/admin/items', label: 'Items', badge: 'low_stock' },
  { path: '/admin/suppliers', label: 'Suppliers', badge: '' },
  { path: '/admin/requests', label: 'Requests', badge: 'pending_requests' },
  { path: '/admin/issues/new', label: 'Issue Item', badge: '' },
  { path: '/admin/issues', label: 'Issues', badge: 'overdue' },
  { path: '/admin/notices', label: 'Notices', badge: '' },
];

export function AdminLayout({ children }: AdminLayoutProps) {
  const { user, logout } = useAuth();
  const location = useLocation();
  const [counts, setCounts] = useState<Record<string, number>>({});

  useEffect(() => {
    adminService.getOverview().then((overview) => {
      setCounts({
        low_stock: overview.low_stock_items,
        pending_requests: overview.pending_requests,
        overdue: overview.overdue_items,
      });
    });
    adminService.getPendingUsers('', { page: 1, page_size: 1 }).then((res) => {
      setCounts(prev => ({ ...prev, pending_users: res.total_count }));
    });
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link to="/admin" className="text-xl font-bold text-blue-600 flex items-center gap-2">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.486-1.027.764-2.155.764-3.317 0-3.866-3.134-7-7-7s-7 3.134-7 7c0 1.162.278 2.29.764 3.317M13 20h3v-3h-3v3zM7 20h3v-3H7v3z" />
                </svg>
                Admin
              </Link>
              <nav className="ml-8 flex space-x-1">
                {adminNav.map((item) => {
                  const isActive = item.exact
                    ? location.pathname === item.path
                    : location.pathname.startsWith(item.path);
                  const badge = item.badge ? counts[item.badge] : 0;
                  return (
                    <Link
                      key={item.path}
                      to={item.path}
                      className={`relative px-3 py-2 rounded-md text-sm font-medium ${
                        isActive
                          ? 'bg-blue-100 text-blue-700'
                          : 'text-gray-600 hover:bg-gray-100'
                      }`}
                    >
                      {item.label}
                      {badge > 0 && (
                        <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                          {badge > 9 ? '9+' : badge}
                        </span>
                      )}
                    </Link>
                  );
                })}
              </nav>
            </div>
            <div className="flex items-center space-x-4">
              <Link to="/dashboard" className="text-sm text-gray-500 hover:text-gray-700 flex items-center gap-1">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                </svg>
                User Dashboard
              </Link>
              <span className="text-sm text-gray-700">{user?.username}</span>
              <button
                onClick={logout}
                className="text-sm text-gray-500 hover:text-gray-700 flex items-center gap-1"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                Logout
              </button>
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