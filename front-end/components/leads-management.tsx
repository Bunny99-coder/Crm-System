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
import { Textarea } from "@/components/ui/textarea"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { MoreHorizontal, Search, Plus, Edit, Trash2, Users, Target, TrendingUp, UserCheck } from "lucide-react"
import { api, type Lead, type Contact, type Property, type User } from "@/lib/api"
import { useAuth, ROLE_RECEPTION, ROLE_SALES_AGENT } from "@/lib/auth" // Import useAuth and role constants
import { useToast } from "@/hooks/use-toast" // Import useToast

export function LeadsManagement() {
  const { toast } = useToast() // Initialize toast

  const [leads, setLeads] = useState<Lead[]>([])
  const [contacts, setContacts] = useState<Contact[]>([])
  const [properties, setProperties] = useState<Property[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [searchTerm, setSearchTerm] = useState("")
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [selectedLead, setSelectedLead] = useState<Lead | null>(null)
  const [formData, setFormData] = useState({
    contact_id: "",
    property_id: "",
    source_id: "",
    status_id: "",
    assigned_to: "",
    notes: "",
  })

  const { user, hasRole } = useAuth() // Use the useAuth hook

  // Load data on component mount or when user changes
  useEffect(() => {
    if (user) {
      loadAllData(user)
    }
  }, [user])

  const loadAllData = async (currentUser: User) => {
    console.log("loadAllData called for user:", currentUser.username, "role:", currentUser.role_id)
    try {
      setIsLoading(true)
      let leadsData: Lead[]
      if (currentUser.role_id === ROLE_SALES_AGENT) {
        console.log("Fetching leads for sales agent, assigned to:", currentUser.id)
        leadsData = await api.getLeads(currentUser.id)
      } else {
        console.log("Fetching all leads for non-sales agent role.")
        leadsData = await api.getLeads()
      }
      setLeads(leadsData)
    } catch (err) {
      console.error("Error loading leads:", err)
      setError("Failed to load leads")
    }

    try {
      const contactsData = await api.getContacts()
      setContacts(contactsData)
    } catch (err) {
      setError("Failed to load contacts")
    }

    try {
      const propertiesData = await api.getProperties()
      setProperties(propertiesData)
    } catch (err) {
      setError("Failed to load properties")
    }

    try {
      if (currentUser.role_id === ROLE_SALES_AGENT) {
        // If sales agent, only show themselves in the assigned_to dropdown
        setUsers(currentUser ? [currentUser] : [])
      } else {
        // For other roles (e.g., Reception), fetch all users
        const usersData = await api.getUsers()
        setUsers(usersData)
      }
    } catch (err) {
      setError("Failed to load users")
    }

    setIsLoading(false)
  }

  const handleCreateLead = async () => {
    try {
      const leadData = {
        contact_id: Number.parseInt(formData.contact_id),
        property_id: Number.parseInt(formData.property_id),
        source_id: Number.parseInt(formData.source_id),
        status_id: Number.parseInt(formData.status_id),
        assigned_to: Number.parseInt(formData.assigned_to),
        notes: formData.notes,
      }
      console.log("Attempting to create lead with data:", JSON.stringify(leadData, null, 2))
      await api.createLead(leadData)
      setIsCreateDialogOpen(false)
      resetForm()
      loadAllData(user!)
        toast({
          title: 'Lead created',
          description: 'The new lead has been successfully added.',
        });
      } catch (error: any) {
        console.error('Failed to create lead:', error);
        toast({
          title: 'Error',
          description: `Failed to create lead: ${error.message || 'Unknown error'}`,
          variant: 'destructive',
        });
      } finally {
        setIsLoading(false);
      }
  }

  const handleEditLead = async () => {
    if (!selectedLead?.id) return

    try {
      const leadData = {
        contact_id: Number.parseInt(formData.contact_id),
        property_id: Number.parseInt(formData.property_id),
        source_id: Number.parseInt(formData.source_id),
        status_id: Number.parseInt(formData.status_id),
        assigned_to: Number.parseInt(formData.assigned_to),
        notes: formData.notes,
      }
      console.log("Attempting to update lead with ID:", selectedLead.id, "and data:", leadData)
      await api.updateLead(selectedLead.id, leadData)
      setIsEditDialogOpen(false)
      resetForm()
      setSelectedLead(null)
      loadAllData(user!)
    } catch (err) {
      console.error("Failed to update lead:", err)
      setError("Failed to update lead")
    }
  }

  const handleDeleteLead = async (id: number) => {
    if (!confirm("Are you sure you want to delete this lead?")) return

    try {
      console.log("Attempting to delete lead with ID:", id)
      await api.deleteLead(id)
      loadAllData(user!)
    } catch (err) {
      console.error("Failed to delete lead:", err)
      setError("Failed to delete lead")
    }
  }

  const resetForm = () => {
    setFormData({
      contact_id: "",
      property_id: "",
      source_id: "",
      status_id: "",
      assigned_to: "",
      notes: "",
    })
  }

  const openEditDialog = (lead: Lead) => {
    setSelectedLead(lead)
    setFormData({
      contact_id: lead.contact_id.toString(),
      property_id: lead.property_id.toString(),
      source_id: lead.source_id.toString(),
      status_id: lead.status_id.toString(),
      assigned_to: lead.assigned_to.toString(),
      notes: lead.notes,
    })
    setIsEditDialogOpen(true)
  }

  const getContactName = (contactId: number) => {
    const contact = contacts.find((c) => c.id === contactId)
    return contact ? `${contact.first_name} ${contact.last_name}` : `Contact ${contactId}`
  }

  const getPropertyName = (propertyId: number) => {
    const property = properties.find((p) => p.id === propertyId)
    return property ? property.name : `Property ${propertyId}`
  }

  const getUserName = (userId: number) => {
    const user = users.find((u) => u.id === userId)
    return user ? user.username : `User ${userId}`
  }

  const getStatusBadge = (statusId: number) => {
    // Mock status mapping - in real app this would come from API
    const statusMap: { [key: number]: { label: string; color: string } } = {
      1: { label: "New", color: "bg-blue-50 text-blue-700 border-blue-200" },
      2: { label: "Contacted", color: "bg-yellow-50 text-yellow-700 border-yellow-200" },
      3: { label: "Qualified", color: "bg-green-50 text-green-700 border-green-200" },
      4: { label: "Converted", color: "bg-purple-50 text-purple-700 border-purple-200" },
      5: { label: "Lost", color: "bg-red-50 text-red-700 border-red-200" },
    }

    const status = statusMap[statusId] || {
      label: `Status ${statusId}`,
      color: "bg-gray-50 text-gray-700 border-gray-200",
    }
    return (
      <Badge variant="outline" className={status.color}>
        {status.label}
      </Badge>
    )
  }

  const filteredLeads = leads.filter(
    (lead) =>
      getContactName(lead.contact_id).toLowerCase().includes(searchTerm.toLowerCase()) ||
      getPropertyName(lead.property_id).toLowerCase().includes(searchTerm.toLowerCase()) ||
      getUserName(lead.assigned_to).toLowerCase().includes(searchTerm.toLowerCase()) ||
      lead.notes.toLowerCase().includes(searchTerm.toLowerCase()),
  )

  // Mock status counts - in real app this would be calculated from actual status data
  const newLeads = leads.filter((l) => l.status_id === 1).length
  const qualifiedLeads = leads.filter((l) => l.status_id === 3).length
  const convertedLeads = leads.filter((l) => l.status_id === 4).length

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Loading leads...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold text-balance text-foreground">Leads Management</h1>
          <p className="text-muted-foreground text-pretty">Track and manage your sales leads through the pipeline</p>
        </div>
        {hasRole(ROLE_RECEPTION) && (
          <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button className="bg-cyan-600 hover:bg-cyan-700">
                <Plus className="mr-2 h-4 w-4" />
                Add Lead
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[500px]">
              <DialogHeader>
                <DialogTitle>Create New Lead</DialogTitle>
                <DialogDescription>Add a new lead to your sales pipeline.</DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="contact_id" className="text-right">
                    Contact
                  </Label>
                  <Select
                    value={formData.contact_id}
                    onValueChange={(value) => setFormData((prev) => ({ ...prev, contact_id: value }))}
                  >
                    <SelectTrigger className="col-span-3">
                      <SelectValue placeholder="Select contact" />
                    </SelectTrigger>
                    <SelectContent>
                      {contacts.map((contact) => (
                        <SelectItem key={contact.id} value={contact.id!.toString()}>
                          {contact.first_name} {contact.last_name}
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
                  <Label htmlFor="source_id" className="text-right">
                    Source ID
                  </Label>
                  <Input
                    id="source_id"
                    type="number"
                    value={formData.source_id}
                    onChange={(e) => setFormData((prev) => ({ ...prev, source_id: e.target.value }))}
                    className="col-span-3"
                    placeholder="Lead source ID"
                    required
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="status_id" className="text-right">
                    Status
                  </Label>
                  <Select
                    value={formData.status_id}
                    onValueChange={(value) => setFormData((prev) => ({ ...prev, status_id: value }))}
                  >
                    <SelectTrigger className="col-span-3">
                      <SelectValue placeholder="Select status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="1">New</SelectItem>
                      <SelectItem value="2">Contacted</SelectItem>
                      <SelectItem value="3">Qualified</SelectItem>
                      <SelectItem value="4">Converted</SelectItem>
                      <SelectItem value="5">Lost</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="assigned_to" className="text-right">
                    Assigned To
                  </Label>
                  <Select
                    value={formData.assigned_to}
                    onValueChange={(value) => setFormData((prev) => ({ ...prev, assigned_to: value }))}
                  >
                    <SelectTrigger className="col-span-3">
                      <SelectValue placeholder="Select user" />
                    </SelectTrigger>
                    <SelectContent>
                      {users.map((user) => (
                        <SelectItem key={user.id} value={user.id!.toString()}>
                          {user.username}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="notes" className="text-right">
                    Notes
                  </Label>
                  <Textarea
                    id="notes"
                    value={formData.notes}
                    onChange={(e) => setFormData((prev) => ({ ...prev, notes: e.target.value }))}
                    className="col-span-3"
                    placeholder="Lead notes and comments"
                    rows={3}
                  />
                </div>
              </div>
              <DialogFooter>
                <Button type="submit" onClick={handleCreateLead} className="bg-cyan-600 hover:bg-cyan-700">
                  Create Lead
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
            <Users className="h-5 w-5 text-cyan-600" />
            <h3 className="font-semibold">Total Leads</h3>
          </div>
          <p className="text-2xl font-bold text-cyan-600 mt-2">{leads.length}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <Target className="h-5 w-5 text-blue-600" />
            <h3 className="font-semibold">New Leads</h3>
          </div>
          <p className="text-2xl font-bold text-blue-600 mt-2">{newLeads}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <UserCheck className="h-5 w-5 text-green-600" />
            <h3 className="font-semibold">Qualified</h3>
          </div>
          <p className="text-2xl font-bold text-green-600 mt-2">{qualifiedLeads}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5 text-purple-600" />
            <h3 className="font-semibold">Converted</h3>
          </div>
          <p className="text-2xl font-bold text-purple-600 mt-2">{convertedLeads}</p>
        </div>
      </div>

      {/* Search */}
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search leads..."
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
              <TableHead className="text-card-foreground">Contact</TableHead>
              <TableHead className="text-card-foreground">Property</TableHead>
              <TableHead className="text-card-foreground">Status</TableHead>
              <TableHead className="text-card-foreground">Assigned To</TableHead>
              <TableHead className="text-card-foreground">Notes</TableHead>
              <TableHead className="w-[50px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredLeads.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                  {searchTerm ? "No leads found matching your search." : "No leads yet. Create your first lead!"}
                </TableCell>
              </TableRow>
            ) : (
              filteredLeads.map((lead) => (
                <TableRow key={lead.id} className="hover:bg-muted/50">
                  <TableCell className="font-medium text-card-foreground">{getContactName(lead.contact_id)}</TableCell>
                  <TableCell className="text-card-foreground">{getPropertyName(lead.property_id)}</TableCell>
                  <TableCell>{getStatusBadge(lead.status_id)}</TableCell>
                  <TableCell className="text-card-foreground">{getUserName(lead.assigned_to)}</TableCell>
                  <TableCell className="text-card-foreground max-w-xs truncate">{lead.notes || "No notes"}</TableCell>
                  <TableCell>
                    {hasRole(ROLE_RECEPTION) && (
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => openEditDialog(lead)}>
                            <Edit className="mr-2 h-4 w-4" />
                            Edit
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            className="text-destructive"
                            onClick={() => lead.id && handleDeleteLead(lead.id)}
                          >
                            <Trash2 className="mr-2 h-4 w-4" />
                            Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    )}
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Edit Lead</DialogTitle>
            <DialogDescription>Update the lead information.</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_contact_id" className="text-right">
                Contact
              </Label>
              <Select
                value={formData.contact_id}
                onValueChange={(value) => setFormData((prev) => ({ ...prev, contact_id: value }))}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select contact" />
                </SelectTrigger>
                <SelectContent>
                  {contacts.map((contact) => (
                    <SelectItem key={contact.id} value={contact.id!.toString()}>
                      {contact.first_name} {contact.last_name}
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
              <Label htmlFor="edit_source_id" className="text-right">
                Source ID
              </Label>
              <Input
                id="edit_source_id"
                type="number"
                value={formData.source_id}
                onChange={(e) => setFormData((prev) => ({ ...prev, source_id: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_status_id" className="text-right">
                Status
              </Label>
              <Select
                value={formData.status_id}
                onValueChange={(value) => setFormData((prev) => ({ ...prev, status_id: value }))}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1">New</SelectItem>
                  <SelectItem value="2">Contacted</SelectItem>
                  <SelectItem value="3">Qualified</SelectItem>
                  <SelectItem value="4">Converted</SelectItem>
                  <SelectItem value="5">Lost
</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_assigned_to" className="text-right">
                Assigned To
              </Label>
              <Select
                value={formData.assigned_to}
                onValueChange={(value) => setFormData((prev) => ({ ...prev, assigned_to: value }))}
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select user" />
                </SelectTrigger>
                <SelectContent>
                  {users.map((user) => (
                    <SelectItem key={user.id} value={user.id!.toString()}>
                      {user.username}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_notes" className="text-right">
                Notes
              </Label>
              <Textarea
                id="edit_notes"
                value={formData.notes}
                onChange={(e) => setFormData((prev) => ({ ...prev, notes: e.target.value }))}
                className="col-span-3"
                rows={3}
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="submit" onClick={handleEditLead} className="bg-cyan-600 hover:bg-cyan-700">
              Update Lead
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
