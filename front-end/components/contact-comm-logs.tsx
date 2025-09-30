"use client"

import { useState, useEffect, useCallback } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Plus, Phone, Mail, MessageSquare, Calendar, User } from "lucide-react"
import { api, type CommLog } from "@/lib/api"

interface ContactCommLogsProps {
  contactId: number
}

export function ContactCommLogs({ contactId }: ContactCommLogsProps) {
  const [commLogs, setCommLogs] = useState<CommLog[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isCreating, setIsCreating] = useState(false)
  const [editingLogId, setEditingLogId] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [formData, setFormData] = useState({
    interaction_type: "Email" as CommLog["interaction_type"],
    interaction_date: "",
    notes: "",
  })

  const loadCommLogs = useCallback(async () => {
    let isMounted = true // Flag to prevent state updates on unmounted components
    try {
      setIsLoading(true)
      setError(null)
      const commLogsData = await api.getCommLogsForContact(contactId)
      if (isMounted) {
        setCommLogs(commLogsData)
      }
    } catch (err: any) {
      if (isMounted) {
        setError(err.message || "Failed to load communication logs")
        console.error("Failed to load communication logs:", err)
      }
    }
    finally {
      if (isMounted) {
        setIsLoading(false)
      }
    }
    return () => {
      isMounted = false
    }
  }, [contactId]) // Dependency on contactId

  useEffect(() => {
    const cleanup = loadCommLogs() // Call loadCommLogs and get its cleanup function
    return () => {
      if (typeof cleanup === 'function') {
        cleanup(); // Execute the cleanup function returned by loadCommLogs
      }
    };
  }, [loadCommLogs]) // Dependency on loadCommLogs

  const handleCreateCommLog = async () => {
    if (!formData.interaction_type || !formData.interaction_date) {
      setError("Interaction type and date are required.")
      return
    }

    try {
      setError(null)
      const createdLog = await api.createContactCommLog(contactId, {
        interaction_type: formData.interaction_type,
        interaction_date: new Date(formData.interaction_date).toISOString(),
        notes: formData.notes || null, // Ensure notes is null if empty string
      })
      setCommLogs((prev) => [...prev, createdLog])
      setFormData({
        interaction_type: "Email",
        interaction_date: "",
        notes: "",
      })
      setIsCreating(false)
    } catch (err: any) {
      setError(err.message || "Failed to create communication log")
      console.error("Failed to create communication log:", err)
    }
  }

  const startEditCommLog = (log: CommLog) => {
    setEditingLogId(log.id!)
    // Format date for datetime-local input
    const formattedDate = new Date(log.interaction_date).toISOString().slice(0, 16)
    setFormData({
      interaction_type: log.interaction_type,
      interaction_date: formattedDate,
      notes: log.notes || "",
    })
  }

  const cancelEdit = () => {
    setEditingLogId(null)
    setFormData({
      interaction_type: "Email",
      interaction_date: "",
      notes: "",
    })
    setError(null)
  }

  const handleUpdateCommLog = async (logId: number) => {
    if (!formData.interaction_type || !formData.interaction_date) {
      setError("Interaction type and date are required.")
      return
    }

    try {
      setError(null)
      const updatedLog = await api.updateContactCommLog(contactId, logId, {
        interaction_type: formData.interaction_type,
        interaction_date: new Date(formData.interaction_date).toISOString(),
        notes: formData.notes || null,
      })
      setCommLogs((prev) => prev.map((log) => (log.id === logId ? updatedLog : log)))
      cancelEdit()
    } catch (err: any) {
      setError(err.message || "Failed to update communication log")
      console.error("Failed to update communication log:", err)
    }
  }

  const handleDeleteCommLog = async (logId: number) => {
    try {
      setError(null)
      await api.deleteContactCommLog(contactId, logId)
      setCommLogs((prev) => prev.filter((log) => log.id !== logId))
    } catch (err: any) {
      setError(err.message || "Failed to delete communication log")
      console.error("Failed to delete communication log:", err)
    }
  }

  const getCommTypeIcon = (type: CommLog["interaction_type"]) => {
    switch (type) {
      case "Email":
        return <Mail className="h-4 w-4" />
      case "Call":
        return <Phone className="h-4 w-4" />
      case "SMS":
        return <MessageSquare className="h-4 w-4" />
      case "Meeting":
        return <Calendar className="h-4 w-4" />
      default:
        return <MessageSquare className="h-4 w-4" />
    }
  }

  const getCommTypeBadge = (type: CommLog["interaction_type"]) => {
    const colors = {
      Email: "bg-blue-50 text-blue-700 border-blue-200",
      Call: "bg-green-50 text-green-700 border-green-200",
      Meeting: "bg-purple-50 text-purple-700 border-purple-200",
      SMS: "bg-yellow-50 text-yellow-700 border-yellow-200",
      Other: "bg-gray-50 text-gray-700 border-gray-200",
    }
    return (
      <Badge variant="outline" className={colors[type] || colors.Other}>
        {type}
      </Badge>
    )
  }

  if (isLoading) {
    return <div className="text-center py-8">Loading communication logs...</div>
  }

  return (
    <div className="space-y-4">
      {/* Error display */}
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* Create Communication Log */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Communication Logs</CardTitle>
              <CardDescription>Track all communications with this contact</CardDescription>
            </div>
            {!isCreating && editingLogId === null && (
              <Button onClick={() => setIsCreating(true)} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                <Plus className="mr-2 h-4 w-4" />
                Log Communication
              </Button>
            )}
          </div>
        </CardHeader>
        {(isCreating || editingLogId !== null) && (
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <Select
                value={formData.interaction_type}
                onValueChange={(value: CommLog["interaction_type"]) =>
                  setFormData((prev) => ({ ...prev, interaction_type: value }))
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Email">Email</SelectItem>
                  <SelectItem value="Call">Call</SelectItem>
                  <SelectItem value="Meeting">Meeting</SelectItem>
                  <SelectItem value="SMS">SMS</SelectItem>
                  <SelectItem value="Other">Other</SelectItem>
                </SelectContent>
              </Select>
              <Input
                type="datetime-local"
                value={formData.interaction_date}
                onChange={(e) => setFormData((prev) => ({ ...prev, interaction_date: e.target.value }))}
              />
            </div>
            <Textarea
              placeholder="Communication notes..."
              value={formData.notes}
              onChange={(e) => setFormData((prev) => ({ ...prev, notes: e.target.value }))}
              rows={3}
            />
            <div className="flex gap-2">
              {editingLogId !== null ? (
                <Button
                  onClick={() => handleUpdateCommLog(editingLogId)}
                  size="sm"
                  className="bg-blue-600 hover:bg-blue-700"
                >
                  Save Changes
                </Button>
              ) : (
                <Button onClick={handleCreateCommLog} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                  Log Communication
                </Button>
              )}
              <Button onClick={cancelEdit} variant="outline" size="sm">
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
                      {getCommTypeIcon(log.interaction_type)}
                      <h4 className="font-medium">{log.interaction_type}</h4>
                    </div>
                    <div className="flex gap-2">
                      {getCommTypeBadge(log.interaction_type)}
                      <Button size="sm" variant="outline" onClick={() => startEditCommLog(log)}>
                        Edit
                      </Button>
                      <Button size="sm" variant="destructive" onClick={() => handleDeleteCommLog(log.id!)}>
                        Delete
                      </Button>
                    </div>
                  </div>
                  {log.notes && <p className="text-sm text-muted-foreground leading-relaxed">{log.notes}</p>}
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Calendar className="h-4 w-4" />
                      {new Date(log.interaction_date).toLocaleDateString()} {new Date(log.interaction_date).toLocaleTimeString()}
                    </div>
                    <div className="flex items-center gap-1">
                      <User className="h-4 w-4" />
                      User {log.user_id}
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