"use client"

import { useState, useEffect } from "react"
import { useRouter } from 'next/navigation';
import { api } from "./api"

// ======================
// Role Constants
// ======================
export const ROLE_SALES_AGENT = 1
export const ROLE_RECEPTION = 2 // Manager role

// ======================
// Types
// ======================
export interface User {
  id?: number
  username: string
  email: string
  role_id: number
}

export interface LoginResponse {
  token: string
  user: User
}

// ======================
// Auth Manager (Singleton)
// ======================
class AuthManager {
  private static instance: AuthManager
  private token: string | null = null
  private user: User | null = null

  private constructor() {
    // Do not load from localStorage here to avoid SSR issues
  }

  public static getInstance(): AuthManager {
    if (!AuthManager.instance) {
      AuthManager.instance = new AuthManager()
    }
    return AuthManager.instance
  }

  public initFromLocalStorage(): void {
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("jwt_token")
      if (this.token) {
        this.user = this.decodeUserFromToken(this.token)
      }
    }
  }

  public setAuth(token: string): void {
    this.token = token
    this.user = this.decodeUserFromToken(token)

    if (typeof window !== "undefined") {
      localStorage.setItem("jwt_token", token)
    }
  }

  public clearAuth(): void {
    this.token = null
    this.user = null

    if (typeof window !== "undefined") {
      localStorage.removeItem("jwt_token")
    }
  }

  public getToken(): string | null {
    if (!this.token && typeof window !== "undefined") {
      this.token = localStorage.getItem("jwt_token")
    }
    return this.token
  }

  public getUser(): User | null {
    if (!this.user && this.token) { // Only decode if token exists
      this.user = this.decodeUserFromToken(this.token)
    }
    return this.user
  }

  public isAuthenticated(): boolean {
    return this.getToken() !== null && this.getUser() !== null && !this.isTokenExpired()
  }

  public hasRole(roleId: number): boolean {
    return this.getUser()?.role_id === roleId
  }

  // ======================
  // Helpers
  // ======================
  private decodeUserFromToken(token: string): User | null {
    try {
      const payload = JSON.parse(atob(token.split(".")[1]))
      return {
        id: payload.user_id,
        username: payload.username,
        role_id: payload.role_id,
        email: payload.email || "",
      }
    } catch (error) {
      console.error("Failed to decode JWT:", error)
      return null
    }
  }

  public isTokenExpired(): boolean {
    if (!this.getToken()) return true
    try {
      const payload = JSON.parse(atob(this.getToken()!.split(".")[1]))
      const currentTime = Math.floor(Date.now() / 1000)
      return payload.exp && payload.exp < currentTime
    } catch {
      return true
    }
  }

  public getAuthHeaders(): Record<string, string> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    }
    if (this.getToken() && !this.isTokenExpired()) {
      headers["Authorization"] = `Bearer ${this.getToken()}`
    }
    return headers
  }
}

export const authManager = AuthManager.getInstance()
// Call initFromLocalStorage immediately after getting the instance
// This ensures localStorage is checked as early as possible on the client side.
if (typeof window !== "undefined") {
  authManager.initFromLocalStorage()
  ;(window as any).authManager = authManager; // Expose for debugging
}

// ======================
// React Hook
// ======================

export function useAuth() {
  const [loading, setLoading] = useState(true)
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false)
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    authManager.initFromLocalStorage()
    setIsAuthenticated(authManager.isAuthenticated())
    setUser(authManager.getUser())
    setLoading(false)
  }, [])

  const login = async (username: string, password: string): Promise<LoginResponse> => {
    try {
      const data = await api.login(username, password)
      authManager.setAuth(data.token)
      setIsAuthenticated(true)
      setUser(authManager.getUser())
      return data
    } catch (error) {
      console.error("Login error:", error)
      throw error
    } 
  }

  const logout = () => {
    authManager.clearAuth()
    setIsAuthenticated(false)
    setUser(null)
    
    router.push("/login")
  }


  const checkAuth = () => {
    const authenticated = authManager.isAuthenticated()
    if (!authenticated && isAuthenticated) {
      logout()
    }
    setIsAuthenticated(authenticated)
    setUser(authManager.getUser())
  }

  return {
    loading,
    isAuthenticated,
    user,
    login,
    logout,
    checkAuth,
    hasRole: (roleId: number) => authManager.hasRole(roleId),
  }
}
