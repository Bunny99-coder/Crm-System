import { useState, useEffect, useContext, createContext } from 'react';
import { api } from '../lib/api'; // Import the API client

export const ROLE_RECEPTION = 'reception';
export const ROLE_SALES_AGENT = 'sales_agent';

// Define role IDs based on your backend configuration
export const ROLE_ID_RECEPTION = 1; // Assuming 1 for Receptionist
export const ROLE_ID_SALES_AGENT = 2; // Assuming 2 for Sales Agent

interface AuthContextType {
  user: { id: number; name: string; roles: string[]; role_id: number } | null; // Add role_id
  hasRole: (role: string) => boolean;
  isLoading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<AuthContextType['user']>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulate loading user from local storage or API
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
    setIsLoading(false);
  }, []);

  const hasRole = (role: string) => {
    return user?.roles.includes(role) || false;
  };

  const login = async (username: string, password: string) => {
    setIsLoading(true);
    try {
      const response = await api.login(username, password);
      const userRoles: string[] = [];
      let roleId: number = 0;

      if (response.user.role_id === ROLE_ID_RECEPTION) {
        userRoles.push(ROLE_RECEPTION);
        roleId = ROLE_ID_RECEPTION;
      } else if (response.user.role_id === ROLE_ID_SALES_AGENT) {
        userRoles.push(ROLE_SALES_AGENT);
        roleId = ROLE_ID_SALES_AGENT;
      }

      const authenticatedUser = {
        id: response.user.id,
        name: response.user.username,
        roles: userRoles,
        role_id: roleId,
      };
      localStorage.setItem('user', JSON.stringify(authenticatedUser));
      setUser(authenticatedUser);
    } catch (error) {
      console.error("Login failed:", error);
      logout(); // Clear any partial auth state
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem('user');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, hasRole, isLoading, login, logout }}>
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
