package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// parseUUIDParam returns the URL path parameter as a UUID string.
// Returns ("", false) if the parameter is missing or not a valid UUID.
func parseUUIDParam(r *http.Request, param string) (string, bool) {
	s := chi.URLParam(r, param)
	if s == "" {
		return "", false
	}
	_, err := uuid.Parse(s)
	if err != nil {
		return "", false
	}
	return s, true
}
