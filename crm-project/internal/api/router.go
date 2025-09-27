
package api

import (
	"crm-project/internal/api/handlers"
	"crm-project/internal/config"
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
	tempHandler *handlers.TempHandler, // Add this line
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

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(jwtSecret))

			r.Get("/users", userHandler.GetAllUsers)
			r.Post("/users", userHandler.CreateUser)
			r.Get("/users/{userId}", userHandler.GetUserByID)

			r.Get("/contacts", contactHandler.GetAllContacts)
			r.Post("/contacts", contactHandler.CreateContact)
			r.Get("/contacts/{contactId}", contactHandler.GetContactByID)
			r.Put("/contacts/{contactId}", contactHandler.UpdateContact)
			r.Delete("/contacts/{contactId}", contactHandler.DeleteContact)

			r.Get("/properties", propertyHandler.GetAllProperties)
			r.Post("/properties", propertyHandler.CreateProperty)
			r.Get("/properties/{propertyId}", propertyHandler.GetPropertyByID)
			r.Put("/properties/{propertyId}", propertyHandler.UpdateProperty)
			r.Delete("/properties/{propertyId}", propertyHandler.DeleteProperty)

			r.Get("/leads", leadHandler.GetAllLeads)
			r.Post("/leads", leadHandler.CreateLead)
			r.Get("/leads/{id}", leadHandler.GetLeadByID)
			r.Put("/leads/{id}", leadHandler.UpdateLead)
			r.Delete("/leads/{id}", leadHandler.DeleteLead)

			r.Get("/deals", dealHandler.GetAllDeals)
			r.Post("/deals", dealHandler.CreateDeal)
			r.Get("/deals/{id}", dealHandler.GetDealByID)
			r.Put("/deals/{id}", dealHandler.UpdateDeal)
			r.Delete("/deals/{id}", dealHandler.DeleteDeal)

			r.Get("/tasks", taskHandler.GetAllTasks)
			r.Post("/tasks", taskHandler.CreateTask)
			r.Get("/tasks/{id}", taskHandler.GetTaskByID)
			r.Put("/tasks/{id}", taskHandler.UpdateTask)
			r.Delete("/tasks/{id}", taskHandler.DeleteTask)

			r.Get("/contacts/{contactId}/notes", noteHandler.GetContactNotes)
			r.Post("/contacts/{contactId}/notes", noteHandler.CreateNote)
			r.Get("/contacts/{contactId}/notes/{noteId}", noteHandler.GetNoteByID)
			r.Put("/contacts/{contactId}/notes/{noteId}", noteHandler.UpdateNote)
			r.Delete("/contacts/{contactId}/notes/{noteId}", noteHandler.DeleteNote)

			r.Get("/events", eventHandler.GetAllEvents)
			r.Post("/events", eventHandler.CreateEvent)
			r.Get("/events/{eventId}", eventHandler.GetEventByID)
			r.Put("/events/{eventId}", eventHandler.UpdateEvent)
			r.Delete("/events/{eventId}", eventHandler.DeleteEvent)

			r.Post("/comm-logs", commLogHandler.CreateCommLog)
			r.Get("/comm-logs/{logId}", commLogHandler.GetCommLogByID)
			r.Put("/comm-logs/{logId}", commLogHandler.UpdateCommLog)
			r.Delete("/comm-logs/{logId}", commLogHandler.DeleteCommLog)

			r.Get("/contacts/{contactId}/comm-logs", commLogHandler.GetLogsForContact)
			r.Post("/contacts/{contactId}/comm-logs", commLogHandler.CreateContactCommLog)
			r.Put("/contacts/{contactId}/comm-logs/{logId}", commLogHandler.UpdateContactCommLog)
			r.Delete("/contacts/{contactId}/comm-logs/{logId}", commLogHandler.DeleteContactCommLog)

			r.Get("/users/{userId}/notes", noteHandler.GetUserNotes)
			r.Get("/users/{userId}/events", eventHandler.GetEventsForUser)

			r.Get("/reports/employee-leads", reportHandler.GetEmployeeLeadReport)
			r.Get("/reports/employee-sales", reportHandler.GetEmployeeSalesReport)
			r.Get("//reports/source-leads", reportHandler.GetSourceLeadReport)
			r.Get("/reports/source-sales", reportHandler.GetSourceSalesReport)
			r.Get("/reports/my-sales", reportHandler.GetMySalesReport)
			r.Get("/reports/deals-pipeline", reportHandler.GetDealsPipelineReport)
			})
		})

	return r
}
