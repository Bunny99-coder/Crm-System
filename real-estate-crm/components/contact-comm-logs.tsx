"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Plus, Phone, Mail, MessageSquare, Calendar, User } from "lucide-react"
import { api } from "@/lib/api"

interface CommLog {
  id: number
  communication_type: "Email" | "Phone" | "Meeting" | "SMS" | "Other"
  subject: string
  content: string
  communication_date: string
  direction: "Inbound" | "Outbound"
  created_by: number
  created_at: string
}

interface ContactCommLogsProps {
  contactId: number
}

export function ContactCommLogs({ contactId }: ContactCommLogsProps) {
  const [commLogs, setCommLogs] = useState<CommLog[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isCreating, setIsCreating] = useState(false)
  const [formData, setFormData] = useState({
    communication_type: "Email" as CommLog["communication_type"],
    subject: "",
    content: "",
    communication_date: "",
    direction: "Outbound" as CommLog["direction"],
  })

  useEffect(() => {
    loadCommLogs()
  }, [contactId])

  const loadCommLogs = async () => {
    try {
      setIsLoading(true)
      const commLogsData = await api.get(`/contacts/${contactId}/communications`)
      setCommLogs(commLogsData)
    } catch (err) {
      console.error("Failed to load communication logs:", err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateCommLog = async () => {
    if (!formData.subject.trim() || !formData.communication_date) return

    try {
      await api.post(`/contacts/${contactId}/communications`, formData)
      setFormData({
        communication_type: "Email",
        subject: "",
        content: "",
        communication_date: "",
        direction: "Outbound",
      })
      setIsCreating(false)
      loadCommLogs()
    } catch (err) {
      console.error("Failed to create communication log:", err)
    }
  }

  const getCommTypeIcon = (type: string) => {
    switch (type) {
      case "Email":
        return <Mail className="h-4 w-4" />
      case "Phone":
        return <Phone className="h-4 w-4" />
      case "SMS":
        return <MessageSquare className="h-4 w-4" />
      case "Meeting":
        return <Calendar className="h-4 w-4" />
      default:
        return <MessageSquare className="h-4 w-4" />
    }
  }

  const getCommTypeBadge = (type: string) => {
    const colors = {
      Email: "bg-blue-50 text-blue-700 border-blue-200",
      Phone: "bg-green-50 text-green-700 border-green-200",
      Meeting: "bg-purple-50 text-purple-700 border-purple-200",
      SMS: "bg-yellow-50 text-yellow-700 border-yellow-200",
      Other: "bg-gray-50 text-gray-700 border-gray-200",
    }
    return (
      <Badge variant="outline" className={colors[type as keyof typeof colors] || colors.Other}>
        {type}
      </Badge>
    )
  }

  const getDirectionBadge = (direction: string) => {
    return (
      <Badge
        variant="outline"
        className={
          direction === "Inbound"
            ? "bg-orange-50 text-orange-700 border-orange-200"
            : "bg-cyan-50 text-cyan-700 border-cyan-200"
        }
      >
        {direction}
      </Badge>
    )
  }

  if (isLoading) {
    return <div className="text-center py-8">Loading communication logs...</div>
  }

  return (
    <div className="space-y-4">
      {/* Create Communication Log */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Communication Logs</CardTitle>
              <CardDescription>Track all communications with this contact</CardDescription>
            </div>
            {!isCreating && (
              <Button onClick={() => setIsCreating(true)} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                <Plus className="mr-2 h-4 w-4" />
                Log Communication
              </Button>
            )}
          </div>
        </CardHeader>
        {isCreating && (
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-3">
              <Select
                value={formData.communication_type}
                onValueChange={(value: CommLog["communication_type"]) =>
                  setFormData((prev) => ({ ...prev, communication_type: value }))
                }
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Email">Email</SelectItem>
                  <SelectItem value="Phone">Phone</SelectItem>
                  <SelectItem value="Meeting">Meeting</SelectItem>
                  <SelectItem value="SMS">SMS</SelectItem>
                  <SelectItem value="Other">Other</SelectItem>
                </SelectContent>
              </Select>
              <Select
                value={formData.direction}
                onValueChange={(value: CommLog["direction"]) => setFormData((prev) => ({ ...prev, direction: value }))}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Inbound">Inbound</SelectItem>
                  <SelectItem value="Outbound">Outbound</SelectItem>
                </SelectContent>
              </Select>
              <Input
                type="datetime-local"
                value={formData.communication_date}
                onChange={(e) => setFormData((prev) => ({ ...prev, communication_date: e.target.value }))}
              />
            </div>
            <Input
              placeholder="Subject/Title"
              value={formData.subject}
              onChange={(e) => setFormData((prev) => ({ ...prev, subject: e.target.value }))}
            />
            <Textarea
              placeholder="Communication details/notes"
              value={formData.content}
              onChange={(e) => setFormData((prev) => ({ ...prev, content: e.target.value }))}
              rows={3}
            />
            <div className="flex gap-2">
              <Button onClick={handleCreateCommLog} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                Log Communication
              </Button>
              <Button
                onClick={() => {
                  setIsCreating(false)
                  setFormData({
                    communication_type: "Email",
                    subject: "",
                    content: "",
                    communication_date: "",
                    direction: "Outbound",
                  })
                }}
                variant="outline"
                size="sm"
              >
                Cancel
              </Button>
            </div>
          </CardContent>
        )}
      </Card>

      {/* Communication Logs List */}
      <div className="space-y-4">
        {commLogs.length === 0 ? (
          <Card>
            <CardContent className="text-center py-8">
              <p className="text-muted-foreground">
                No communication logs yet. Log your first communication to get started.
              </p>
            </CardContent>
          </Card>
        ) : (
          commLogs.map((log) => (
            <Card key={log.id}>
              <CardContent className="pt-6">
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      {getCommTypeIcon(log.communication_type)}
                      <h4 className="font-medium">{log.subject}</h4>
                    </div>
                    <div className="flex gap-2">
                      {getCommTypeBadge(log.communication_type)}
                      {getDirectionBadge(log.direction)}
                    </div>
                  </div>
                  {log.content && <p className="text-sm text-muted-foreground leading-relaxed">{log.content}</p>}
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Calendar className="h-4 w-4" />
                      {new Date(log.communication_date).toLocaleString()}
                    </div>
                    <div className="flex items-center gap-1">
                      <User className="h-4 w-4" />
                      User {log.created_by}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  )
}
