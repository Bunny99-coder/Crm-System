
package api

import (
	"crm-project/internal/api/handlers"
	"crm-project/internal/config"
	"crm-project/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	cfg *config.Config,
	jwtSecret string,
	authHandler *handlers.AuthHandler,
	contactHandler *handlers.ContactHandler,
	userHandler *handlers.UserHandler,
	propertyHandler *handlers.PropertyHandler,
	leadHandler *handlers.LeadHandler,
	dealHandler *handlers.DealHandler,
	reportHandler *handlers.ReportHandler,
	taskHandler *handlers.TaskHandler,
	commLogHandler *handlers.CommLogHandler,
	noteHandler *handlers.NoteHandler,
	eventHandler *handlers.EventHandler,
	tempHandler *handlers.TempHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Temporary route to grant reception role
	r.Post("/temp/grant-reception-role", tempHandler.GrantReceptionRole)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/register", authHandler.Register)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(jwtSecret))

			// User Routes
			r.Get("/users", userHandler.GetAllUsers)
			r.Post("/users", userHandler.CreateUser)
			r.Get("/users/{userId}", userHandler.GetUserByID)

			// Contact Routes
			r.Get("/contacts", contactHandler.GetAllContacts)
			r.Post("/contacts", contactHandler.CreateContact)
			r.Get("/contacts/{contactId}", contactHandler.GetContactByID)
			r.Put("/contacts/{contactId}", contactHandler.UpdateContact)
			r.Delete("/contacts/{contactId}", contactHandler.DeleteContact)

			// Property Routes
			r.Get("/properties", propertyHandler.GetAllProperties)
			r.Post("/properties", propertyHandler.CreateProperty)
			r.Get("/properties/{propertyId}", propertyHandler.GetPropertyByID)
			r.Put("/properties/{propertyId}", propertyHandler.UpdateProperty)
			r.Delete("/properties/{propertyId}", propertyHandler.DeleteProperty)

			// Lead Routes
			r.Get("/leads", leadHandler.GetAllLeads)
			r.Post("/leads", leadHandler.CreateLead)
			r.Get("/leads/{id}", leadHandler.GetLeadByID)
			r.Put("/leads/{id}", leadHandler.UpdateLead)
			r.Delete("/leads/{id}", leadHandler.DeleteLead)

			// Deal Routes
			r.Get("/deals", dealHandler.GetAllDeals)
			r.Post("/deals", dealHandler.CreateDeal)
			r.Get("/deals/{id}", dealHandler.GetDealByID)
			r.Put("/deals/{id}", dealHandler.UpdateDeal)
			r.Delete("/deals/{id}", dealHandler.DeleteDeal)

			// Task Routes
			// Reception and Sales Agents can view and update tasks
			r.Group(func(r chi.Router) {
				r.Use(AuthorizeRole(util.RoleReception, util.RoleSalesAgent))
				r.Get("/tasks", taskHandler.GetAllTasks)
				r.Get("/tasks/{id}", taskHandler.GetTaskByID)
				r.Put("/tasks/{id}", taskHandler.UpdateTask)
			})

			// Only Reception can create and delete tasks
			r.Group(func(r chi.Router) {
				r.Use(AuthorizeRole(util.RoleReception))
				r.Post("/tasks", taskHandler.CreateTask)
				r.Delete("/tasks/{id}", taskHandler.DeleteTask)
			})

			// Note Routes
			r.Get("/contacts/{contactId}/notes", noteHandler.GetContactNotes)
			r.Post("/contacts/{contactId}/notes", noteHandler.CreateNote)
			r.Get("/contacts/{contactId}/notes/{noteId}", noteHandler.GetNoteByID)
			r.Put("/contacts/{contactId}/notes/{noteId}", noteHandler.UpdateNote)
			r.Delete("/contacts/{contactId}/notes/{noteId}", noteHandler.DeleteNote)

			// Event Routes
			r.Get("/events", eventHandler.GetAllEvents)
			r.Post("/events", eventHandler.CreateEvent)
			r.Get("/events/{eventId}", eventHandler.GetEventByID)
			r.Put("/events/{eventId}", eventHandler.UpdateEvent)
			r.Delete("/events/{eventId}", eventHandler.DeleteEvent)

			// Communication Log Routes
			r.Post("/comm-logs", commLogHandler.CreateCommLog)
			r.Get("/comm-logs/{logId}", commLogHandler.GetCommLogByID)
			r.Put("/comm-logs/{logId}", commLogHandler.UpdateCommLog)
			r.Delete("/comm-logs/{logId}", commLogHandler.DeleteCommLog)

			// Contact-specific Communication Log Routes
			r.Get("/contacts/{contactId}/comm-logs", commLogHandler.GetLogsForContact)
			r.Post("/contacts/{contactId}/comm-logs", commLogHandler.CreateContactCommLog)
			r.Put("/contacts/{contactId}/comm-logs/{logId}", commLogHandler.UpdateContactCommLog)
			r.Delete("/contacts/{contactId}/comm-logs/{logId}", commLogHandler.DeleteContactCommLog)

			// User-specific Notes and Events
			r.Get("/users/{userId}/notes", noteHandler.GetUserNotes)
			r.Get("/users/{userId}/events", eventHandler.GetEventsForUser)

			// Report Routes
			r.Get("/reports/employee-leads", reportHandler.GetEmployeeLeadReport)
			r.Get("/reports/employee-sales", reportHandler.GetEmployeeSalesReport)
			r.Get("/reports/source-leads", reportHandler.GetSourceLeadReport)
			r.Get("/reports/source-sales", reportHandler.GetSourceSalesReport)
			r.Get("/reports/my-sales", reportHandler.GetMySalesReport)
			r.Get("/reports/deals-pipeline", reportHandler.GetDealsPipelineReport)
		})
	})

	return r
}
