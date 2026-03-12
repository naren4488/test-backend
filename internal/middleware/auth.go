package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	apperr "test-backend/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"

// GetUserID returns the authenticated user's ID from the request context.
// Returns ("", false) if not set (e.g. no auth middleware or not logged in).
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok && id != ""
}

// Auth validates the Bearer JWT and sets the user ID in context.
// Returns 401 if missing or invalid. Requires JWTSecret to be non-empty.
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if jwtSecret == "" {
				writeAuthError(w, apperr.NewUnauthorized("Authentication is not configured."))
				return
			}
			auth := r.Header.Get("Authorization")
			if auth == "" {
				writeAuthError(w, apperr.NewUnauthorized("Missing or invalid token."))
				return
			}
			const prefix = "Bearer "
			if !strings.HasPrefix(auth, prefix) {
				writeAuthError(w, apperr.NewUnauthorized("Missing or invalid token."))
				return
			}
			tokenStr := strings.TrimSpace(auth[len(prefix):])
			if tokenStr == "" {
				writeAuthError(w, apperr.NewUnauthorized("Missing or invalid token."))
				return
			}
			claims := &jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(*jwt.Token) (any, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				writeAuthError(w, apperr.NewUnauthorized("Missing or invalid token."))
				return
			}
			userID := claims.Subject
			if userID == "" {
				writeAuthError(w, apperr.NewUnauthorized("Missing or invalid token."))
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeAuthError(w http.ResponseWriter, err *apperr.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": err.Code, "message": err.Message},
	})
}
