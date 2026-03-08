package response

import (
	"encoding/json"
	"net/http"

	apperr "test-backend/pkg/errors"
)

// ErrorBody is the JSON shape for error responses.
type ErrorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Success writes a 200 JSON response with message and data.
func Success(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": message, "data": data})
}

// Created writes a 201 JSON response with message and data.
func Created(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": message, "data": data})
}

// OkWithMessage writes a 200 JSON response with only a message (e.g. for delete, reset-password).
func OkWithMessage(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": message})
}

// NoContent writes 204 No Content (kept for compatibility; prefer OkWithMessage for a message).
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error writes status and JSON error body from an error.
// If err is an AppError, uses its status and message; otherwise 500.
func Error(w http.ResponseWriter, err error) {
	appErr := apperr.AsAppError(err)
	if appErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.StatusCode)
		_ = json.NewEncoder(w).Encode(ErrorBody{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: appErr.Code, Message: appErr.Message},
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorBody{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INTERNAL_ERROR", Message: "An unexpected error occurred. Please try again later."},
		})
}

// ValidationError writes 400 with a list of field errors (e.g. from validator).
func ValidationError(w http.ResponseWriter, message string, fieldErrors map[string]string) {
	if message == "" {
		message = "One or more fields failed validation. Check the 'fields' object for details."
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	body := map[string]any{
		"error": map[string]any{
			"code":    "VALIDATION_ERROR",
			"message": message,
			"fields":  fieldErrors,
		},
	}
	_ = json.NewEncoder(w).Encode(body)
}
