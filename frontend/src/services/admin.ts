import api from './api';
import { AdminOverview, AnalyticsData, AdminFilter, PaginationParams } from '../types';

export const adminService = {
  getDashboard: async () => {
    const response = await api.get('/admin/dashboard');
    return response.data;
  },

  getOverview: async (): Promise<AdminOverview> => {
    const response = await api.get<AdminOverview>('/admin/overview');
    return response.data;
  },

  getAnalytics: async (): Promise<AnalyticsData> => {
    const response = await api.get<AnalyticsData>('/admin/analytics');
    return response.data;
  },

  listUsers: async (filter: AdminFilter, page: PaginationParams) => {
    const params = new URLSearchParams();
    params.append('page', String(page.page));
    params.append('page_size', String(page.page_size));
    if (filter.status) params.append('status', filter.status);
    if (filter.search) params.append('search', filter.search);
    if (filter.department) params.append('department', filter.department);
    if (filter.role_id) params.append('role_id', String(filter.role_id));
    const response = await api.get(`/admin/users/list?${params}`);
    return response.data;
  },

  toggleUser: async (userId: string, active: boolean) => {
    const response = await api.post('/admin/users/toggle', { user_id: userId, active });
    return response.data;
  },

  getPendingUsers: async (search: string, page: PaginationParams) => {
    const params = new URLSearchParams();
    params.append('page', String(page.page));
    params.append('page_size', String(page.page_size));
    if (search) params.append('search', search);
    const response = await api.get(`/admin/pending-users?${params}`);
    return response.data;
  },

  approveUser: async (pendingUserId: string, roleId: number = 2) => {
    const response = await api.post('/admin/pending-users/approve', { pending_user_id: pendingUserId, role_id: roleId });
    return response.data;
  },

  rejectUser: async (pendingUserId: string, reason?: string) => {
    const response = await api.post('/admin/pending-users/reject', { pending_user_id: pendingUserId, reason });
    return response.data;
  },

  listSuppliers: async (search: string, page: PaginationParams) => {
    const params = new URLSearchParams();
    params.append('page', String(page.page));
    params.append('page_size', String(page.page_size));
    if (search) params.append('search', search);
    const response = await api.get(`/admin/suppliers/list?${params}`);
    return response.data;
  },

  toggleSupplier: async (supplierId: string, active: boolean) => {
    const response = await api.post('/admin/suppliers/toggle', { supplier_id: supplierId, active });
    return response.data;
  },

  listNotices: async (activeOnly: boolean, page: PaginationParams) => {
    const params = new URLSearchParams();
    params.append('page', String(page.page));
    params.append('page_size', String(page.page_size));
    if (activeOnly) params.append('active', 'true');
    const response = await api.get(`/admin/notices/list?${params}`);
    return response.data;
  },

  listRequests: async (filter: AdminFilter, page: PaginationParams) => {
    const params = new URLSearchParams();
    params.append('page', String(page.page));
    params.append('page_size', String(page.page_size));
    if (filter.status) params.append('status', filter.status);
    if (filter.search) params.append('search', filter.search);
    const response = await api.get(`/admin/requests/list?${params}`);
    return response.data;
  },

  manageRequest: async (requestId: string, action: string, reason?: string) => {
    const response = await api.post('/admin/requests/manage', { request_id: requestId, action, reason });
    return response.data;
  },

  listIssues: async (filter: AdminFilter, page: PaginationParams) => {
    const params = new URLSearchParams();
    params.append('page', String(page.page));
    params.append('page_size', String(page.page_size));
    if (filter.status) params.append('status', filter.status);
    if (filter.search) params.append('search', filter.search);
    if (filter.overdue) params.append('overdue', 'true');
    const response = await api.get(`/admin/issues/list?${params}`);
    return response.data;
  },

  bulkAction: async (type: 'requests' | 'issues', ids: string[], action: string) => {
    const response = await api.post(`/admin/${type}/bulk?type=${type}`, { ids, action });
    return response.data;
  },
};