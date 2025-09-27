// lib/api.ts
import { authManager } from "./auth"

const API_BASE_URL = "http://localhost:8081/api/v1"

// ==========================
// Types (OpenAPI Spec based)
// ==========================
export interface Contact {
  id?: number
  first_name: string
  last_name: string
  email: string
  primary_phone: string
  created_at?: string
  updated_at?: string
  created_by?: number // Add this line
}

export interface Property {
  id?: number
  name: string
  site_id: number
  property_type_id: number
  unit_no: string
  price: number
  status: "Available" | "Pending" | "Sold"
}

export interface Lead {
  id?: number
  contact_id: number
  property_id: number
  source_id: number
  status_id: number
  assigned_to: number
  notes: string
}

export interface Deal {
  id?: number
  lead_id: number
  property_id: number
  stage_id: number
  deal_status: "Pending" | "Closed-Won" | "Closed-Lost"
  deal_amount: number
  created_by?: number // Add this line
}

export interface Task {
  id?: number
  task_name: string
  task_description: string
  due_date: string
  status: "Pending" | "Completed"
  assigned_to: number
}

export interface User {
  id?: number
  username: string
  email: string
  role_id: number
}

// Reports
export interface EmployeeLeadReportRow {
  employee_id: number
  employee_name: string
  counts: {
    new: number
    contacted: number
    qualified: number
    converted: number
    lost: number
  }
}

export interface EmployeeLeadReport {
  rows: EmployeeLeadReportRow[]
  total: {
    new: number
    contacted: number
    qualified: number
    converted: number
    lost: number
  }
}

export interface SourceLeadReportRow {
  lead_date: string
  contact_name: string
  contact_phone: string
  contact_email: string
  lead_source: string
  assigned_employee: string
  lead_status: string
}

export interface SourceSalesReportRow {
  source_name: string
  number_of_sales: number
  total_sales_amount: number
}

export interface EmployeeSalesReportRow {
  employee_id: number
  employee_name: string
  number_of_sales: number
  total_sales_amount: number
}

export interface DealsPipelineReportRow {
  stage_name: string
  deal_count: number
  total_amount: number
  avg_days_in_stage: number
}

// Activity & subtypes
export interface Activity {
  id?: number
  entity_type: "contact" | "deal"
  entity_id: number
  activity_type: "note" | "task" | "event" | "comm_log"
  title: string
  description?: string
  due_date?: string
  status?: "Pending" | "Completed"
  communication_type?: "Email" | "Call" | "Meeting" | "SMS"
  created_at?: string
  created_by?: number
}

export interface Note {
  id: number
  content: string
  created_at?: string
  updated_at?: string
  created_by: number
}


export interface Event {
  id?: number
  event_name: string
  event_description?: string
  event_date: string
  created_by: number
}

export interface CommLog {
  id?: number
  contact_id?: number
  user_id?: number
  lead_id?: number
  deal_id?: number
  interaction_date: string
  interaction_type: "Call" | "Email" | "Meeting" | "SMS" | "Other"
  notes?: string
  created_at?: string
  deleted_at?: string
}

// ==========================
// API Client
// ==========================





class ApiClient {
  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    console.log('Requesting endpoint:', endpoint, 'with options:', options);
    console.log("API_BASE_URL:", API_BASE_URL)
    const url = `${API_BASE_URL}${endpoint}`
    console.log("Constructed URL:", url)

    const headers: Record<string, string> = {
      ...authManager.getAuthHeaders(),
      ...(options.headers as Record<string, string>),
    }

    const config: RequestInit = {
      ...options,
      headers,
    }

    try {
      const response = await fetch(url, config)

      if (response.status === 401) {
        authManager.clearAuth()
        window.location.href = "/login"
        throw new Error("Unauthorized - please login again")
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const contentType = response.headers.get("content-type")
      if (contentType && contentType.includes("application/json")) {
        return response.json()
      }

      return {} as T
    } catch (error) {
      console.error(`API request failed for ${endpoint}:`, error)
      throw error
    }
  }

  // ========== Auth ==========
  async login(username: string, password: string): Promise<{ token: string; user: User }> {
    return this.request("/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    })
  }

  async registerUser(userData: { username: string; email: string; password: string; role_id: number }): Promise<void> {
    return this.request("/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(userData),
    })
  }


// Contact Notes Methods
  // ==========================
  getContactNotes(contactId: number) {
    return this.request<Note[]>(`/contacts/${contactId}/notes`)
  }

