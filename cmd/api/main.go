package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"test-backend/internal/config"
	"test-backend/internal/database"
	"test-backend/internal/handler"
	"test-backend/internal/middleware"
	"test-backend/internal/repository"
	"test-backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		slog.Error("database init", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	r := chi.NewRouter()

	// CORS: allow frontend origins (must be before other middleware so OPTIONS gets CORS headers)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Middleware: request ID first (so logger can use it), then logger, then recoverer
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", healthHandler(db))

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	taskRepo := repository.NewTaskRepository(db)
	taskSvc := service.NewTaskService(taskRepo, userRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)
	authHandler := handler.NewAuthHandler(userSvc, cfg.JWTSecret, cfg.JWTExpiryHours)

	// Public routes (no JWT)
	r.Post("/api/v1/login", authHandler.Login)
	r.Post("/api/v1/users", userHandler.Create)

	// Protected routes (JWT required)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/users", userHandler.List)
		r.Post("/users/{id}/tasks", taskHandler.Create)
		r.Get("/users/{id}/tasks", taskHandler.ListByUser)
		r.Post("/users/{id}/reset-password", userHandler.ResetPassword)
		r.Get("/users/{id}", userHandler.GetByID)
		r.Put("/users/{id}", userHandler.Update)
		r.Delete("/users/{id}", userHandler.Delete)
		r.Get("/tasks", taskHandler.ListAll)
		r.Get("/tasks/{id}", taskHandler.GetByID)
		r.Put("/tasks/{id}", taskHandler.Update)
		r.Delete("/tasks/{id}", taskHandler.Delete)
	})

	addr := ":" + strconv.Itoa(cfg.Port)
	server := &http.Server{Addr: addr, Handler: r}

	go func() {
		slog.Info("server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}

// healthHandler returns 200 with status and optional DB ping.
func healthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := map[string]string{
			"message": "Service is healthy.",
			"status":  "ok",
		}
		if err := db.Ping(); err != nil {
			body["database"] = "disconnected"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(body)
			return
		}
		body["database"] = "connected"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(body)
	}
}
