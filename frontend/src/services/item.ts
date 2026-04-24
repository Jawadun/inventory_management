import api from './api';
import { ItemListResponse, ItemFilter, CreateItemRequest, UpdateItemRequest, Item, ItemHistory, Category, Supplier } from '../types';

export const itemService = {
  list: async (filter?: ItemFilter, page = 1, pageSize = 20): Promise<ItemListResponse> => {
    const params = new URLSearchParams();
    params.append('page', String(page));
    params.append('page_size', String(pageSize));
    if (filter?.category_id) params.append('category_id', filter.category_id);
    if (filter?.supplier_id) params.append('supplier_id', filter.supplier_id);
    if (filter?.status) params.append('status', filter.status);
    if (filter?.search) params.append('search', filter.search);
    if (filter?.low_stock) params.append('low_stock', 'true');
    const response = await api.get<ItemListResponse>(`/items?${params}`);
    return response.data;
  },

  get: async (id: string): Promise<Item> => {
    const response = await api.get<Item>(`/items/get?id=${id}`);
    return response.data;
  },

  create: async (data: CreateItemRequest): Promise<Item> => {
    const response = await api.post<Item>('/items/create', data);
    return response.data;
  },

  update: async (id: string, data: UpdateItemRequest): Promise<Item> => {
    const response = await api.put<Item>(`/items/update?id=${id}`, data);
    return response.data;
  },

  adjustQuantity: async (id: string, data: { quantity_change: number; reason: string }): Promise<Item> => {
    const response = await api.patch<Item>(`/items/adjust?id=${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/items/delete?id=${id}`);
  },

  archive: async (id: string): Promise<void> => {
    await api.post(`/items/archive?id=${id}`);
  },

  restore: async (id: string): Promise<void> => {
    await api.post(`/items/restore?id=${id}`);
  },

  getHistory: async (id: string): Promise<ItemHistory[]> => {
    const response = await api.get<ItemHistory[]>(`/items/history?id=${id}`);
    return response.data;
  },

  listCategories: async (): Promise<Category[]> => {
    const response = await api.get<Category[]>('/categories');
    return response.data;
  },

  createCategory: async (data: { name: string; description?: string }): Promise<Category> => {
    const response = await api.post<Category>('/categories/create', data);
    return response.data;
  },

  updateCategory: async (id: string, data: { name?: string; description?: string }): Promise<Category> => {
    const response = await api.put<Category>(`/categories/update?id=${id}`, data);
    return response.data;
  },

  deleteCategory: async (id: string, moveToId?: string): Promise<void> => {
    const params = new URLSearchParams();
    if (moveToId) params.append('move_to_category', moveToId);
    await api.delete(`/categories/delete?id=${id}?${params}`);
  },

  listSuppliers: async (): Promise<Supplier[]> => {
    const response = await api.get<Supplier[]>('/suppliers');
    return response.data;
  },

  createSupplier: async (data: { name: string; contact_person?: string; phone?: string; email?: string; address?: string; notes?: string }): Promise<Supplier> => {
    const response = await api.post<Supplier>('/suppliers/create', data);
    return response.data;
  },

  updateSupplier: async (id: string, data: { name?: string; contact_person?: string; phone?: string; email?: string; address?: string; notes?: string }): Promise<Supplier> => {
    const response = await api.put<Supplier>(`/suppliers/update?id=${id}`, data);
    return response.data;
  },

  deleteSupplier: async (id: string, moveToId?: string): Promise<void> => {
    const params = new URLSearchParams();
    if (moveToId) params.append('move_to_supplier', moveToId);
    await api.delete(`/suppliers/delete?id=${id}?${params}`);
  },
};