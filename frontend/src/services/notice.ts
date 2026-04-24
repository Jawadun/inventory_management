import api from './api';
import { Notice, CreateNoticeRequest, UpdateNoticeRequest } from '../types';

export const noticeService = {
  list: async (activeOnly = false): Promise<Notice[]> => {
    const response = await api.get<Notice[]>(`/notices?active=${activeOnly}`);
    return response.data;
  },

  get: async (id: string): Promise<Notice> => {
    const response = await api.get<Notice>(`/notices/get?id=${id}`);
    return response.data;
  },

  create: async (data: CreateNoticeRequest): Promise<Notice> => {
    const response = await api.post<Notice>('/admin/notices/create', data);
    return response.data;
  },

  update: async (id: string, data: UpdateNoticeRequest): Promise<Notice> => {
    const response = await api.put<Notice>(`/admin/notices/update?id=${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/admin/notices/delete?id=${id}`);
  },
};