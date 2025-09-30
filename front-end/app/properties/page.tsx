"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import DashboardSidebar from "@/components/dashboard-sidebar"
import DashboardHeader from "@/components/dashboard-header"
import { PropertiesManagement } from "@/components/properties-management"
import { useAuth } from "@/lib/auth"

export default function PropertiesPage() {
  const { isAuthenticated, loading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push("/login")
    }
  }, [isAuthenticated, router, loading])

  if (loading || !isAuthenticated) {
    return <div className="flex items-center justify-center h-screen">Loading...</div> // Or a loading spinner
  }

  return (
    <div className="flex h-screen bg-background">
      {/* Sidebar */}
      <DashboardSidebar />

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <DashboardHeader />

        {/* Content */}
        <main className="flex-1 overflow-y-auto p-6">
          <PropertiesManagement />
        </main>
      </div>
    </div>
  )
}
