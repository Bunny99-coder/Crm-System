// Replace the entire contents of src/context/AuthContext.tsx
import React, { createContext, useState, useContext, useEffect } from 'react';
import type { ReactNode } from 'react';
import { login as apiLogin } from '../services/apiService';
import type { LoginCredentials } from '../services/apiService';
import { useNavigate } from 'react-router-dom';
import { getRoleFromToken } from '../util/auth';

interface AuthContextType {
  token: string | null;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  userRole: number | null;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('authToken'));
  const [userRole, setUserRole] = useState<number | null>(() => getRoleFromToken(localStorage.getItem('authToken')));
  const navigate = useNavigate();

  const login = async (credentials: LoginCredentials) => {
    const response = await apiLogin(credentials);
    localStorage.setItem('authToken', response.token);
    const role = getRoleFromToken(response.token);
    setToken(response.token);
    setUserRole(role);
    navigate('/');
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    setToken(null);
    setUserRole(null);
    navigate('/');
  };

  const isAuthenticated = !!token;

  return (
    <AuthContext.Provider value={{ token, login, logout, isAuthenticated, userRole }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};