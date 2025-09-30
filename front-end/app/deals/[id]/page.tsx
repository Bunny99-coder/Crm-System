"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowLeft, DollarSign, Calendar, User, Building, Edit, ExternalLink } from "lucide-react"
import { api, type Deal, type Lead, type Property } from "@/lib/api"
import { DealNotes } from "@/components/deal-notes"
import { DealTasks } from "@/components/deal-tasks"
import { DealEvents } from "@/components/deal-events"
import { DealActivity } from "@/components/deal-activity"
import { PageBreadcrumb } from "@/components/page-breadcrumb"

export default function DealDetailPage() {
  const params = useParams()
  const router = useRouter()
  const dealId = Number(params.id)

  const [deal, setDeal] = useState<Deal | null>(null)
  const [lead, setLead] = useState<Lead | null>(null)
  const [property, setProperty] = useState<Property | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    loadDealDetails()
  }, [dealId])

 const loadDealDetails = async () => {
  try {
    setIsLoading(true)

    const allDeals = await api.getDeals()
    const dealData = allDeals.find(d => d.id === dealId)
    if (!dealData) throw new Error("Deal not found")
    setDeal(dealData)

    const allLeads = await api.getLeads()
    const leadData = allLeads.find(l => l.id === dealData.lead_id)
    setLead(leadData || null)

    const allProperties = await api.getProperties()
    const propertyData = allProperties.find(p => p.id === dealData.property_id)
    setProperty(propertyData || null)
  } catch (err) {
    setError("Failed to load deal details")
  } finally {
    setIsLoading(false)
  }
}


  const getStatusBadge = (status: string) => {
    switch (status) {
      case "Pending":
        return (
          <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">
            Pending
          </Badge>
        )
      case "Closed-Won":
        return (
          <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
            Closed-Won
          </Badge>
        )
      case "Closed-Lost":
        return (
          <Badge variant="outline" className="bg-red-50 text-red-700 border-red-200">
            Closed-Lost
          </Badge>
        )
      default:
        return <Badge variant="outline">{status}</Badge>
    }
  }

  const getStageName = (stageId: number) => {
    const stageMap: { [key: number]: string } = {
      1: "Initial Contact",
      2: "Qualification",
      3: "Proposal",
      4: "Negotiation",
      5: "Closing",
    }
    return stageMap[stageId] || `Stage ${stageId}`
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Loading deal details...</div>
      </div>
    )
  }

  if (error || !deal) {
    return (
      <div className="flex flex-col items-center justify-center h-64 space-y-4">
        <div className="text-destructive">{error || "Deal not found"}</div>
        <Button onClick={() => router.push("/deals")} variant="outline">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Deals
        </Button>
      </div>
    )
  }

  const breadcrumbItems = [{ label: "Deals", href: "/deals" }, { label: `Deal #${deal.id}` }]

  return (
    <div className="space-y-6">
      <PageBreadcrumb items={breadcrumbItems} />

      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Button onClick={() => router.push("/deals")} variant="outline" size="sm">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-balance">Deal #{deal.id}</h1>
            <p className="text-muted-foreground">Manage deal details and track progress</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm">
            <Edit className="mr-2 h-4 w-4" />
            Edit Deal
          </Button>
          {getStatusBadge(deal.deal_status)}
        </div>
      </div>

      {/* Deal Overview */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Deal Value</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-cyan-600">${deal.deal_amount.toLocaleString()}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Stage</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">{getStageName(deal.stage_id)}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Lead</CardTitle>
            <User className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <div className="text-2xl font-bold text-purple-600">Lead #{deal.lead_id}</div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => router.push(`/leads/${deal.lead_id}`)}
                className="p-1 h-auto"
              >
                <ExternalLink className="h-4 w-4" />
              </Button>
            </div>
            {lead && <p className="text-sm text-muted-foreground mt-1">{lead.notes.substring(0, 50)}...</p>}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Property</CardTitle>
            <Building className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <div className="text-2xl font-bold text-green-600">
                {property?.name || `Property ${deal.property_id}`}
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => router.push(`/properties/${deal.property_id}`)}
                className="p-1 h-auto"
              >
                <ExternalLink className="h-4 w-4" />
              </Button>
            </div>
            {property && <p className="text-sm text-muted-foreground mt-1">{property.unit_no}</p>}
          </CardContent>
        </Card>
      </div>

      {/* Tabbed Content */}
      <Tabs defaultValue="notes" className="space-y-4">
        <TabsList>
          <TabsTrigger value="notes">Notes</TabsTrigger>
          <TabsTrigger value="tasks">Tasks</TabsTrigger>
          <TabsTrigger value="events">Events</TabsTrigger>
          <TabsTrigger value="activity">Activity</TabsTrigger>
        </TabsList>

        <TabsContent value="notes">
          <DealNotes dealId={dealId} />
        </TabsContent>

        <TabsContent value="tasks">
          <DealTasks dealId={dealId} />
        </TabsContent>

        <TabsContent value="events">
          <DealEvents dealId={dealId} />
        </TabsContent>

        <TabsContent value="activity">
          <DealActivity dealId={dealId} />
        </TabsContent>
      </Tabs>
    </div>
  )
}
