package api

import (
    "crm-project/internal/api/handlers"
    "crm-project/internal/dto"
    "log/slog"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    chiMiddleware "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
)

// Context key for user ID
type contextKey string

const UserIDContextKey contextKey = "userID"

func NewRouter(
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
) http.Handler {

    r := chi.NewRouter()

    // --- Custom Handlers for Debugging ---
    r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
        slog.Warn("Method not allowed", "method", r.Method, "url", r.URL.Path, "headers", r.Header)
        w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    })
    r.NotFound(func(w http.ResponseWriter, r *http.Request) {
        slog.Warn("Route not found", "method", r.Method, "url", r.URL.Path, "headers", r.Header)
        http.Error(w, "Not Found", http.StatusNotFound)
    })

    // --- Global Middleware ---
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        AllowCredentials: true,
    }))
    r.Use(chiMiddleware.Logger)
    r.Use(chiMiddleware.Recoverer)
    r.Use(TimeoutMiddleware(30 * time.Second))
    slog.Info("Starting router setup")

    r.Route("/api/v1", func(r chi.Router) {
        slog.Info("Setting up /api/v1 routes")

        // --- Public Routes ---
        r.Post("/auth/login", authHandler.Login)
        r.Post("/users", userHandler.CreateUser)
        slog.Info("Public routes registered: /auth/login, /users")

        // --- Protected Routes ---
        r.Group(func(r chi.Router) {
            r.Use(AuthMiddleware(jwtSecret))
            slog.Info("Protected group started with AuthMiddleware")

			r.Post("/auth/logout", authHandler.Logout)
			slog.Info("Registered POST /auth/logout")

            // --- Contacts ---
            r.Get("/contacts", contactHandler.GetAllContacts)
            r.Post("/contacts", contactHandler.CreateContact)
            r.Get("/contacts/{contactId}", contactHandler.GetContactByID)
            r.Put("/contacts/{contactId}", contactHandler.UpdateContact)
            r.Delete("/contacts/{contactId}", contactHandler.DeleteContact)
            slog.Info("Contact routes registered")

            // --- Contact Nested Resources ---
            r.Route("/contacts/{contactId}", func(r chi.Router) {
                slog.Info("Setting up /contacts/{contactId} sub-routes")

                // Notes
                r.Route("/notes", func(r chi.Router) {
                    r.Get("/", noteHandler.GetContactNotes) // Changed from "" to "/"
                    slog.Info("Registered GET /contacts/{contactId}/notes")
                    r.Post("/", noteHandler.CreateNote) // Changed from "" to "/"
                    slog.Info("Registered POST /contacts/{contactId}/notes")
                    r.Route("/{noteId}", func(r chi.Router) {
                        r.Get("/", noteHandler.GetNoteByID) // Changed from "" to "/"
                        slog.Info("Registered GET /contacts/{contactId}/notes/{noteId}")
                        r.Put("/", noteHandler.UpdateNote) // Changed from "" to "/"
                        slog.Info("Registered PUT /contacts/{contactId}/notes/{noteId}")
                        r.Delete("/", noteHandler.DeleteNote) // Changed from "" to "/"
                        slog.Info("Registered DELETE /contacts/{contactId}/notes/{noteId}")
                    })
                })
                slog.Info("Contact notes sub-routes registered")

                // Communication Logs
                r.Route("/comm-logs", func(r chi.Router) {
                    r.Get("/", commLogHandler.GetLogsForContact) // Changed from "" to "/"
                    slog.Info("Registered GET /contacts/{contactId}/comm-logs")
                    r.Post("/", commLogHandler.CreateContactCommLog) // Changed from "" to "/"
                    slog.Info("Registered POST /contacts/{contactId}/comm-logs")
                    r.Route("/{logId}", func(r chi.Router) {
                        r.Get("/", commLogHandler.GetCommLogByID) // Changed from "" to "/"
                        slog.Info("Registered GET /contacts/{contactId}/comm-logs/{logId}")
                        r.Put("/", commLogHandler.UpdateContactCommLog) // Changed from "" to "/"
                        slog.Info("Registered PUT /contacts/{contactId}/comm-logs/{logId}")
                        r.Delete("/", commLogHandler.DeleteContactCommLog) // Changed from "" to "/"
                        slog.Info("Registered DELETE /contacts/{contactId}/comm-logs/{logId}")
                    })
                })
                slog.Info("Contact comm-logs sub-routes registered")
            })

            // --- Properties ---
            r.Get("/properties", propertyHandler.GetAllProperties)
            r.Post("/properties", propertyHandler.CreateProperty)
            r.Get("/properties/{propertyId}", propertyHandler.GetPropertyByID)
            r.Put("/properties/{propertyId}", propertyHandler.UpdateProperty)
            r.Delete("/properties/{propertyId}", propertyHandler.DeleteProperty)
            slog.Info("Property routes registered")

            // --- Users ---
            r.Get("/users", userHandler.GetAllUsers)
            r.Get("/users/{userId}", userHandler.GetUserByID)
            r.Get("/users/{userId}/events", eventHandler.GetEventsForUser)
            slog.Info("User routes registered")

            // --- Leads ---
            r.Get("/leads", leadHandler.GetAllLeads)
            r.Post("/leads", leadHandler.CreateLead)
            r.Get("/leads/{leadId}", leadHandler.GetLeadByID)
            r.Put("/leads/{leadId}", leadHandler.UpdateLead)
            r.Delete("/leads/{leadId}", leadHandler.DeleteLead)
            slog.Info("Lead routes registered")

            // --- Deals ---
            r.Get("/deals", dealHandler.GetAllDeals)
            r.Post("/deals", dealHandler.CreateDeal)
            r.Get("/deals/{dealId}", dealHandler.GetDealByID)
            r.Put("/deals/{dealId}", dealHandler.UpdateDeal)
            r.Delete("/deals/{dealId}", dealHandler.DeleteDeal)
            slog.Info("Deal base routes registered")

            // --- Deal Nested Resources ---
            r.Route("/deals/{dealId}", func(r chi.Router) {
                slog.Info("Setting up /deals/{dealId} sub-routes")

                // Notes
                r.Route("/notes", func(r chi.Router) {
                    r.Get("/", noteHandler.GetDealNotes) // Changed from "" to "/"
                    slog.Info("Registered GET /deals/{dealId}/notes")
                    r.Post("/", noteHandler.CreateDealNote) // Changed from "" to "/"
                    slog.Info("Registered POST /deals/{dealId}/notes")
                    r.Route("/{noteId}", func(r chi.Router) {
                        r.Get("/", noteHandler.GetDealNoteByID) // Changed from "" to "/"
                        slog.Info("Registered GET /deals/{dealId}/notes/{noteId}")
                        r.Put("/", noteHandler.UpdateDealNote) // Changed from "" to "/"
                        slog.Info("Registered PUT /deals/{dealId}/notes/{noteId}")
                        r.Delete("/", noteHandler.DeleteDealNote) // Changed from "" to "/"
                        slog.Info("Registered DELETE /deals/{dealId}/notes/{noteId}")
                    })
                })
                slog.Info("Deal notes sub-routes registered")

                // Tasks
                r.Route("/tasks", func(r chi.Router) {
                    r.Get("/", taskHandler.GetDealTasks) // Changed from "" to "/"
                    slog.Info("Registered GET /deals/{dealId}/tasks")
                    r.Post("/", taskHandler.CreateDealTask) // Changed from "" to "/"
                    slog.Info("Registered POST /deals/{dealId}/tasks")
                    r.Route("/{taskId}", func(r chi.Router) {
                        r.Get("/", taskHandler.GetTaskByID) // Changed from "" to "/"
                        slog.Info("Registered GET /deals/{dealId}/tasks/{taskId}")
                        r.Put("/", taskHandler.UpdateDealTask) // Changed from "" to "/"
                        slog.Info("Registered PUT /deals/{dealId}/tasks/{taskId}")
                        r.Delete("/", taskHandler.DeleteDealTask) // Changed from "" to "/"
                        slog.Info("Registered DELETE /deals/{dealId}/tasks/{taskId}")
                    })
                })
                slog.Info("Deal tasks sub-routes registered")

                // Events
                r.Route("/events", func(r chi.Router) {
                    r.Get("/", eventHandler.GetDealEvents) // Changed from "" to "/"
                    slog.Info("Registered GET /deals/{dealId}/events")
                    r.Post("/", eventHandler.CreateDealEvent) // Changed from "" to "/"
                    slog.Info("Registered POST /deals/{dealId}/events")
                    r.Route("/{eventId}", func(r chi.Router) {
                        r.Get("/", eventHandler.GetEventByID) // Changed from "" to "/"
                        slog.Info("Registered GET /deals/{dealId}/events/{eventId}")
                        r.Put("/", eventHandler.UpdateDealEvent) // Changed from "" to "/"
                        slog.Info("Registered PUT /deals/{dealId}/events/{eventId}")
                        r.Delete("/", eventHandler.DeleteDealEvent) // Changed from "" to "/"
                        slog.Info("Registered DELETE /deals/{dealId}/events/{eventId}")
                    })
                })
                slog.Info("Deal events sub-routes registered")

                // Communication Logs
                r.Route("/comm-logs", func(r chi.Router) {
                    r.Get("/", commLogHandler.GetDealCommLogs) // Changed from "" to "/"
                    slog.Info("Registered GET /deals/{dealId}/comm-logs")
                    r.Post("/", commLogHandler.CreateDealCommLog) // Changed from "" to "/"
                    slog.Info("Registered POST /deals/{dealId}/comm-logs")
                    r.Route("/{logId}", func(r chi.Router) {
                        r.Get("/", commLogHandler.GetCommLogByID) // Changed from "" to "/"
                        slog.Info("Registered GET /deals/{dealId}/comm-logs/{logId}")
                        r.Put("/", commLogHandler.UpdateDealCommLog) // Changed from "" to "/"
                        slog.Info("Registered PUT /deals/{dealId}/comm-logs/{logId}")
                        r.Delete("/", commLogHandler.DeleteDealCommLog) // Changed from "" to "/"
                        slog.Info("Registered DELETE /deals/{dealId}/comm-logs/{logId}")
                    })
                })
                slog.Info("Deal comm-logs sub-routes registered")
            })
            slog.Info("All deal nested sub-routes registered")

            // --- Reports (role-protected) ---
            r.Group(func(r chi.Router) {
                r.Use(AuthorizeRole(dto.RoleReception))
                r.Get("/reports/employee-leads", reportHandler.GetEmployeeLeadReport)
                r.Get("/reports/source-leads", reportHandler.GetSourceLeadReport)
                r.Get("/reports/employee-sales", reportHandler.GetEmployeeSalesReport)
                r.Get("/reports/source-sales", reportHandler.GetSourceSalesReport)
                slog.Info("Report routes registered with role authorization")
            })

            // --- Global Routes (non-nested) ---
            r.Mount("/tasks", TaskRoutes(taskHandler))
            slog.Info("Mounted global /tasks routes")
            r.Mount("/events", EventRoutes(eventHandler))
            slog.Info("Mounted global /events routes")
            r.Mount("/comm-logs", CommLogRoutes(commLogHandler))
            slog.Info("Mounted global /comm-logs routes")
        })
        slog.Info("Protected group completed")
    })

    slog.Info("Router setup completed successfully")

    // Debug: Walk and log all registered routes
    chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
        slog.Info("Registered route", "method", method, "route", route)
        return nil
    })

    return r
}

