import api from './api';
import { RequestListResponse, CreateRequestRequest, ItemRequest } from '../types';

export const requestService = {
  list: async (filter?: { status?: string; user_id?: string }, page = 1, pageSize = 20): Promise<RequestListResponse> => {
    const params = new URLSearchParams();
    params.append('page', String(page));
    params.append('page_size', String(pageSize));
    if (filter?.status) params.append('status', filter.status);
    if (filter?.user_id) params.append('user_id', filter.user_id);
    const response = await api.get<RequestListResponse>(`/requests?${params}`);
    return response.data;
  },

  get: async (id: string): Promise<ItemRequest> => {
    const response = await api.get<ItemRequest>(`/requests/get?id=${id}`);
    return response.data;
  },

  create: async (data: CreateRequestRequest): Promise<ItemRequest> => {
    const response = await api.post<ItemRequest>('/requests/create', data);
    return response.data;
  },

  cancel: async (id: string): Promise<ItemRequest> => {
    const response = await api.post<ItemRequest>(`/requests/cancel?id=${id}`);
    return response.data;
  },

  approve: async (id: string, notes?: string): Promise<ItemRequest> => {
    const response = await api.post<ItemRequest>(`/requests/approve?id=${id}`, { approved: true, notes });
    return response.data;
  },

  reject: async (id: string, reason: string): Promise<ItemRequest> => {
    const response = await api.post<ItemRequest>(`/requests/reject?id=${id}`, { approved: false, rejection_reason: reason });
    return response.data;
  },

  fulfill: async (id: string): Promise<ItemRequest> => {
    const response = await api.post<ItemRequest>(`/requests/fulfill?id=${id}`);
    return response.data;
  },

  getPending: async (page = 1, pageSize = 20): Promise<RequestListResponse> => {
    const response = await api.get<RequestListResponse>(`/requests/pending?page=${page}&page_size=${pageSize}`);
    return response.data;
  },

  getMyRequests: async (page = 1, pageSize = 20): Promise<RequestListResponse> => {
    const response = await api.get<RequestListResponse>(`/requests?page=${page}&page_size=${pageSize}`);
    return response.data;
  },

  getByItem: async (itemId: string, page = 1, pageSize = 20): Promise<RequestListResponse> => {
    const response = await api.get<RequestListResponse>(`/requests?item_id=${itemId}&page=${page}&page_size=${pageSize}`);
    return response.data;
  },
};