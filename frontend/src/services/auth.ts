import api from './api';
import { LoginRequest, LoginResponse, Claims } from '../types';

export const authService = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data);
    return response.data;
  },

  register: async (data: { username: string; password: string; full_name: string; email: string; department: string }): Promise<void> => {
    await api.post('/auth/register', data);
  },

  logout: async (): Promise<void> => {
    const refreshToken = localStorage.getItem('refresh_token');
    if (refreshToken) {
      try {
        await api.post('/auth/logout', { refresh_token: refreshToken });
      } catch {
        // ignore
      }
    }
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  },

  getClaims: (): Claims | null => {
    const token = localStorage.getItem('access_token');
    if (!token) return null;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return {
        user_id: payload.user_id,
        username: payload.username,
        role_id: payload.role_id,
      };
    } catch {
      return null;
    }
  },

  isAuthenticated: (): boolean => {
    const claims = authService.getClaims();
    if (!claims) return false;
    return Date.now() < claims.exp * 1000;
  },
};