// --- Helper Functions ---

func TaskRoutes(h *handlers.TaskHandler) http.Handler {
    r := chi.NewRouter()
    r.Get("/", h.GetAllTasks)
    r.Post("/", h.CreateTask)
    r.Get("/{taskId}", h.GetTaskByID)
    r.Put("/{taskId}", h.UpdateTask)
    r.Delete("/{taskId}", h.DeleteTask)
    slog.Info("TaskRoutes helper registered")
    return r
}

func EventRoutes(h *handlers.EventHandler) http.Handler {
    r := chi.NewRouter()
    r.Get("/", h.GetAllEvents)
    r.Post("/", h.CreateEvent)
    r.Get("/{eventId}", h.GetEventByID)
    r.Put("/{eventId}", h.UpdateEvent)
    r.Delete("/{eventId}", h.DeleteEvent)
    slog.Info("EventRoutes helper registered")
    return r
}

func CommLogRoutes(h *handlers.CommLogHandler) http.Handler {
    r := chi.NewRouter()
    r.Get("/", h.GetAllCommLogs)
    r.Post("/", h.CreateCommLog)
    r.Get("/{logId}", h.GetCommLogByID)
    r.Put("/{logId}", h.UpdateCommLog)
    r.Delete("/{logId}", h.DeleteCommLog)
    slog.Info("CommLogRoutes helper registered")
    return r
}