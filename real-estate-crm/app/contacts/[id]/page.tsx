"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowLeft, User, Mail, Phone, Calendar, Edit } from "lucide-react"
import { getContact, type Contact } from "@/lib/api"
import ContactNotes from "@/components/contact-notes"
import { ContactCommLogs } from "@/components/contact-comm-logs"
import { ContactActivity } from "@/components/contact-activity"
import { PageBreadcrumb } from "@/components/page-breadcrumb"

export default function ContactDetailPage() {
  const params = useParams()
  const router = useRouter()
  const contactId = Number(params.id)

  const [contact, setContact] = useState<Contact | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")

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
          <Button variant="outline" size="sm">
            <Edit className="mr-2 h-4 w-4" />
            Edit Contact
          </Button>
          <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
            Active
          </Badge>
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
          <TabsTrigger value="activity">Activity</TabsTrigger>
        </TabsList>

        <TabsContent value="notes">
          <ContactNotes contactId={contactId} />
        </TabsContent>

        <TabsContent value="communications">
          <ContactCommLogs contactId={contactId} />
        </TabsContent>

        <TabsContent value="activity">
          <ContactActivity contactId={contactId} />
        </TabsContent>
      </Tabs>
    </div>
  )
}
