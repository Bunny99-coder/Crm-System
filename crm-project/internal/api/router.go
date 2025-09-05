// Replace the entire contents of your internal/api/router.go file with this.
package api

import (
	"crm-project/internal/api/handlers"
	"crm-project/internal/dto"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)
func NewRouter(jwtSecret string, authHandler *handlers.AuthHandler, contactHandler *handlers.ContactHandler, userHandler *handlers.UserHandler, propertyHandler *handlers.PropertyHandler, leadHandler *handlers.LeadHandler, dealHandler *handlers.DealHandler, reportHandler *handlers.ReportHandler, taskHandler *handlers.TaskHandler, commLogHandler *handlers.CommLogHandler, noteHandler *handlers.NoteHandler, eventHandler *handlers.EventHandler) http.Handler {
	r := chi.NewRouter()

	// --- Global Middleware ---
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(TimeoutMiddleware(30 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		// --- Public Routes ---
		r.Post("/auth/login", authHandler.Login)
		r.Post("/users", userHandler.CreateUser)

		// --- Protected Routes ---
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(jwtSecret))

			// --- Routes for ANY authenticated user ---
			// Everyone can VIEW lists of core data. Service layer filters for Sales Agents.
			r.Get("/contacts", contactHandler.GetAllContacts)
			r.Get("/contacts/{contactId}", contactHandler.GetContactByID)
			r.Get("/properties", propertyHandler.GetAllProperties)
			r.Get("/properties/{propertyId}", propertyHandler.GetPropertyByID)
			r.Get("/users", userHandler.GetAllUsers)
			r.Get("/users/{userId}", userHandler.GetUserByID)
			r.Get("/leads", leadHandler.GetAllLeads)
			r.Get("/deals", dealHandler.GetAllDeals)
			
			// Full CRUD for personal/organizational items for everyone.
			r.Mount("/tasks", TaskRoutes(taskHandler))
			r.Mount("/notes", NoteRoutes(noteHandler))
			r.Mount("/events", EventRoutes(eventHandler))
			r.Mount("/comm-logs", CommLogRoutes(commLogHandler))
			
			// Nested GET routes for everyone.
			r.Get("/contacts/{contactId}/comm-logs", commLogHandler.GetLogsForContact)
			r.Get("/users/{userId}/notes", noteHandler.GetNotesByUser)
			r.Get("/users/{userId}/events", eventHandler.GetEventsForUser)

			// --- Reception ONLY Routes (Managerial Role) ---
			r.Group(func(r chi.Router) {
				r.Use(AuthorizeRole(dto.RoleReception))
				
				// Full CRUD for Contacts, Properties, and Leads.
				r.Post("/contacts", contactHandler.CreateContact)
				r.Put("/contacts/{contactId}", contactHandler.UpdateContact)
				r.Delete("/contacts/{contactId}", contactHandler.DeleteContact)
				
				r.Post("/properties", propertyHandler.CreateProperty)
				r.Put("/properties/{propertyId}", propertyHandler.UpdateProperty)
				r.Delete("/properties/{propertyId}", propertyHandler.DeleteProperty)

				r.Post("/leads", leadHandler.CreateLead)
				r.Get("/leads/{id}", leadHandler.GetLeadByID)
				r.Put("/leads/{id}", leadHandler.UpdateLead)
				r.Delete("/leads/{id}", leadHandler.DeleteLead)

				// Full access to Reports.
				r.Get("/reports/employee-leads", reportHandler.GetEmployeeLeadReport)
				r.Get("/reports/source-leads", reportHandler.GetSourceLeadReport)
				r.Get("/reports/employee-sales", reportHandler.GetEmployeeSalesReport)
				r.Get("/reports/source-sales", reportHandler.GetSourceSalesReport)
			})

			// --- Sales Agent ONLY Routes ---
			r.Group(func(r chi.Router) {
				r.Use(AuthorizeRole(dto.RoleSalesAgent))
				
				// Sales Agents manage the full Deal lifecycle.
				r.Post("/deals", dealHandler.CreateDeal)
				r.Get("/deals/{id}", dealHandler.GetDealByID)
				r.Put("/deals/{id}", dealHandler.UpdateDeal)
				r.Delete("/deals/{id}", dealHandler.DeleteDeal)
			})
		})
	})
	return r
}
// --- Helper Functions (No changes needed below) ---

func TaskRoutes(h *handlers.TaskHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.GetAllTasks)
	r.Post("/", h.CreateTask)
	r.Get("/{id}", h.GetTaskByID)
	r.Put("/{id}", h.UpdateTask)
	r.Delete("/{id}", h.DeleteTask)
	return r
}
func NoteRoutes(h *handlers.NoteHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.GetAllNotes)
	r.Post("/", h.CreateNote)
	r.Get("/{noteId}", h.GetNoteByID)
	r.Put("/{noteId}", h.UpdateNote)
	r.Delete("/{noteId}", h.DeleteNote)
	return r
}
func EventRoutes(h *handlers.EventHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.GetAllEvents)
	r.Post("/", h.CreateEvent)
	r.Get("/{eventId}", h.GetEventByID)
	r.Put("/{eventId}", h.UpdateEvent)
	r.Delete("/{eventId}", h.DeleteEvent)
	return r
}
func CommLogRoutes(h *handlers.CommLogHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.CreateLog)
	r.Get("/{logId}", h.GetLogByID)
	r.Put("/{logId}", h.UpdateLog)
	r.Delete("/{logId}", h.DeleteLog)
	return r
}