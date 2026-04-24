import { Routes, Route, Navigate } from 'react-router-dom';
import { PublicLayout, AppLayout, AdminLayout } from '../components/layout';
import { ProtectedRoute } from '../components/auth';

import { LoginPage } from '../pages/auth/LoginPage';
import { StatsPage } from '../pages/public/StatsPage';
import { NoticesPage } from '../pages/public/NoticesPage';

import { DashboardPage } from '../pages/user/DashboardPage';
import { InventoryPage } from '../pages/user/InventoryPage';
import { RequestFormPage } from '../pages/user/RequestFormPage';
import { MyRequestsPage } from '../pages/user/MyRequestsPage';
import { MyIssuesPage } from '../pages/user/MyIssuesPage';

import { AdminDashboardPage } from '../pages/admin/AdminDashboardPage';
import { AdminUsersPage } from '../pages/admin/AdminUsersPage';
import { AdminItemsPage } from '../pages/admin/AdminItemsPage';
import { AdminItemFormPage } from '../pages/admin/AdminItemFormPage';
import { AdminItemDetailsPage } from '../pages/admin/AdminItemDetailsPage';
import { AdminSuppliersPage } from '../pages/admin/AdminSuppliersPage';
import { AdminRequestsPage } from '../pages/admin/AdminRequestsPage';
import { AdminIssueFormPage } from '../pages/admin/AdminIssueFormPage';
import { AdminIssuesPage } from '../pages/admin/AdminIssuesPage';
import { AdminNoticesPage } from '../pages/admin/AdminNoticesPage';

import { Role } from '../types';

export function AppRoutes() {
  return (
    <Routes>
      <Route element={<PublicLayout />}>
        <Route path="/" element={<StatsPage />} />
        <Route path="/notices" element={<NoticesPage />} />
        <Route path="/login" element={<LoginPage />} />
      </Route>

      <Route
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/inventory" element={<InventoryPage />} />
        <Route path="/inventory/:id" element={<InventoryPage />} />
        <Route path="/requests/new" element={<RequestFormPage />} />
        <Route path="/requests" element={<MyRequestsPage />} />
        <Route path="/issues" element={<MyIssuesPage />} />
      </Route>

      <Route
        path="/admin"
        element={
          <ProtectedRoute requiredRole={1 as Role}>
            <AdminLayout />
          </ProtectedRoute>
        }
      >
        <Route index element={<AdminDashboardPage />} />
        <Route path="users" element={<AdminUsersPage />} />
        <Route path="items" element={<AdminItemsPage />} />
        <Route path="items/new" element={<AdminItemFormPage />} />
        <Route path="items/:id" element={<AdminItemDetailsPage />} />
        <Route path="items/:id/edit" element={<AdminItemFormPage isEdit={true} />} />
        <Route path="suppliers" element={<AdminSuppliersPage />} />
        <Route path="requests" element={<AdminRequestsPage />} />
        <Route path="issues/new" element={<AdminIssueFormPage />} />
        <Route path="issues" element={<AdminIssuesPage />} />
        <Route path="notices" element={<AdminNoticesPage />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}