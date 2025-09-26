import { useState, useEffect, useContext, createContext } from 'react';

export const ROLE_RECEPTION = 'reception';
export const ROLE_SALES_AGENT = 'sales_agent';

interface AuthContextType {
  user: { id: string; name: string; roles: string[] } | null;
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
    // Simulate API call
    return new Promise<void>((resolve) => {
      setTimeout(() => {
        const mockUser = { id: '1', name: username, roles: [ROLE_RECEPTION] }; // Default to reception for testing
        localStorage.setItem('user', JSON.stringify(mockUser));
        setUser(mockUser);
        setIsLoading(false);
        resolve();
      }, 1000);
    });
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
