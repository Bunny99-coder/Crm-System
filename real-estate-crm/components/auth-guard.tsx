"use client"

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/lib/auth'

interface AuthGuardProps {
  children: React.ReactNode
}

export default function AuthGuard({ children }: AuthGuardProps) {
  const router = useRouter()
  const { loading, isAuthenticated } = useAuth()

  useEffect(() => {
    // Only redirect if not loading and not authenticated
    if (!loading && !isAuthenticated) {
      router.push('/login')
    }
  }, [loading, isAuthenticated, router])

  if (loading) {
    return <div className="flex items-center justify-center min-h-screen">Loading authentication...</div>
  }

  if (!isAuthenticated) {
    // This state should ideally be caught by the useEffect and redirected
    // but as a fallback, if isAuthenticated is false after loading, show a message
    return <div className="flex items-center justify-center min-h-screen">Not authenticated. Redirecting...</div>
  }

  return <>{children}</>
}
