"use client"

import type React from "react"

import { createContext, useContext, useEffect, useState } from "react"
import { authManager, type User } from "@/lib/auth"
import { useRouter } from "next/navigation"

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  useEffect(() => {
    // Check authentication status on mount
    const checkAuth = () => {
      const authenticated = authManager.isAuthenticated() && !authManager.isTokenExpired()
      const currentUser = authManager.getUser()

      if (!authenticated && authManager.getToken()) {
        // Token exists but is expired, clear it
        authManager.clearAuth()
      }

      setIsAuthenticated(authenticated)
      setUser(currentUser)
      setIsLoading(false)
    }

    checkAuth()
  }, [])

  const logout = () => {
    authManager.clearAuth()
    setIsAuthenticated(false)
    setUser(null)
    router.push("/login")
  }

  return <AuthContext.Provider value={{ user, isAuthenticated, isLoading, logout }}>{children}</AuthContext.Provider>
}

export function useAuthContext() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuthContext must be used within an AuthProvider")
  }
  return context
}
