import React, { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { Role } from '../../types';

interface ProtectedRouteProps {
  children: ReactNode;
  requiredRole?: Role;
  fallbackPath?: string;
}

export function ProtectedRoute({ children, requiredRole, fallbackPath = '/login' }: ProtectedRouteProps) {
  const { user, isAuthenticated, isLoading } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to={fallbackPath} state={{ from: location }} replace />;
  }

  if (requiredRole) {
    const hasAccess = checkRoleAccess(user?.role_id, requiredRole);
    if (!hasAccess) {
      return <Navigate to="/" replace />;
    }
  }

  return <>{children}</>;
}

function checkRoleAccess(userRole: number | undefined, requiredRole: Role): boolean {
  if (!userRole) return false;
  if (requiredRole === 1) return userRole === 1;
  if (requiredRole === 2) return userRole === 1 || userRole === 2;
  return true;
}