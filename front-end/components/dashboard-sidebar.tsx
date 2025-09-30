"use client"

import { useState } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Home, Users, Building, UserCheck, Handshake, CheckSquare, BarChart3, Settings, Menu, X } from "lucide-react"
import { useAuth, ROLE_SALES_AGENT, ROLE_RECEPTION } from "@/lib/auth" // Import useAuth and role constants

interface NavItem {
  name: string
  href: string
  icon: React.ElementType
  roles?: number[] // Optional: roles that can see this item
}

const navigation: NavItem[] = [
  { name: "Dashboard", href: "/", icon: Home },
  { name: "Contacts", href: "/contacts", icon: Users },
  { name: "Properties", href: "/properties", icon: Building, roles: [ROLE_RECEPTION, ROLE_SALES_AGENT] }, // Both can see
  { name: "Leads", href: "/leads", icon: UserCheck, roles: [ROLE_RECEPTION, ROLE_SALES_AGENT] }, // Both can see
  { name: "Deals", href: "/deals", icon: Handshake },
  { name: "Tasks", href: "/tasks", icon: CheckSquare },
  { name: "Reports", href: "/reports", icon: BarChart3, roles: [ROLE_RECEPTION] }, // Only manager can see
  { name: "Settings", href: "/settings", icon: Settings, roles: [ROLE_RECEPTION] }, // Only manager can see
]

export default function DashboardSidebar() {
  const pathname = usePathname()
  const [isCollapsed, setIsCollapsed] = useState(false)
  const { user, hasRole } = useAuth() // Use the useAuth hook

  const filteredNavigation = navigation.filter(item => {
    if (!item.roles) return true; // If no roles specified, visible to all authenticated users
    return item.roles.some(role => hasRole(role));
  });

  return (
    <div
      className={cn(
        "flex flex-col h-screen bg-sidebar border-r border-sidebar-border transition-all duration-300",
        isCollapsed ? "w-16" : "w-64",
      )}
    >
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-sidebar-border">
        {!isCollapsed && <h1 className="text-xl font-bold text-sidebar-foreground">Real Estate CRM</h1>}
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setIsCollapsed(!isCollapsed)}
          className="text-sidebar-foreground hover:bg-sidebar-accent"
        >
          {isCollapsed ? <Menu className="h-4 w-4" /> : <X className="h-4 w-4" />}
        </Button>
      </div>

      {/* Navigation */}
      <ScrollArea className="flex-1 px-3 py-4">
        <nav className="space-y-2">
          {filteredNavigation.map((item) => {
            const isActive = pathname === item.href
            return (
              <Link key={item.name} href={item.href}>
                <Button
                  variant={isActive ? "default" : "ghost"}
                  className={cn(
                    "w-full justify-start gap-3 text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
                    isActive && "bg-primary text-primary-foreground hover:bg-primary/90",
                    isCollapsed && "justify-center px-2",
                  )}
                >
                  <item.icon className="h-4 w-4 flex-shrink-0" />
                  {!isCollapsed && <span>{item.name}</span>}
                </Button>
              </Link>
            )
          })}
        </nav>
      </ScrollArea>

      {/* User Profile */}
      <div className="p-4 border-t border-sidebar-border">
        <div className={cn("flex items-center gap-3", isCollapsed && "justify-center")}>
          <div className="w-8 h-8 bg-primary rounded-full flex items-center justify-center">
            <span className="text-sm font-medium text-primary-foreground">{user?.username.substring(0, 2).toUpperCase()}</span>
          </div>
          {!isCollapsed && (
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-sidebar-foreground truncate">{user?.username}</p>
              <p className="text-xs text-sidebar-foreground/60 truncate">{user?.role_id === ROLE_RECEPTION ? "Manager" : "Sales Agent"}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
