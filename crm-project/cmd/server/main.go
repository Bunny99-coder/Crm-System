package main

import (
	"context"
	"crm-project/internal/api"
	"crm-project/internal/api/handlers"
	"crm-project/internal/config"
	"crm-project/internal/repository/postgres"
	"crm-project/internal/service"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
    "github.com/go-chi/chi/v5" // <-- Make sure chi is imported here
"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	// --- Initialize Logger & Config ---
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("Test log to confirm logging is working")
	cfg, err := config.Load("config.yml", logger)
	if err != nil {
		logger.Error("could not load configuration", "error", err)
		os.Exit(1)
	}
	logger.Info("configuration loaded successfully")

	// --- Connect to Database ---
	db, err := sqlx.Connect("pgx", cfg.Database.URL)
	if err != nil {
		logger.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		logger.Error("could not ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("successfully connected to the database!")

	// --- Dependency Injection ---
	// Repository Layer
	contactRepo := postgres.NewContactRepo(db)
	userRepo := postgres.NewUserRepo(db)
	propertyRepo := postgres.NewPropertyRepo(db)
	leadRepo := postgres.NewLeadRepo(db)
	dealRepo := postgres.NewDealRepo(db)
taskRepo := postgres.NewTaskRepository(db)	
	commLogRepo := postgres.NewCommLogRepository(db) // Corrected from NewCommLogRepo
noteRepo := postgres.NewNoteRepository(db) // <- pass the underlying *sql.DB
eventRepo := postgres.NewEventRepository(db)



	// Service Layer
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, logger)
	contactService := service.NewContactService(contactRepo, logger)
	userService := service.NewUserService(userRepo, logger)
	propertyService := service.NewPropertyService(propertyRepo, logger)
	leadService := service.NewLeadService(leadRepo, contactRepo, userRepo, propertyRepo, logger)
	dealService := service.NewDealService(dealRepo, leadRepo, propertyRepo, logger)
	reportService := service.NewReportService(userRepo, leadRepo, dealRepo, logger)
taskService := service.NewTaskService(taskRepo)	
	commLogService := service.NewCommLogService(commLogRepo) // Corrected to match service constructor
noteService := service.NewNoteService(noteRepo)
eventService := service.NewEventService(eventRepo)
	// Handler Layer



	authHandler := handlers.NewAuthHandler(authService, logger)
	contactHandler := handlers.NewContactHandler(contactService, logger)
	userHandler := handlers.NewUserHandler(userService, logger)
	propertyHandler := handlers.NewPropertyHandler(propertyService, logger)
	leadHandler := handlers.NewLeadHandler(leadService, logger)
	dealHandler := handlers.NewDealHandler(dealService, logger)
	reportHandler := handlers.NewReportHandler(reportService, logger)
	
taskHandler := handlers.NewTaskHandler(taskService)	
	commLogHandler := handlers.NewCommLogHandler(commLogService) // Corrected to match handler constructor
noteHandler := handlers.NewNoteHandler(noteService)
eventHandler := handlers.NewEventHandler(eventService)	// Router
	router := api.NewRouter(
		cfg.Auth.JWTSecret,
		authHandler,
		contactHandler,
		userHandler,
		propertyHandler,
		leadHandler,
		dealHandler,
		reportHandler,
		taskHandler,
		commLogHandler,
		noteHandler,
		eventHandler,
	)
	
	
    // --- PASTE THIS ENTIRE BLOCK OF CODE HERE ---
	logger.Info("--- Registered Routes ---")
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s \t %s\n", method, route)
		return nil
	}
if err := chi.Walk(router.(chi.Router), walkFunc); err != nil {
		logger.Error("failed to walk routes", "error", err)
	}
	logger.Info("-------------------------")
    // --- END OF THE BLOCK TO PASTE ---






	
	// --- Create and Start the HTTP Server ---
	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	go func() {
		logger.Info("server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen and serve error", "error", err)
			os.Exit(1)
		}
	}()

	// --- Graceful Shutdown Logic ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("shutdown signal received, starting graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited gracefully")
}