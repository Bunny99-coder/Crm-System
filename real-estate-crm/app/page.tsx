import { DashboardSidebar } from "@/components/dashboard-sidebar"
import { DashboardHeader } from "@/components/dashboard-header"
import { DashboardStats } from "@/components/dashboard-stats"
import { ContactsTable } from "@/components/contacts-table"

export default function DashboardPage() {
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
            <h1 className="text-3xl font-bold text-balance text-foreground">Welcome</h1>
            <p className="text-muted-foreground text-pretty">
              Here's what's happening with your real estate business today.
            </p>
          </div>

          {/* Stats */}
          <DashboardStats />

          {/* Recent Contacts */}
          <ContactsTable />
        </main>
      </div>
    </div>
  )
}
