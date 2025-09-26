"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import DashboardSidebar from "@/components/dashboard-sidebar"
import DashboardHeader from "@/components/dashboard-header"
import { DashboardStats } from "@/components/dashboard-stats"
import { ContactsTable } from "@/components/contacts-table"
import { useAuth, ROLE_SALES_AGENT, ROLE_RECEPTION } from "@/lib/auth"

export default function DashboardPage() {
  const { user, hasRole, isAuthenticated, loading } = useAuth()
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
        <main className="flex-1 overflow-y-auto p-6 space-y-6">
          {/* Welcome Section */}
          <div className="space-y-2">
            <h1 className="text-3xl font-bold text-balance text-foreground">Welcome, {user?.username}!</h1>
            <p className="text-muted-foreground text-pretty">
              Here's what's happening with your real estate business today.
            </p>
          </div>

          {hasRole(ROLE_RECEPTION) && (
            <>
              {/* Stats */}
              <DashboardStats />

              {/* Recent Contacts */}
              <ContactsTable />
            </>
          )}

          {hasRole(ROLE_SALES_AGENT) && (
            <div className="rounded-lg border bg-card text-card-foreground shadow-sm p-6">
              <h2 className="text-2xl font-semibold mb-4">Your Sales Agent Dashboard</h2>
              <p className="text-muted-foreground">Here you will see reports and activities relevant to your assigned leads and deals.</p>
              {/* Placeholder for Sales Agent specific reports/widgets */}
              <div className="mt-4 p-4 border rounded-md border-dashed text-muted-foreground">
                <p>Sales Agent specific content coming soon...</p>
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
