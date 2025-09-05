import axios from 'axios';

// ================== CONFIGURATION ==================
const API_URL = 'http://localhost:8080/api/v1';

const apiClient = axios.create({
  baseURL: API_URL,
});

// Axios Request Interceptor to automatically add the JWT token to every request.
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// ================== AUTH ==================
export interface LoginCredentials {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

export const login = async (credentials: LoginCredentials): Promise<LoginResponse> => {
  try {
    const response = await apiClient.post<LoginResponse>('/auth/login', credentials);
    return response.data;
  } catch (error) {
    console.error('Login failed:', error);
    if (axios.isAxiosError(error) && error.response) {
      throw new Error(error.response.data || 'Invalid username or password');
    }
    throw new Error('An unexpected error occurred during login.');
  }
};

// ================== CONTACTS ==================
export interface Contact {
  id: number;
  first_name: string;
  last_name: string;
  email?: string;
  primary_phone: string;
}

export type CreateContactPayload = Omit<Contact, 'id'>;

export const getContacts = async (): Promise<Contact[]> => {
  const response = await apiClient.get<Contact[]>('/contacts');
  return response.data;
};

export const createContact = async (contactData: CreateContactPayload): Promise<Contact> => {
  const response = await apiClient.post<Contact>('/contacts', contactData);
  return response.data;
};

// ================== USERS (for Select dropdowns) ==================
export interface UserSelectItem {
  id: number;
  username: string;
}
export const getUsers = async (): Promise<UserSelectItem[]> => {
  const response = await apiClient.get<UserSelectItem[]>('/users');
  return response.data;
};


// Add these to the CONTACTS section in apiService.ts

export type UpdateContactPayload = Omit<Contact, 'id'>;

export const updateContact = async (id: number, contactData: UpdateContactPayload): Promise<Contact> => {
  try {
    // Note the URL includes the contact's ID
    const response = await apiClient.put<Contact>(`/contacts/${id}`, contactData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update contact ${id}:`, error);
    throw error;
  }
};

export const deleteContact = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/contacts/${id}`);
  } catch (error) {
    console.error(`Failed to delete contact ${id}:`, error);
    throw error;
  }
};





// ================== PROPERTIES ==================
export interface Property {
  id: number;
  name: string;
  site_id: number;
  property_type_id: number;
  unit_no?: string;
  price: number;
  status: string;
}

export type CreatePropertyPayload = Omit<Property, 'id'>;

export const getProperties = async (): Promise<Property[]> => {
  const response = await apiClient.get<Property[]>('/properties');
  return response.data;
};

export const createProperty = async (propertyData: CreatePropertyPayload): Promise<Property> => {
  const response = await apiClient.post<Property>('/properties', propertyData);
  return response.data;
};

// ================== LEADS ==================
export interface Lead {
  id: number;
  contact_id: number;
  property_id?: number;
  source_id: number;
  status_id: number;
  assigned_to: number;
  notes?: string;
  created_at: string;
}

export type CreateLeadPayload = Omit<Lead, 'id' | 'created_at'>;

export const getLeads = async (): Promise<Lead[]> => {
  const response = await apiClient.get<Lead[]>('/leads');
  return response.data;
};

export const createLead = async (leadData: CreateLeadPayload): Promise<Lead> => {
  const response = await apiClient.post<Lead>('/leads', leadData);
  return response.data;
};


// Add these to the LEADS section in apiService.ts

export type UpdateLeadPayload = Omit<Lead, 'id' | 'created_at'>;

export const updateLead = async (id: number, leadData: UpdateLeadPayload): Promise<Lead> => {
  try {
    const response = await apiClient.put<Lead>(`/leads/${id}`, leadData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update lead ${id}:`, error);
    throw error;
  }
};

export const deleteLead = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/leads/${id}`);
  } catch (error) {
    console.error(`Failed to delete lead ${id}:`, error);
    throw error;
  }
};


// Add these to the PROPERTIES section in apiService.ts

export type UpdatePropertyPayload = Omit<Property, 'id'>;

export const updateProperty = async (id: number, propertyData: UpdatePropertyPayload): Promise<Property> => {
  try {
    const response = await apiClient.put<Property>(`/properties/${id}`, propertyData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update property ${id}:`, error);
    throw error;
  }
};

export const deleteProperty = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/properties/${id}`);
  } catch (error) {
    console.error(`Failed to delete property ${id}:`, error);
    throw error;
  }
};


// ================== DEALS ==================

// This interface now correctly matches the full Deal model from our Go backend.
export interface Deal {
  id: number;
  lead_id: number;
  property_id: number;
  stage_id: number;
  deal_status: string;
  deal_amount: number;
  deal_date: string;
  closing_date?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

// This is the corrected, explicit payload type for creating a new Deal.
// It includes 'notes' and other optional fields.
export interface CreateDealPayload {
  lead_id: number;
  property_id: number;
  stage_id: number;
  deal_status: string;
  deal_amount: number;
  closing_date?: string;
  notes?: string;
}

export const getDeals = async (): Promise<Deal[]> => {
  const response = await apiClient.get<Deal[]>('/deals');
  return response.data;
};

export const createDeal = async (dealData: CreateDealPayload): Promise<Deal> => {
  const response = await apiClient.post<Deal>('/deals', dealData);
  return response.data;
};

// Add these to the DEALS section in apiService.ts

export type UpdateDealPayload = CreateDealPayload; // The payload for creating and updating is the same

export const updateDeal = async (id: number, dealData: UpdateDealPayload): Promise<Deal> => {
  try {
    const response = await apiClient.put<Deal>(`/deals/${id}`, dealData);
    return response.data;
  } catch (error) {
    console.error(`Failed to update deal ${id}:`, error);
    throw error;
  }
};

export const deleteDeal = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/deals/${id}`);
  } catch (error) {
    console.error(`Failed to delete deal ${id}:`, error);
    throw error;
  }
};





// Add this new section to apiService.ts

// ================== TASKS ==================
export interface Task {
  id: number;
  task_name: string;
  task_description?: string;
  due_date: string; // Should be a string
  status: string;
  assigned_to: number; // Should be a number
}
export interface CreateTaskPayload {
  task_name: string;
  task_description?: string;
  due_date: string; // Must be an ISO string
  status: string;
  assigned_to: number;
}
export type UpdateTaskPayload = Omit<Task, 'id'>;

export const getTasks = async (): Promise<Task[]> => {
  const response = await apiClient.get<Task[]>('/tasks');
  return response.data;
};

export const createTask = async (taskData: CreateTaskPayload): Promise<Task> => {
  const response = await apiClient.post<Task>('/tasks', taskData);
  return response.data;
};

export const updateTask = async (id: number, taskData: UpdateTaskPayload): Promise<Task> => {
  const response = await apiClient.put<Task>(`/tasks/${id}`, taskData);
  return response.data;
};

export const deleteTask = async (id: number): Promise<void> => {
  await apiClient.delete(`/tasks/${id}`);
};






// Add this new section to apiService.ts

// ================== COMM LOGS ==================
export interface CommLog {
  id: number;
  contact_id: number;
  user_id: number;
  interaction_date: string;
  interaction_type: string;
  notes?: string;
}

export type CreateCommLogPayload = Omit<CommLog, 'id'>;
export type UpdateCommLogPayload = Omit<CommLog, 'id' | 'contact_id' | 'user_id'>; // Usually you only update notes/date/type

// Get all logs for ONE specific contact
export const getCommLogsForContact = async (contactId: number): Promise<CommLog[]> => {
  const response = await apiClient.get<CommLog[]>(`/contacts/${contactId}/comm-logs`);
  return response.data;
};

export const createCommLog = async (logData: CreateCommLogPayload): Promise<CommLog> => {
  const response = await apiClient.post<CommLog>('/comm-logs', logData);
  return response.data;
};

export const updateCommLog = async (id: number, logData: UpdateCommLogPayload): Promise<CommLog> => {
  const response = await apiClient.put<CommLog>(`/comm-logs/${id}`, logData);
  return response.data;
};

export const deleteCommLog = async (id: number): Promise<void> => {
  await apiClient.delete(`/comm-logs/${id}`);
};

// We also need a function to get a SINGLE contact's details
export const getContactById = async (id: number): Promise<Contact> => {
    const response = await apiClient.get<Contact>(`/contacts/${id}`);
    return response.data;
}




// ================== NOTES ==================
export interface Note {
  id: number;
  user_id: number;
  note_date: string;
  note_text: string;
}
export type CreateNotePayload = Omit<Note, 'id' | 'created_at' | 'note_date'>;
export type UpdateNotePayload = Omit<Note, 'id' | 'created_at' | 'user_id'>;

// Gets ALL notes (for a future admin dashboard, perhaps).
export const getNotes = async (): Promise<Note[]> => {
  try {
    const response = await apiClient.get<Note[]>('/notes');
    return response.data;
  } catch (error) { throw error; }
};

// Gets notes for a SPECIFIC user.
export const getNotesForUser = async (userId: number): Promise<Note[]> => {
  try {
    const response = await apiClient.get<Note[]>(`/users/${userId}/notes`);
    return response.data;
  } catch (error) { throw error; }
};

// THIS FUNCTION WAS INCOMPLETE
export const createNote = async (noteData: CreateNotePayload): Promise<Note> => {
  try {
    const response = await apiClient.post<Note>('/notes', noteData);
    return response.data;
  } catch (error) { throw error; }
};

// THIS FUNCTION WAS INCOMPLETE
export const updateNote = async (id: number, noteData: UpdateNotePayload): Promise<Note> => {
  try {
    const response = await apiClient.put<Note>(`/notes/${id}`, noteData);
    return response.data;
  } catch (error) { throw error; }
};

// THIS FUNCTION WAS INCOMPLETE
export const deleteNote = async (id: number): Promise<void> => {
  try {
    await apiClient.delete(`/notes/${id}`);
  } catch (error) { throw error; }
};









// ================== EVENTS ==================
export interface Event {
  id: number;
  event_name: string;
  event_description?: string;
  start_time: string;
  end_time: string;
  location?: string;
  organizer_id: number;
}

export type CreateEventPayload = Omit<Event, 'id'>;
export type UpdateEventPayload = Omit<Event, 'id'>;

// In src/services/apiService.ts


export const getEvents = async (): Promise<Event[]> => {
  try { // <-- Make sure this try/catch block is here
    const response = await apiClient.get<Event[]>('/events');
    return response.data;
  } catch (error) {
    console.error('Failed to fetch events:', error);
    throw error;
  }
};


export const createEvent = async (eventData: CreateEventPayload): Promise<Event> => {
  const response = await apiClient.post<Event>('/events', eventData);
  return response.data;
};

export const updateEvent = async (id: number, eventData: UpdateEventPayload): Promise<Event> => {
  const response = await apiClient.put<Event>(`/events/${id}`, eventData);
  return response.data;
};

export const deleteEvent = async (id: number): Promise<void> => {
  await apiClient.delete(`/events/${id}`);
};



// Add this new section to apiService.ts

// ================== REPORTS ==================
// Define the types for the report data structures
interface EmployeeLeadRow {
  employee_name: string;
  counts: { new: number; contacted: number; qualified: number; converted: number };
}
interface EmployeeLeadReport {
    rows: EmployeeLeadRow[];
}
export const getEmployeeLeadReport = async (): Promise<EmployeeLeadReport> => {
  const response = await apiClient.get<EmployeeLeadReport>('/reports/employee-leads');
  return response.data;
};

interface EmployeeSaleRow {
    employee_name: string;
    number_of_sales: number;
    total_sales_amount: number;
}
export const getEmployeeSalesReport = async (): Promise<EmployeeSaleRow[]> => {
    const response = await apiClient.get<EmployeeSaleRow[]>('/reports/employee-sales');
    return response.data;
}

// Add this to the REPORTS section in apiService.ts

export interface SourceSaleRow {
    source_name: string;
    number_of_sales: number;
    total_sales_amount: number;
}
export const getSourceSalesReport = async (): Promise<SourceSaleRow[]> => {
    const response = await apiClient.get<SourceSaleRow[]>('/reports/source-sales');
    return response.data;
}