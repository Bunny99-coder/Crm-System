"use client"

import { useState, useEffect } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { MoreHorizontal, Search, Plus, Edit, Trash2, DollarSign, TrendingUp, CheckCircle, Eye } from "lucide-react"
import { api, type Deal, type Lead, type Property } from "@/lib/api"
import { useAuth, ROLE_SALES_AGENT, ROLE_RECEPTION } from "@/lib/auth" // Import useAuth and role constants

export function DealsManagement() {
  const [deals, setDeals] = useState<Deal[]>([])
  const [leads, setLeads] = useState<Lead[]>([])
  const [properties, setProperties] = useState<Property[]>([])
  const [searchTerm, setSearchTerm] = useState("")
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [selectedDeal, setSelectedDeal] = useState<Deal | null>(null)
  const [formData, setFormData] = useState({
    lead_id: "",
    property_id: "",
    stage_id: "",
    deal_status: "Pending" as "Pending" | "Closed-Won" | "Closed-Lost",
    deal_amount: "",
  })

  const { user, hasRole } = useAuth() // Use the useAuth hook

  // Load data on component mount
  useEffect(() => {
    loadAllData()
  }, [])

  const loadAllData = async () => {
    try {
      setIsLoading(true)
      const dealsData = await api.getDeals()
      setDeals(dealsData || [])
    } catch (err) {
      setError("Failed to load deals")
    }

    try {
      const leadsData = await api.getLeads()
      setLeads(leadsData || [])
    } catch (err) {
      setError("Failed to load leads")
    }

    try {
      const propertiesData = await api.getProperties()
      setProperties(propertiesData || [])
    } catch (err) {
      setError("Failed to load properties")
    }

    setIsLoading(false)
  }

  const handleCreateDeal = async () => {
    try {
      const dealData = {
        lead_id: Number.parseInt(formData.lead_id),
        property_id: Number.parseInt(formData.property_id),
        stage_id: Number.parseInt(formData.stage_id),
        deal_status: formData.deal_status,
        deal_amount: Number.parseFloat(formData.deal_amount),
      }
      await api.createDeal(dealData)
      setIsCreateDialogOpen(false)
      resetForm()
      loadAllData()
    } catch (err) {
      setError("Failed to create deal")
    }
  }

  const handleEditDeal = async () => {
    if (!selectedDeal?.id) return

    try {
      const dealData = {
        lead_id: Number.parseInt(formData.lead_id),
        property_id: Number.parseInt(formData.property_id),
        stage_id: Number.parseInt(formData.stage_id),
        deal_status: formData.deal_status,
        deal_amount: Number.parseFloat(formData.deal_amount),
      }
      await api.updateDeal(selectedDeal.id, dealData)
      setIsEditDialogOpen(false)
      resetForm()
      setSelectedDeal(null)
      loadAllData()
    } catch (err) {
      setError("Failed to update deal")
    }
  }

  const handleDeleteDeal = async (id: number) => {
    if (!confirm("Are you sure you want to delete this deal?")) return

    try {
      await api.deleteDeal(id)
      loadAllData()
    } catch (err) {
      setError("Failed to delete deal")
    }
  }

  const resetForm = () => {
    setFormData({
      lead_id: "",
      property_id: "",
      stage_id: "",
      deal_status: "Pending",
      deal_amount: "",
    })
  }

  const openEditDialog = (deal: Deal) => {
    setSelectedDeal(deal)
    setFormData({
      lead_id: deal.lead_id.toString(),
      property_id: deal.property_id.toString(),
      stage_id: deal.stage_id.toString(),
      deal_status: deal.deal_status,
      deal_amount: deal.deal_amount.toString(),
    })
    setIsEditDialogOpen(true)
  }

  const getLeadInfo = (leadId: number) => {
    const lead = leads.find((l) => l.id === leadId)
    return lead ? `Lead ${leadId}` : `Lead ${leadId}`
  }

  const getPropertyName = (propertyId: number) => {
    const property = properties.find((p) => p.id === propertyId)
    return property ? property.name : `Property ${propertyId}`
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

  const getStageBadge = (stageId: number) => {
    // Mock stage mapping - in real app this would come from API
    const stageMap: { [key: number]: { label: string; color: string } } = {
      1: { label: "Initial Contact", color: "bg-blue-50 text-blue-700 border-blue-200" },
      2: { label: "Qualification", color: "bg-purple-50 text-purple-700 border-purple-200" },
      3: { label: "Proposal", color: "bg-orange-50 text-orange-700 border-orange-200" },
      4: { label: "Negotiation", color: "bg-yellow-50 text-yellow-700 border-yellow-200" },
      5: { label: "Closing", color: "bg-green-50 text-green-700 border-green-200" },
    }

    const stage = stageMap[stageId] || { label: `Stage ${stageId}`, color: "bg-gray-50 text-gray-700 border-gray-200" }
    return (
      <Badge variant="outline" className={stage.color}>
        {stage.label}
      </Badge>
    )
  }

  const filteredDeals = deals.filter(
    (deal) =>
      getLeadInfo(deal.lead_id).toLowerCase().includes(searchTerm.toLowerCase()) ||
      getPropertyName(deal.property_id).toLowerCase().includes(searchTerm.toLowerCase()) ||
      deal.deal_status.toLowerCase().includes(searchTerm.toLowerCase()) ||
      deal.deal_amount.toString().includes(searchTerm),
  )

  const pendingDeals = deals.filter((d) => d.deal_status === "Pending").length
  const wonDeals = deals.filter((d) => d.deal_status === "Closed-Won").length
  const lostDeals = deals.filter((d) => d.deal_status === "Closed-Lost").length
  const totalValue = deals
    .filter((d) => d.deal_status === "Closed-Won")
    .reduce((sum, deal) => sum + deal.deal_amount, 0)

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Loading deals...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <h1 className="text-3xl font-bold text-balance text-foreground">Deals Management</h1>
                <p className="text-muted-foreground text-pretty">Track and manage your sales deals and revenue pipeline</p>
              </div>
              {(hasRole(ROLE_RECEPTION) || hasRole(ROLE_SALES_AGENT)) && (
                <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                  <DialogTrigger asChild>
                    <Button className="bg-cyan-600 hover:bg-cyan-700">
                      <Plus className="mr-2 h-4 w-4" />
                      Add Deal
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="sm:max-w-[500px]">
                    <DialogHeader>
                      <DialogTitle>Create New Deal</DialogTitle>
                      <DialogDescription>Add a new deal to your sales pipeline.</DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4 py-4">
                      <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="lead_id" className="text-right">
                          Lead
                        </Label>
                        <Select
                          value={formData.lead_id}
                          onValueChange={(value) => setFormData((prev) => ({ ...prev, lead_id: value }))}
                        >
                          <SelectTrigger className="col-span-3">
                            <SelectValue placeholder="Select lead" />
                          </SelectTrigger>
                          <SelectContent>
                            {leads.map((lead) => (
                              <SelectItem key={lead.id} value={lead.id!.toString()}>
      Lead {lead.id} - {lead.notes ? lead.notes.substring(0, 30) : "No notes"}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="property_id" className="text-right">
                          Property
                        </Label>
                        <Select
                          value={formData.property_id}
                          onValueChange={(value) => setFormData((prev) => ({ ...prev, property_id: value }))}
                        >
                          <SelectTrigger className="col-span-3">
                            <SelectValue placeholder="Select property" />
                          </SelectTrigger>
                          <SelectContent>
                            {properties.map((property) => (
                              <SelectItem key={property.id} value={property.id!.toString()}>
                                {property.name} - {property.unit_no}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="stage_id" className="text-right">
                          Stage
                        </Label>
                        <Select
                          value={formData.stage_id}
                          onValueChange={(value) => setFormData((prev) => ({ ...prev, stage_id: value }))}
                        >
                          <SelectTrigger className="col-span-3">
                            <SelectValue placeholder="Select stage" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="1">Initial Contact</SelectItem>
                            <SelectItem value="2">Qualification</SelectItem>
                            <SelectItem value="3">Proposal</SelectItem>
                            <SelectItem value="4">Negotiation</SelectItem>
                            <SelectItem value="5">Closing</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="deal_status" className="text-right">
                          Status
                        </Label>
                        <Select
                          value={formData.deal_status}
                          onValueChange={(value: "Pending" | "Closed-Won" | "Closed-Lost") =>
                            setFormData((prev) => ({ ...prev, deal_status: value }))
                          }
                        >
                          <SelectTrigger className="col-span-3">
                            <SelectValue placeholder="Select status" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="Pending">Pending</SelectItem>
                            <SelectItem value="Closed-Won">Closed-Won</SelectItem>
                            <SelectItem value="Closed-Lost">Closed-Lost</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="grid grid-cols-4 items-center gap-4">
                        <Label htmlFor="deal_amount" className="text-right">
                          Amount
                        </Label>
                        <Input
                          id="deal_amount"
                          type="number"
                          step="0.01"
                          value={formData.deal_amount}
                          onChange={(e) => setFormData((prev) => ({ ...prev, deal_amount: e.target.value }))}
                          className="col-span-3"
                          placeholder="Deal amount"
                          required
                        />
                      </div>
                    </div>
                    <DialogFooter>
                      <Button type="submit" onClick={handleCreateDeal} className="bg-cyan-600 hover:bg-cyan-700">
                        Create Deal
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              )}
            </div>
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4">
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5 text-cyan-600" />
            <h3 className="font-semibold">Total Deals</h3>
          </div>
          <p className="text-2xl font-bold text-cyan-600 mt-2">{deals.length}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <DollarSign className="h-5 w-5 text-yellow-600" />
            <h3 className="font-semibold">Pending</h3>
          </div>
          <p className="text-2xl font-bold text-yellow-600 mt-2">{pendingDeals}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <CheckCircle className="h-5 w-5 text-green-600" />
            <h3 className="font-semibold">Won</h3>
          </div>
          <p className="text-2xl font-bold text-green-600 mt-2">{wonDeals}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <DollarSign className="h-5 w-5 text-purple-600" />
            <h3 className="font-semibold">Total Value</h3>
          </div>
          <p className="text-2xl font-bold text-purple-600 mt-2">${totalValue.toLocaleString()}</p>
        </div>
      </div>

      {/* Search */}
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search deals..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="pl-10"
        />
      </div>

      {/* Table */}
      <div className="rounded-md border border-border bg-card">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-muted/50">
              <TableHead className="text-card-foreground">Lead</TableHead>
              <TableHead className="text-card-foreground">Property</TableHead>
              <TableHead className="text-card-foreground">Stage</TableHead>
              <TableHead className="text-card-foreground">Status</TableHead>
              <TableHead className="text-card-foreground">Amount</TableHead>
              <TableHead className="w-[50px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredDeals.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                  {searchTerm ? "No deals found matching your search." : "No deals yet. Create your first deal!"}
                </TableCell>
              </TableRow>
            ) : (
              filteredDeals.map((deal) => {
                const canEditDelete = hasRole(ROLE_RECEPTION) || (user && deal.created_by === user.id);
                return (
                  <TableRow key={deal.id} className="hover:bg-muted/50">
                    <TableCell className="font-medium text-card-foreground">{getLeadInfo(deal.lead_id)}</TableCell>
                    <TableCell className="text-card-foreground">{getPropertyName(deal.property_id)}</TableCell>
                    <TableCell>{getStageBadge(deal.stage_id)}</TableCell>
                    <TableCell>{getStatusBadge(deal.deal_status)}</TableCell>
                    <TableCell className="text-card-foreground font-medium">
                      ${deal.deal_amount.toLocaleString()}
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => (window.location.href = `/deals/${deal.id}`)}>
                            <Eye className="mr-2 h-4 w-4" />
                            View Details
                          </DropdownMenuItem>
                          {canEditDelete && (
                            <DropdownMenuItem onClick={() => openEditDialog(deal)}>
                              <Edit className="mr-2 h-4 w-4" />
                              Edit
                            </DropdownMenuItem>
                          )}
                          {canEditDelete && (
                            <DropdownMenuItem
                              className="text-destructive"
                              onClick={() => deal.id && handleDeleteDeal(deal.id)}
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Delete
                            </DropdownMenuItem>
                          )}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </div>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Edit Deal</DialogTitle>
            <DialogDescription>Update the deal information.</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_lead_id" className="text-right">
                Lead
              </Label>
              <Select
                value={formData.lead_id}
                onValueChange={(value) => setFormData((prev) => ({ ...prev, lead_id: value }))}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select lead" />
                </SelectTrigger>
                <SelectContent>
                  {leads.map((lead) => (
                    <SelectItem key={lead.id} value={lead.id!.toString()}>
Lead {lead.id} - {lead.notes ? lead.notes.substring(0, 30) : "No notes"}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_property_id" className="text-right">
                Property
              </Label>
              <Select
                value={formData.property_id}
                onValueChange={(value) => setFormData((prev) => ({ ...prev, property_id: value }))}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select property" />
                </SelectTrigger>
                <SelectContent>
                  {properties.map((property) => (
                    <SelectItem key={property.id} value={property.id!.toString()}>
                      {property.name} - {property.unit_no}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_stage_id" className="text-right">
                Stage
              </Label>
              <Select
                value={formData.stage_id}
                onValueChange={(value) => setFormData((prev) => ({ ...prev, stage_id: value }))}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select stage" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1">Initial Contact</SelectItem>
                  <SelectItem value="2">Qualification</SelectItem>
                  <SelectItem value="3">Proposal</SelectItem>
                  <SelectItem value="4">Negotiation</SelectItem>
                  <SelectItem value="5">Closing</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_deal_status" className="text-right">
                Status
              </Label>
              <Select
                value={formData.deal_status}
                onValueChange={(value: "Pending" | "Closed-Won" | "Closed-Lost") =>
                  setFormData((prev) => ({ ...prev, deal_status: value }))
                }
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Pending">Pending</SelectItem>
                  <SelectItem value="Closed-Won">Closed-Won</SelectItem>
                  <SelectItem value="Closed-Lost">Closed-Lost</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_deal_amount" className="text-right">
                Amount
              </Label>
              <Input
                id="edit_deal_amount"
                type="number"
                step="0.01"
                value={formData.deal_amount}
                onChange={(e) => setFormData((prev) => ({ ...prev, deal_amount: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="submit" onClick={handleEditDeal} className="bg-cyan-600 hover:bg-cyan-700">
              Update Deal
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
