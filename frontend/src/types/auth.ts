export type Role = 1 | 2 | 3;

export const RoleNames: Record<Role, string> = {
  1: 'admin',
  2: 'authorized_user',
  3: 'viewer',
};

export interface User {
  id: string;
  username: string;
  full_name: string;
  email?: string;
  department?: string;
  employee_id?: string;
  phone?: string;
  role_id: Role;
  role?: Role;
  is_active: boolean;
  created_at: string;
}

export interface AuthUser extends User {
  access_token: string;
  refresh_token: string;
  expires_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_at: string;
  user: User;
}

export interface Claims {
  user_id: string;
  username: string;
  role_id: Role;
}