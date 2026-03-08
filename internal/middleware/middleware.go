package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// RequestIDKey is the context key for the request ID.
type RequestIDKey struct{}

// RequestID sets a unique X-Request-ID on each request and stores it in context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), RequestIDKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger logs each request with method, path, request ID, status and duration.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		reqID := r.Context().Value(RequestIDKey{})
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"request_id", reqID,
		)
		next.ServeHTTP(ww, r)
		slog.Info("response",
			"status", ww.Status(),
			"request_id", reqID,
		)
	})
}

// Recoverer recovers from panics, logs the panic, and returns HTTP 500.
func Recoverer(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}
