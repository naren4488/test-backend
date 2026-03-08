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

	// Middleware: request ID first (so logger can use it), then logger, then recoverer
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", healthHandler(db))

	// Users API (full paths so /api/v1/users and /api/v1/users/ both work)
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	r.Post("/api/v1/users", userHandler.Create)
	r.Get("/api/v1/users", userHandler.List)
	r.Post("/api/v1/users/{id}/reset-password", userHandler.ResetPassword)
	r.Get("/api/v1/users/{id}", userHandler.GetByID)
	r.Put("/api/v1/users/{id}", userHandler.Update)
	r.Delete("/api/v1/users/{id}", userHandler.Delete)

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
