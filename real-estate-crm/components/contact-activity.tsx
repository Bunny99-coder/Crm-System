"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Clock, FileText, MessageSquare, Phone, Mail } from "lucide-react"
import { api, type Activity } from "@/lib/api"

interface ContactActivityProps {
  contactId: number
}

export function ContactActivity({ contactId }: ContactActivityProps) {
  const [activities, setActivities] = useState<Activity[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    loadActivity()
  }, [contactId])

  const loadActivity = async () => {
    try {
      setIsLoading(true)
      const activityData = await api.getActivities("contact", contactId)
      setActivities(activityData)
    } catch (err) {
      console.error("Failed to load activity:", err)
    } finally {
      setIsLoading(false)
    }
  }

  const getActivityIcon = (type: string) => {
    switch (type) {
      case "note":
        return <FileText className="h-4 w-4 text-blue-600" />
      case "comm_log":
        return <MessageSquare className="h-4 w-4 text-green-600" />
      case "task":
        return <Phone className="h-4 w-4 text-orange-600" />
      case "event":
        return <Mail className="h-4 w-4 text-purple-600" />
      default:
        return <Clock className="h-4 w-4 text-gray-600" />
    }
  }

  const getActivityBadge = (type: string) => {
    const colors = {
      note: "bg-blue-50 text-blue-700 border-blue-200",
      comm_log: "bg-green-50 text-green-700 border-green-200",
      task: "bg-orange-50 text-orange-700 border-orange-200",
      event: "bg-purple-50 text-purple-700 border-purple-200",
    }

    const labels = {
      note: "Note",
      comm_log: "Communication",
      task: "Task",
      event: "Event",
    }

    return (
      <Badge
        variant="outline"
        className={colors[type as keyof typeof colors] || "bg-gray-50 text-gray-700 border-gray-200"}
      >
        {labels[type as keyof typeof labels] || type}
      </Badge>
    )
  }

  const formatTimeAgo = (dateString: string) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60))

    if (diffInHours < 1) return "Just now"
    if (diffInHours < 24) return `${diffInHours}h ago`
    if (diffInHours < 168) return `${Math.floor(diffInHours / 24)}d ago`
    return date.toLocaleDateString()
  }

  if (isLoading) {
    return <div className="text-center py-8">Loading activity...</div>
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Activity Feed</CardTitle>
        <CardDescription>Recent activity and updates for this contact</CardDescription>
      </CardHeader>
      <CardContent>
        {activities.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-muted-foreground">
              No activity yet. Activity will appear here as you work with this contact.
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {activities.map((activity) => (
              <div key={activity.id} className="flex gap-4 p-4 rounded-lg border bg-card">
                <div className="flex-shrink-0">
                  <Avatar className="h-8 w-8">
                    <AvatarFallback className="text-xs">{getActivityIcon(activity.activity_type)}</AvatarFallback>
                  </Avatar>
                </div>
                <div className="flex-1 space-y-2">
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-medium">{activity.title}</p>
                    {getActivityBadge(activity.activity_type)}
                  </div>
                  {activity.description && <p className="text-sm text-muted-foreground">{activity.description}</p>}
                  <div className="flex items-center gap-2 text-xs text-muted-foreground">
                    <span>User {activity.created_by}</span>
                    <span>•</span>
                    <span>{formatTimeAgo(activity.created_at || new Date().toISOString())}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
