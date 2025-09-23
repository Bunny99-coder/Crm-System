"use client"

import { useState } from "react"

// ======================
// Types
// ======================
export interface LoginResponse {
  token: string
}

export interface User {
  id: number
  username: string
  email?: string
  role_id: number
}

// ======================
// Auth Manager (Singleton)
// ======================
class AuthManager {
  private static instance: AuthManager
  private token: string | null = null
  private user: User | null = null

  private constructor() {
    // Load from localStorage
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("jwt_token")
      if (this.token) {
        this.user = this.decodeUserFromToken(this.token)
      }
    }
  }

  public static getInstance(): AuthManager {
    if (!AuthManager.instance) {
      AuthManager.instance = new AuthManager()
    }
    return AuthManager.instance
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
    return this.token
  }

  public getUser(): User | null {
    return this.user
  }

  public isAuthenticated(): boolean {
    return this.token !== null && this.user !== null && !this.isTokenExpired()
  }

  public hasRole(roleId: number): boolean {
    return this.user?.role_id === roleId
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
    if (!this.token) return true
    try {
      const payload = JSON.parse(atob(this.token.split(".")[1]))
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
    if (this.token && !this.isTokenExpired()) {
      headers["Authorization"] = `Bearer ${this.token}`
    }
    return headers
  }
}

export const authManager = AuthManager.getInstance()

// ======================
// React Hook
// ======================
export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(authManager.isAuthenticated())
  const [user, setUser] = useState(authManager.getUser())

  const login = async (username: string, password: string): Promise<LoginResponse> => {
    try {
      const response = await fetch("http://localhost:8080/api/v1/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      })

      if (!response.ok) {
        throw new Error("Login failed")
      }

      const data: LoginResponse = await response.json()
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
    isAuthenticated,
    user,
    login,
    logout,
    checkAuth,
    hasRole: (roleId: number) => authManager.hasRole(roleId),
  }
}
