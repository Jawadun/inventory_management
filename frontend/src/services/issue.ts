import api from './api';
import { IssueListResponse, IssueRecord, CreateIssueRequest, CreateReturnRequest } from '../types';

export const issueService = {
  list: async (filter?: { status?: string; recipient_id?: string }, page = 1, pageSize = 20): Promise<IssueListResponse> => {
    const params = new URLSearchParams();
    params.append('page', String(page));
    params.append('page_size', String(pageSize));
    if (filter?.status) params.append('status', filter.status);
    if (filter?.recipient_id) params.append('recipient_id', filter.recipient_id);
    const response = await api.get<IssueListResponse>(`/issues?${params}`);
    return response.data;
  },

  get: async (id: string): Promise<IssueRecord> => {
    const response = await api.get<IssueRecord>(`/issues/get?id=${id}`);
    return response.data;
  },

  create: async (data: CreateIssueRequest): Promise<IssueRecord> => {
    const response = await api.post<IssueRecord>('/issues/create', data);
    return response.data;
  },

  approve: async (id: string): Promise<IssueRecord> => {
    const response = await api.post<IssueRecord>(`/issues/approve?id=${id}`);
    return response.data;
  },

  reject: async (id: string, reason: string): Promise<IssueRecord> => {
    const response = await api.post<IssueRecord>(`/issues/reject?id=${id}`, { notes: reason });
    return response.data;
  },

  return: async (id: string, data: CreateReturnRequest): Promise<IssueRecord> => {
    const response = await api.post<IssueRecord>(`/issues/return?id=${id}`, data);
    return response.data;
  },

  getOverdue: async (page = 1, pageSize = 20): Promise<IssueListResponse> => {
    const response = await api.get<IssueListResponse>(`/issues/overdue?page=${page}&page_size=${pageSize}`);
    return response.data;
  },

  getMyIssues: async (page = 1, pageSize = 20): Promise<IssueListResponse> => {
    const response = await api.get<IssueListResponse>(`/issues?page=${page}&page_size=${pageSize}`);
    return response.data;
  },

  getActive: async (page = 1, pageSize = 20): Promise<IssueListResponse> => {
    const response = await api.get<IssueListResponse>(`/issues?status=issued&page=${page}&page_size=${pageSize}`);
    return response.data;
  },

  getByItem: async (itemId: string, page = 1, pageSize = 20): Promise<IssueListResponse> => {
    const response = await api.get<IssueListResponse>(`/issues?item_id=${itemId}&page=${page}&page_size=${pageSize}`);
    return response.data;
  },
};