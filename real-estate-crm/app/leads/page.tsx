"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import DashboardSidebar from "@/components/dashboard-sidebar"
import DashboardHeader from "@/components/dashboard-header"
import { LeadsManagement } from "@/components/leads-management"
import { useAuth, ROLE_SALES_AGENT, ROLE_RECEPTION } from "@/lib/auth"

export default function LeadsPage() {
  const { isAuthenticated, user, loading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push("/login")
    } else if (!loading && user?.role_id !== ROLE_SALES_AGENT) {
      // Redirect non-sales agents to their appropriate dashboard or a general dashboard
      if (user?.role_id === ROLE_RECEPTION) {
        router.push("/") // Receptionists go to the main dashboard
      } else {
        router.push("/") // Fallback for other roles
      }
    }
  }, [isAuthenticated, user, router, loading])

  if (loading || !isAuthenticated || user?.role_id !== ROLE_SALES_AGENT) {
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
          <LeadsManagement />
        </main>
      </div>
    </div>
  )
}
