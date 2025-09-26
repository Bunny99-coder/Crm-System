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
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { MoreHorizontal, Search, Plus, Edit, Trash2, User, Mail, Phone, Eye } from "lucide-react"
import { api, type Contact } from "@/lib/api"
import { useAuth, ROLE_SALES_AGENT, ROLE_RECEPTION } from "@/lib/auth" // Import useAuth and role constants

export function ContactsManagement() {
  const [contacts, setContacts] = useState<Contact[]>([])
  const [searchTerm, setSearchTerm] = useState("")
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [selectedContact, setSelectedContact] = useState<Contact | null>(null)
  const [formData, setFormData] = useState({
    first_name: "",
    last_name: "",
    email: "",
    primary_phone: "",
  })

  const { user, hasRole } = useAuth() // Use the useAuth hook

  // Load contacts on component mount
  useEffect(() => {
    loadContacts()
  }, [])

  const loadContacts = async () => {
    try {
      setIsLoading(true)
      const data = await api.getContacts()
      setContacts(data)
    } catch (err) {
      setError("Failed to load contacts")
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateContact = async () => {
    try {
      await api.createContact(formData)
      setIsCreateDialogOpen(false)
      resetForm()
      loadContacts()
    } catch (err) {
      setError("Failed to create contact")
    }
  }

  const handleEditContact = async () => {
    if (!selectedContact?.id) return

    try {
      await api.updateContact(selectedContact.id, formData)
      setIsEditDialogOpen(false)
      resetForm()
      setSelectedContact(null)
      loadContacts()
    } catch (err) {
      setError("Failed to update contact")
    }
  }

  const handleDeleteContact = async (id: number) => {
    if (!confirm("Are you sure you want to delete this contact?")) return

    try {
      await api.deleteContact(id)
      loadContacts()
    } catch (err) {
      setError("Failed to delete contact")
    }
  }

  const resetForm = () => {
    setFormData({
      first_name: "",
      last_name: "",
      email: "",
      primary_phone: "",
    })
  }

  const openEditDialog = (contact: Contact) => {
    setSelectedContact(contact)
    setFormData({
      first_name: contact.first_name,
      last_name: contact.last_name,
      email: contact.email,
      primary_phone: contact.primary_phone,
    })
    setIsEditDialogOpen(true)
  }

  const filteredContacts = contacts.filter(
    (contact) =>
      contact.first_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      contact.last_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      contact.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
      contact.primary_phone.includes(searchTerm),
  )

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Loading contacts...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold text-balance text-foreground">Contacts Management</h1>
          <p className="text-muted-foreground text-pretty">Manage your client contacts and their information</p>
        </div>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button className="bg-cyan-600 hover:bg-cyan-700">
              <Plus className="mr-2 h-4 w-4" />
              Add Contact
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Create New Contact</DialogTitle>
              <DialogDescription>Add a new contact to your CRM system.</DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="first_name" className="text-right">
                  First Name
                </Label>
                <Input
                  id="first_name"
                  value={formData.first_name}
                  onChange={(e) => setFormData((prev) => ({ ...prev, first_name: e.target.value }))}
                  className="col-span-3"
                  required
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="last_name" className="text-right">
                  Last Name
                </Label>
                <Input
                  id="last_name"
                  value={formData.last_name}
                  onChange={(e) => setFormData((prev) => ({ ...prev, last_name: e.target.value }))}
                  className="col-span-3"
                  required
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="email" className="text-right">
                  Email
                </Label>
                <Input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => setFormData((prev) => ({ ...prev, email: e.target.value }))}
                  className="col-span-3"
                  required
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="primary_phone" className="text-right">
                  Phone
                </Label>
                <Input
                  id="primary_phone"
                  value={formData.primary_phone}
                  onChange={(e) => setFormData((prev) => ({ ...prev, primary_phone: e.target.value }))}
                  className="col-span-3"
                  required
                />
              </div>
            </div>
            <DialogFooter>
              <Button type="submit" onClick={handleCreateContact} className="bg-cyan-600 hover:bg-cyan-700">
                Create Contact
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <User className="h-5 w-5 text-cyan-600" />
            <h3 className="font-semibold">Total Contacts</h3>
          </div>
          <p className="text-2xl font-bold text-cyan-600 mt-2">{contacts.length}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <Mail className="h-5 w-5 text-orange-600" />
            <h3 className="font-semibold">Active Contacts</h3>
          </div>
          <p className="text-2xl font-bold text-orange-600 mt-2">{contacts.length}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <Phone className="h-5 w-5 text-green-600" />
            <h3 className="font-semibold">Recent Contacts</h3>
          </div>
          <p className="text-2xl font-bold text-green-600 mt-2">{Math.min(contacts.length, 5)}</p>
        </div>
      </div>

      {/* Search */}
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search contacts..."
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
              <TableHead className="text-card-foreground">Name</TableHead>
              <TableHead className="text-card-foreground">Email</TableHead>
              <TableHead className="text-card-foreground">Phone</TableHead>
              <TableHead className="text-card-foreground">Status</TableHead>
              <TableHead className="w-[50px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredContacts.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                  {searchTerm
                    ? "No contacts found matching your search."
                    : "No contacts yet. Create your first contact!"}
                </TableCell>
              </TableRow>
            ) : (
              filteredContacts.map((contact) => {
                const canEditDelete = hasRole(ROLE_RECEPTION) || (user && contact.created_by === user.id);
                return (
                  <TableRow key={contact.id} className="hover:bg-muted/50">
                    <TableCell className="font-medium text-card-foreground">
                      {contact.first_name} {contact.last_name}
                    </TableCell>
                    <TableCell className="text-card-foreground">{contact.email}</TableCell>
                    <TableCell className="text-card-foreground">{contact.primary_phone}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
                        Active
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => (window.location.href = `/contacts/${contact.id}`)}>
                            <Eye className="mr-2 h-4 w-4" />
                            View Details
                          </DropdownMenuItem>
                          {canEditDelete && (
                            <DropdownMenuItem onClick={() => openEditDialog(contact)}>
                              <Edit className="mr-2 h-4 w-4" />
                              Edit
                            </DropdownMenuItem>
                          )}
                          {canEditDelete && (
                            <DropdownMenuItem
                              className="text-destructive"
                              onClick={() => contact.id && handleDeleteContact(contact.id)}
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
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Edit Contact</DialogTitle>
            <DialogDescription>Update the contact information.</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_first_name" className="text-right">
                First Name
              </Label>
              <Input
                id="edit_first_name"
                value={formData.first_name}
                onChange={(e) => setFormData((prev) => ({ ...prev, first_name: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_last_name" className="text-right">
                Last Name
              </Label>
              <Input
                id="edit_last_name"
                value={formData.last_name}
                onChange={(e) => setFormData((prev) => ({ ...prev, last_name: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_email" className="text-right">
                Email
              </Label>
              <Input
                id="edit_email"
                type="email"
                value={formData.email}
                onChange={(e) => setFormData((prev) => ({ ...prev, email: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_primary_phone" className="text-right">
                Phone
              </Label>
              <Input
                id="edit_primary_phone"
                value={formData.primary_phone}
                onChange={(e) => setFormData((prev) => ({ ...prev, primary_phone: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="submit" onClick={handleEditContact} className="bg-cyan-600 hover:bg-cyan-700">
              Update Contact
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
