import React, { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

interface LayoutProps {
  children: ReactNode;
}

const navItems = [
  { path: '/dashboard', label: 'Dashboard', roles: [1, 2] },
  { path: '/inventory', label: 'Inventory', roles: [1, 2] },
  { path: '/requests', label: 'My Requests', roles: [1, 2] },
  { path: '/issues', label: 'My Issues', roles: [1, 2] },
  { path: '/notices', label: 'Notices', roles: [1, 2, 3] },
  { path: '/admin', label: 'Admin', roles: [1] },
];

export function AppLayout({ children }: LayoutProps) {
  const { user, logout } = useAuth();
  const location = useLocation();

  const filteredNav = navItems.filter(
    (item) => item.roles.includes(user?.role_id as number)
  );

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link to="/" className="text-xl font-bold text-blue-600">
                Inventory
              </Link>
              <nav className="ml-8 flex space-x-4">
                {filteredNav.map((item) => (
                  <Link
                    key={item.path}
                    to={item.path}
                    className={`px-3 py-2 rounded-md text-sm font-medium ${
                      location.pathname.startsWith(item.path)
                        ? 'bg-blue-100 text-blue-700'
                        : 'text-gray-700 hover:bg-gray-100'
                    }`}
                  >
                    {item.label}
                  </Link>
                ))}
              </nav>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-700">{user?.username}</span>
              <button
                onClick={logout}
                className="text-sm text-gray-500 hover:text-gray-700"
              >
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