  createContactNote(contactId: number, note: { content: string }) {
    return this.request<Note>(`/contacts/${contactId}/notes`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(note),
    })
  }

  updateContactNote(contactId: number, noteId: number, note: { content: string }) {
    return this.request<Note>(`/contacts/${contactId}/notes/${noteId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(note),
    })
  }

  deleteContactNote(contactId: number, noteId: number) {
    return this.request(`/contacts/${contactId}/notes/${noteId}`, { method: "DELETE" })
  }








   getNotesForDeal(dealId: number) {
    return this.request<Note[]>(`/deals/${dealId}/notes`)
  }


createNoteForDeal(dealId: number, note: { content: string }) {
  return this.request<Note>(`/deals/${dealId}/notes`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(note),
  })
}


  updateNoteForDeal(dealId: number, noteId: number, note: { content: string }) {
    return this.request<Note>(`/deals/${dealId}/notes/${noteId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(note),
    })
  }

  deleteNoteForDeal(dealId: number, noteId: number) {
    return this.request(`/deals/${dealId}/notes/${noteId}`, { method: "DELETE" })
  }

  getEventsForDeal(dealId: number) {
    return this.request<Event[]>(`/deals/${dealId}/events`)
  }

  createEventForDeal(dealId: number, event: Omit<Event, "id">) {
    return this.request<Event>(`/deals/${dealId}/events`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(event),
    })
  }

  getTasksForDeal(dealId: number) {
    return this.request<Task[]>(`/deals/${dealId}/tasks`)
  }

  createTaskForDeal(dealId: number, task: Omit<Task, "id">) {
    return this.request<Task>(`/deals/${dealId}/tasks`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(task),
    })
  }

  updateTaskForDeal(dealId: number, taskId: number, task: Partial<Omit<Task, "id">>) {
    return this.request<Task>(`/deals/${dealId}/tasks/${taskId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(task),
    })
  }







  // ========== Contacts ==========
  getContacts() { return this.request<Contact[]>("/contacts") }
  getContactById(id: number) { return this.request<Contact>(`/contacts/${id}`) }
  createContact(contact: Omit<Contact, "id">) {
    return this.request<Contact>("/contacts", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(contact),
    })
  }
  updateContact(id: number, contact: Omit<Contact, "id">) {
    return this.request(`/contacts/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(contact),
    })
  }
  deleteContact(id: number) { return this.request(`/contacts/${id}`, { method: "DELETE" }) }






  // ========== Properties ==========
  getProperties() { return this.request<Property[]>("/properties") }
  getPropertyById(id: number) { return this.request<Property>(`/properties/${id}`) }
  createProperty(property: Omit<Property, "id">) {
    return this.request<Property>("/properties", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(property),
    })
  }
  updateProperty(id: number, property: Omit<Property, "id">) {
    return this.request(`/properties/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(property),
    })
  }
  deleteProperty(id: number) { return this.request(`/properties/${id}`, { method: "DELETE" }) }

  // ========== Users ==========
  getUsers() { return this.request<User[]>("/users") }
  getUserById(id: number) { return this.request<User>(`/users/${id}`) }

  // ========== Leads ==========
  async getLeads(assignedToUserId?: number): Promise<Lead[]> {
    let url = "/leads"
    if (assignedToUserId) {
      url += `?assigned_to=${assignedToUserId}`
    }
    return this.request<Lead[]>(url)
  }
  getLeadById(id: number) { return this.request<Lead>(`/leads/${id}`) }
  createLead(lead: Omit<Lead, "id">) {
    return this.request<Lead>("/leads", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(lead),
    })
  }
  updateLead(id: number, lead: Omit<Lead, "id">) {
    return this.request(`/leads/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(lead),
    })
  }
  deleteLead(id: number) { return this.request(`/leads/${id}`, { method: "DELETE" }) }

  // ========== Deals ==========
  getDeals() { return this.request<Deal[]>("/deals") }
  getDealById(id: number) { return this.request<Deal>(`/deals/${id}`) }
  createDeal(deal: Omit<Deal, "id">) {
    return this.request<Deal>("/deals", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(deal),
    })
  }
  updateDeal(id: number, deal: Omit<Deal, "id">) {
    return this.request(`/deals/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(deal),
    })
  }
  deleteDeal(id: number) { return this.request(`/deals/${id}`, { method: "DELETE" }) }
  getDealActivity(dealId: number) { return this.request<Activity[]>(`/deals/${dealId}/activity`) }
  getContactActivity(contactId: number) { return this.request<Activity[]>(`/contacts/${contactId}/activity`) }

  // ========== Tasks ==========
  async getTasks(assignedToUserId?: number): Promise<Task[]> {
    let url = "/tasks"
    if (assignedToUserId) {
      url += `?assigned_to=${assignedToUserId}`
    }
    return this.request<Task[]>(url)
  }
  getTaskById(id: number) { return this.request<Task>(`/tasks/${id}`) }
  createTask(task: Omit<Task, "id">) {
    return this.request<Task>("/tasks", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(task),
    })
  }
  updateTask(id: number, task: Omit<Task, "id">) {
    return this.request(`/tasks/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(task),
    })
  }
  deleteTask(id: number) { return this.request(`/tasks/${id}`, { method: "DELETE" }) }

  // ========== Notes (global) ==========
  getNotes() { return this.request<Note[]>("/notes") }
  getNoteById(id: number) { return this.request<Note>(`/notes/${id}`) }
  createNote(note: Omit<Note, "id">) {
    return this.request<Note>("/notes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(note),
    })
  }
  updateNote(id: number, note: Omit<Note, "id">) {
    return this.request(`/notes/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(note),
    })
  }
  deleteNote(id: number) { return this.request(`/notes/${id}`, { method: "DELETE" }) }

  // ========== Events ==========
  getEvents() { return this.request<Event[]>("/events") }
  getEventById(id: number) { return this.request<Event>(`/events/${id}`) }
  createEvent(event: Omit<Event, "id">) {
    return this.request<Event>("/events", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(event),
    })
  }
  updateEvent(id: number, event: Omit<Event, "id">) {
    return this.request(`/events/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(event),
    })
  }
  deleteEvent(id: number) { return this.request(`/events/${id}`, { method: "DELETE" }) }

  // ========== Communication Logs ==========
  createCommLog(commLog: Omit<CommLog, "id">) {
    return this.request<CommLog>("/comm-logs", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(commLog),
    })
  }
  getCommLogById(id: number) { return this.request<CommLog>(`/comm-logs/${id}`) }
  updateCommLog(id: number, commLog: Omit<CommLog, "id">) {
    return this.request(`/comm-logs/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(commLog),
    })
  }
  deleteCommLog(id: number) { return this.request(`/comm-logs/${id}`, { method: "DELETE" }) }
  getCommLogsForContact(contactId: number) { return this.request<CommLog[]>(`/contacts/${contactId}/comm-logs`) }

  createContactCommLog(contactId: number, commLog: Omit<CommLog, "id" | "contact_id">) {
    return this.request<CommLog>(`/contacts/${contactId}/comm-logs`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(commLog),
    })
  }

  updateContactCommLog(contactId: number, logId: number, commLog: Partial<Omit<CommLog, "id" | "contact_id">>) {
    return this.request<CommLog>(`/contacts/${contactId}/comm-logs/${logId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(commLog),
    })
  }
  deleteContactCommLog(contactId: number, logId: number) {
    return this.request(`/contacts/${contactId}/comm-logs/${logId}`, { method: "DELETE" })
  }

  // ========== User-specific Notes/Events ==========
  getNotesForUser(userId: number) { return this.request<Note[]>(`/users/${userId}/notes`) }
  getEventsForUser(userId: number) { return this.request<Event[]>(`/users/${userId}/events`) }

  // ========== Reports ==========
  getEmployeeLeadReport() { return this.request<EmployeeLeadReport>("/reports/employee-leads") }
  getSourceLeadReport() { return this.request<SourceLeadReportRow[]>("/reports/source-leads") }
  getEmployeeSalesReport() { return this.request<EmployeeSalesReportRow[]>("/reports/employee-sales") }
  getSourceSalesReport() { return this.request<SourceSalesReportRow[]>("/reports/source-sales") }

  getMySalesReport() { return this.request<EmployeeSalesReportRow[]>("/reports/my-sales") }

  getDealsPipelineReport(period: string) { return this.request<DealsPipelineReportRow[]>(`/reports/deals-pipeline?period=${period}`) }
}

export const api = new ApiClient()

// Helper: get a contact by ID using the api instance
export const getContact = async (id: number) => {
  try {
    return await api.getContactById(id)
  } catch (error) {
    console.error("Failed to fetch contact:", error)
    return null
  }
}
