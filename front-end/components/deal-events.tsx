"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Plus, Calendar, MapPin, Users } from "lucide-react"
import { api, Event } from "@/lib/api"



interface DealEventsProps {
  dealId: number
}

export function DealEvents({ dealId }: DealEventsProps) {
  const [events, setEvents] = useState<Event[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isCreating, setIsCreating] = useState(false)
  const [formData, setFormData] = useState({
    event_name: "",
    event_description: "",
    event_date: "",
  })

  useEffect(() => {
    loadEvents()
  }, [dealId])

  const loadEvents = async () => {
    try {
      setIsLoading(true)
      const eventsData = await api.getEventsForDeal(dealId)
      setEvents(eventsData)
    } catch (err) {
      console.error("Failed to load events:", err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateEvent = async () => {
    if (!formData.event_name.trim() || !formData.event_date) return

    try {
      await api.createEventForDeal(dealId, formData)
      setFormData({
        event_name: "",
        event_description: "",
        event_date: "",
      })
      setIsCreating(false)
      loadEvents()
    } catch (err) {
      console.error("Failed to create event:", err)
    }
  }

  const getEventTypeBadge = (type: string) => {
    const colors = {
      Meeting: "bg-blue-50 text-blue-700 border-blue-200",
      Call: "bg-green-50 text-green-700 border-green-200",
      "Site Visit": "bg-purple-50 text-purple-700 border-purple-200",
      Presentation: "bg-orange-50 text-orange-700 border-orange-200",
      Other: "bg-gray-50 text-gray-700 border-gray-200",
    }
    return (
      <Badge variant="outline" className={colors[type as keyof typeof colors] || colors.Other}>
        {type}
      </Badge>
    )
  }

  const isUpcoming = (eventDate: string) => {
    const eventDateTime = new Date(eventDate)
    return eventDateTime > new Date()
  }

  if (isLoading) {
    return <div className="text-center py-8">Loading events...</div>
  }

  return (
    <div className="space-y-4">
      {/* Create Event */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Deal Events</CardTitle>
              <CardDescription>Schedule and track meetings, calls, and other events</CardDescription>
            </div>
            {!isCreating && (
              <Button onClick={() => setIsCreating(true)} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                <Plus className="mr-2 h-4 w-4" />
                Add Event
              </Button>
            )}
          </div>
        </CardHeader>
        {isCreating && (
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <Input
                placeholder="Event name"
                value={formData.event_name}
                onChange={(e) => setFormData((prev) => ({ ...prev, event_name: e.target.value }))}
              />
              <Input
                type="date"
                value={formData.event_date}
                onChange={(e) => setFormData((prev) => ({ ...prev, event_date: e.target.value }))}
              />
            </div>
            <Textarea
              placeholder="Event description"
              value={formData.event_description}
              onChange={(e) => setFormData((prev) => ({ ...prev, event_description: e.target.value }))}
              rows={2}
            />
            <div className="flex gap-2">
              <Button onClick={handleCreateEvent} size="sm" className="bg-cyan-600 hover:bg-cyan-700">
                Create Event
              </Button>
              <Button
                onClick={() => {
                  setIsCreating(false)
                  setFormData({
                    event_name: "",
                    event_description: "",
                    event_date: "",
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

      {/* Events List */}
      <div className="space-y-4">
        {events.length === 0 ? (
          <Card>
            <CardContent className="text-center py-8">
              <p className="text-muted-foreground">No events scheduled. Add your first event to get started.</p>
            </CardContent>
          </Card>
        ) : (
          events.map((event) => (
            <Card key={event.id} className={isUpcoming(event.event_date) ? "border-cyan-200" : ""}>
              <CardContent className="pt-6">
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <h4 className="font-medium">{event.event_name}</h4>
                    <div className="flex gap-2">
                      {isUpcoming(event.event_date) && (
                        <Badge variant="outline" className="bg-cyan-50 text-cyan-700 border-cyan-200">
                          Upcoming
                        </Badge>
                      )}
                    </div>
                  </div>
                  {event.event_description && <p className="text-sm text-muted-foreground">{event.event_description}</p>}
                  <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Calendar className="h-4 w-4" />
                      {new Date(event.event_date).toLocaleDateString()}
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
