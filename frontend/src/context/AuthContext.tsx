import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authService } from '../services/auth';

interface AuthContextType {
  user: { user_id: string; username: string; role_id: number } | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<{ user_id: string; username: string; role_id: number } | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const claims = authService.getClaims();
    if (claims) {
      setUser(claims);
    }
    setIsLoading(false);
  }, []);

  const login = async (username: string, password: string) => {
    const response = await authService.login({ username, password });
    localStorage.setItem('access_token', response.access_token);
    localStorage.setItem('refresh_token', response.refresh_token);
    setUser({
      user_id: response.user.id,
      username: response.user.username,
      role_id: response.user.role_id,
    });
  };

  const logout = async () => {
    await authService.logout();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
}