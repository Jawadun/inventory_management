import React, { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';

interface PublicLayoutProps {
  children: ReactNode;
}

export function PublicLayout({ children }: PublicLayoutProps) {
  const location = useLocation();
  const isAuthPage = location.pathname === '/login';

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
              {!isAuthPage && (
                <Link
                  to="/login"
                  className="text-sm text-gray-500 hover:text-gray-700"
                >
                  Login
                </Link>
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