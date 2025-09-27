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
import { MoreHorizontal, Search, Plus, Edit, Trash2, Building, DollarSign, MapPin } from "lucide-react"
import { api, type Property } from "@/lib/api"
import { useAuth, ROLE_RECEPTION } from "@/lib/auth" // Import useAuth and role constants

export function PropertiesManagement() {
  const [properties, setProperties] = useState<Property[]>([])
  const [searchTerm, setSearchTerm] = useState("")
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [selectedProperty, setSelectedProperty] = useState<Property | null>(null)
  const [formData, setFormData] = useState({
    name: "",
    site_id: "",
    property_type_id: "",
    unit_no: "",
    price: "",
    status: "Available" as "Available" | "Pending" | "Sold",
  })

  const { hasRole } = useAuth() // Use the useAuth hook

  // Load properties on component mount
  useEffect(() => {
    loadProperties()
  }, [])

  const loadProperties = async () => {
    try {
      setIsLoading(true)
      const data = await api.getProperties()
      setProperties(data)
    } catch (err) {
      setError("Failed to load properties")
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateProperty = async () => {
    try {
      const propertyData = {
        name: formData.name,
        site_id: Number.parseInt(formData.site_id),
        property_type_id: Number.parseInt(formData.property_type_id),
        unit_no: formData.unit_no,
        price: Number.parseFloat(formData.price),
        status: formData.status,
      }
      await api.createProperty(propertyData)
      setIsCreateDialogOpen(false)
      resetForm()
      loadProperties()
    } catch (err) {
      setError("Failed to create property")
    }
  }

  const handleEditProperty = async () => {
    if (!selectedProperty?.id) return

    try {
      const propertyData = {
        name: formData.name,
        site_id: Number.parseInt(formData.site_id),
        property_type_id: Number.parseInt(formData.property_type_id),
        unit_no: formData.unit_no,
        price: Number.parseFloat(formData.price),
        status: formData.status,
      }
      await api.updateProperty(selectedProperty.id, propertyData)
      setIsEditDialogOpen(false)
      resetForm()
      setSelectedProperty(null)
      loadProperties()
    } catch (err) {
      console.error("Error updating property:", err)
      setError("Failed to update property")
    }
  }

  const handleDeleteProperty = async (id: number) => {
    if (!confirm("Are you sure you want to delete this property?")) return

    try {
      await api.deleteProperty(id)
      loadProperties()
    } catch (err) {
      console.error("Error deleting property:", err)
      setError("Failed to delete property")
    }
  }

  const resetForm = () => {
    setFormData({
      name: "",
      site_id: "",
      property_type_id: "",
      unit_no: "",
      price: "",
      status: "Available",
    })
  }

  const openEditDialog = (property: Property) => {
    setSelectedProperty(property)
    setFormData({
      name: property.name,
      site_id: property.site_id.toString(),
      property_type_id: property.property_type_id.toString(),
      unit_no: property.unit_no,
      price: property.price.toString(),
      status: property.status,
    })
    setIsEditDialogOpen(true)
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "Available":
        return (
          <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
            Available
          </Badge>
        )
      case "Pending":
        return (
          <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">
            Pending
          </Badge>
        )
      case "Sold":
        return (
          <Badge variant="outline" className="bg-red-50 text-red-700 border-red-200">
            Sold
          </Badge>
        )
      default:
        return <Badge variant="outline">{status}</Badge>
    }
  }

  const filteredProperties = properties.filter(
    (property) =>
      property.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      property.unit_no.toLowerCase().includes(searchTerm.toLowerCase()) ||
      property.status.toLowerCase().includes(searchTerm.toLowerCase()),
  )

  const availableCount = properties.filter((p) => p.status === "Available").length
  const pendingCount = properties.filter((p) => p.status === "Pending").length
  const soldCount = properties.filter((p) => p.status === "Sold").length

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Loading properties...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold text-balance text-foreground">Properties Management</h1>
          <p className="text-muted-foreground text-pretty">Manage your real estate property listings and inventory</p>
        </div>
        {hasRole(ROLE_RECEPTION) && (
          <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button className="bg-cyan-600 hover:bg-cyan-700">
                <Plus className="mr-2 h-4 w-4" />
                Add Property
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>Create New Property</DialogTitle>
                <DialogDescription>Add a new property to your inventory.</DialogDescription>
              </DialogHeader>
              <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="name" className="text-right">
                    Name
                  </Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData((prev) => ({ ...prev, name: e.target.value }))}
                    className="col-span-3"
                    placeholder="Property name"
                    required
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="site_id" className="text-right">
                    Site ID
                  </Label>
                  <Input
                    id="site_id"
                    type="number"
                    value={formData.site_id}
                    onChange={(e) => setFormData((prev) => ({ ...prev, site_id: e.target.value }))}
                    className="col-span-3"
                    placeholder="Site identifier"
                    required
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="property_type_id" className="text-right">
                    Type ID
                  </Label>
                  <Input
                    id="property_type_id"
                    type="number"
                    value={formData.property_type_id}
                    onChange={(e) => setFormData((prev) => ({ ...prev, property_type_id: e.target.value }))}
                    className="col-span-3"
                    placeholder="Property type ID"
                    required
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="unit_no" className="text-right">
                    Unit No
                  </Label>
                  <Input
                    id="unit_no"
                    value={formData.unit_no}
                    onChange={(e) => setFormData((prev) => ({ ...prev, unit_no: e.target.value }))}
                    className="col-span-3"
                    placeholder="Unit number"
                    required
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="price" className="text-right">
                    Price
                  </Label>
                  <Input
                    id="price"
                    type="number"
                    step="0.01"
                    value={formData.price}
                    onChange={(e) => setFormData((prev) => ({ ...prev, price: e.target.value }))}
                    className="col-span-3"
                    placeholder="Property price"
                    required
                  />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="status" className="text-right">
                    Status
                  </Label>
                  <Select
                    value={formData.status}
                    onValueChange={(value: "Available" | "Pending" | "Sold") =>
                      setFormData((prev) => ({ ...prev, status: value }))
                    }
                  >
                    <SelectTrigger className="col-span-3">
                      <SelectValue placeholder="Select status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Available">Available</SelectItem>
                      <SelectItem value="Pending">Pending</SelectItem>
                      <SelectItem value="Sold">Sold</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <DialogFooter>
                <Button type="submit" onClick={handleCreateProperty} className="bg-cyan-600 hover:bg-cyan-700">
                  Create Property
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
            <Building className="h-5 w-5 text-cyan-600" />
            <h3 className="font-semibold">Total Properties</h3>
          </div>
          <p className="text-2xl font-bold text-cyan-600 mt-2">{properties.length}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <MapPin className="h-5 w-5 text-green-600" />
            <h3 className="font-semibold">Available</h3>
          </div>
          <p className="text-2xl font-bold text-green-600 mt-2">{availableCount}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <DollarSign className="h-5 w-5 text-yellow-600" />
            <h3 className="font-semibold">Pending</h3>
          </div>
          <p className="text-2xl font-bold text-yellow-600 mt-2">{pendingCount}</p>
        </div>
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2">
            <Building className="h-5 w-5 text-red-600" />
            <h3 className="font-semibold">Sold</h3>
          </div>
          <p className="text-2xl font-bold text-red-600 mt-2">{soldCount}</p>
        </div>
      </div>

      {/* Search */}
      <div className="relative max-w-sm">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search properties..."
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
              <TableHead className="text-card-foreground">Property Name</TableHead>
              <TableHead className="text-card-foreground">Unit No</TableHead>
              <TableHead className="text-card-foreground">Price</TableHead>
              <TableHead className="text-card-foreground">Status</TableHead>
              <TableHead className="text-card-foreground">Site ID</TableHead>
              <TableHead className="w-[50px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredProperties.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                  {searchTerm
                    ? "No properties found matching your search."
                    : "No properties yet. Create your first property!"}
                </TableCell>
              </TableRow>
            ) : (
              filteredProperties.map((property) => (
                <TableRow key={property.id} className="hover:bg-muted/50">
                  <TableCell className="font-medium text-card-foreground">{property.name}</TableCell>
                  <TableCell className="text-card-foreground">{property.unit_no}</TableCell>
                  <TableCell className="text-card-foreground">${property.price.toLocaleString()}</TableCell>
                  <TableCell>{getStatusBadge(property.status)}</TableCell>
                  <TableCell className="text-card-foreground">{property.site_id}</TableCell>
                  <TableCell>
                    {hasRole(ROLE_RECEPTION) && (
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="sm">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => openEditDialog(property)}>
                            <Edit className="mr-2 h-4 w-4" />
                            Edit
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            className="text-destructive"
                            onClick={() => property.id && handleDeleteProperty(property.id)}
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
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Edit Property</DialogTitle>
            <DialogDescription>Update the property information.</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_name" className="text-right">
                Name
              </Label>
              <Input
                id="edit_name"
                value={formData.name}
                onChange={(e) => setFormData((prev) => ({ ...prev, name: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_site_id" className="text-right">
                Site ID
              </Label>
              <Input
                id="edit_site_id"
                type="number"
                value={formData.site_id}
                onChange={(e) => setFormData((prev) => ({ ...prev, site_id: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_property_type_id" className="text-right">
                Type ID
              </Label>
              <Input
                id="edit_property_type_id"
                type="number"
                value={formData.property_type_id}
                onChange={(e) => setFormData((prev) => ({ ...prev, property_type_id: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_unit_no" className="text-right">
                Unit No
              </Label>
              <Input
                id="edit_unit_no"
                value={formData.unit_no}
                onChange={(e) => setFormData((prev) => ({ ...prev, unit_no: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_price" className="text-right">
                Price
              </Label>
              <Input
                id="edit_price"
                type="number"
                step="0.01"
                value={formData.price}
                onChange={(e) => setFormData((prev) => ({ ...prev, price: e.target.value }))}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="edit_status" className="text-right">
                Status
              </Label>
              <Select
                value={formData.status}
                onValueChange={(value: "Available" | "Pending" | "Sold") =>
                  setFormData((prev) => ({ ...prev, status: value }))
                }
              >
                <SelectTrigger className="col-span-3">
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Available">Available</SelectItem>
                  <SelectItem value="Pending">Pending</SelectItem>
                  <SelectItem value="Sold">Sold</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button type="submit" onClick={handleEditProperty} className="bg-cyan-600 hover:bg-cyan-700">
              Update Property
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
