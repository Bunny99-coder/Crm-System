"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { ArrowLeft, User, Mail, Phone, Calendar, Edit, Trash2 } from "lucide-react"
import { api, getContact, type Contact } from "@/lib/api"
import ContactNotes from "@/components/contact-notes"
import { ContactCommLogs } from "@/components/contact-comm-logs"
import { PageBreadcrumb } from "@/components/page-breadcrumb"

export default function ContactDetailPage() {
  const params = useParams()
  const router = useRouter()
  const contactId = Number(params.id)

  const [contact, setContact] = useState<Contact | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
  const [formData, setFormData] = useState({
    first_name: "",
    last_name: "",
    email: "",
    primary_phone: "",
  })

  useEffect(() => {
    if (contact) {
      setFormData({
        first_name: contact.first_name,
        last_name: contact.last_name,
        email: contact.email,
        primary_phone: contact.primary_phone,
      })
    }
  }, [contact])

  const handleEditContact = async () => {
    if (!contact?.id) return

    try {
      await api.updateContact(contact.id, formData)
      setIsEditDialogOpen(false)
      loadContactDetails() // Refresh contact details
    } catch (err) {
      setError("Failed to update contact")
    }
  }

  const handleDeleteContact = async () => {
    if (!contact?.id) return

    try {
      await api.deleteContact(contact.id)
      router.push("/contacts") // Redirect to contacts list
    } catch (err) {
      setError("Failed to delete contact")
    }
  }

  const openEditDialog = () => {
    if (!contact) return
    setFormData({
      first_name: contact.first_name,
      last_name: contact.last_name,
      email: contact.email,
      primary_phone: contact.primary_phone,
    })
    setIsEditDialogOpen(true)
  }

  useEffect(() => {
    loadContactDetails()
  }, [contactId])

  const loadContactDetails = async () => {
    try {
      setIsLoading(true)
      const contactData = await getContact(contactId)
      if (!contactData) {
        setError("Contact not found")
      } else {
        setContact(contactData)
      }
    } catch (err) {
      setError("Failed to load contact details")
    } finally {
      setIsLoading(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Loading contact details...</div>
      </div>
    )
  }

  if (error || !contact) {
    return (
      <div className="flex flex-col items-center justify-center h-64 space-y-4">
        <div className="text-destructive">{error || "Contact not found"}</div>
        <Button onClick={() => router.push("/contacts")} variant="outline">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Contacts
        </Button>
      </div>
    )
  }

  const breadcrumbItems = [
    { label: "Contacts", href: "/contacts" },
    { label: `${contact.first_name} ${contact.last_name}` },
  ]

  return (
    <div className="space-y-6">
      <PageBreadcrumb items={breadcrumbItems} />

      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Button onClick={() => router.push("/contacts")} variant="outline" size="sm">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-balance">
              {contact.first_name} {contact.last_name}
            </h1>
            <p className="text-muted-foreground">Contact details and communication history</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={openEditDialog}>
            <Edit className="mr-2 h-4 w-4" />
            Edit Contact
          </Button>
          <Button
            variant="destructive"
            size="sm"
            onClick={() => setIsDeleteDialogOpen(true)}
          >
            <Trash2 className="mr-2 h-4 w-4" />
            Delete Contact
          </Button>
        </div>
      </div>

      {/* Contact Overview */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Contact ID</CardTitle>
            <User className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-cyan-600">#{contact.id}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Email</CardTitle>
            <Mail className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-lg font-medium text-orange-600 truncate">{contact.email}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Phone</CardTitle>
            <Phone className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-lg font-medium text-purple-600">{contact.primary_phone}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Created</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-lg font-medium text-green-600">
              {contact.created_at ? new Date(contact.created_at).toLocaleDateString() : "N/A"}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabbed Content */}
      <Tabs defaultValue="notes" className="space-y-4">
        <TabsList>
          <TabsTrigger value="notes">Notes</TabsTrigger>
          <TabsTrigger value="communications">Communications</TabsTrigger>
        </TabsList>

        <TabsContent value="notes">
          <ContactNotes contactId={contactId} />
        </TabsContent>

        <TabsContent value="communications">
          <ContactCommLogs contactId={contactId} />
        </TabsContent>
      </Tabs>

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

      {/* Delete Confirmation Dialog */}
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Are you sure?</DialogTitle>
            <DialogDescription>
              This action cannot be undone. This will permanently delete the contact.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDeleteDialogOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDeleteContact}>
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
