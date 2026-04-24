export type Role = 1 | 2 | 3;

export const RoleNames: Record<Role, string> = {
  1: 'admin',
  2: 'authorized_user',
  3: 'viewer',
};

export function hasRole(userRole: Role | undefined, requiredRole: Role): boolean {
  if (!userRole) return false;
  if (requiredRole === 1) return userRole === 1;
  if (requiredRole === 2) return userRole === 1 || userRole === 2;
  return true;
}

export function isAdmin(role?: Role): boolean {
  return role === 1;
}

export function isUser(role?: Role): boolean {
  return role === 1 || role === 2;
}