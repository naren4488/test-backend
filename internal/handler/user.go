package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"test-backend/internal/model"
	"test-backend/internal/response"
	"test-backend/internal/service"
	apperr "test-backend/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// UserHandler handles user HTTP requests.
type UserHandler struct {
	svc      *service.UserService
	validate *validator.Validate
}

// NewUserHandler returns a new UserHandler.
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc, validate: validator.New()}
}

// Create handles POST /api/v1/users.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperr.NewBadRequest("Invalid request body. Expected valid JSON."))
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		response.ValidationError(w, "One or more fields failed validation. Check the 'fields' object for details.", validationErrors(err))
		return
	}
	user, err := h.svc.Create(&req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "User created successfully.", user)
}

// List handles GET /api/v1/users.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List()
	if err != nil {
		response.Error(w, err)
		return
	}
	if users == nil {
		users = []*model.User{}
	}
	response.Success(w, "Users retrieved successfully.", users)
}

// GetByID handles GET /api/v1/users/:id.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(w, apperr.NewBadRequest("Invalid user ID. Must be a positive integer."))
		return
	}
	user, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "User retrieved successfully.", user)
}

// Update handles PUT /api/v1/users/:id.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(w, apperr.NewBadRequest("Invalid user ID. Must be a positive integer."))
		return
	}
	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperr.NewBadRequest("Invalid request body. Expected valid JSON."))
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		response.ValidationError(w, "One or more fields failed validation. Check the 'fields' object for details.", validationErrors(err))
		return
	}
	user, err := h.svc.Update(id, &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "User updated successfully.", user)
}

// Delete handles DELETE /api/v1/users/:id.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(w, apperr.NewBadRequest("Invalid user ID. Must be a positive integer."))
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.Error(w, err)
		return
	}
	response.OkWithMessage(w, "User deleted successfully.")
}

// ResetPassword handles POST /api/v1/users/:id/reset-password.
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(w, apperr.NewBadRequest("Invalid user ID. Must be a positive integer."))
		return
	}
	var req model.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperr.NewBadRequest("Invalid request body. Expected valid JSON."))
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		response.ValidationError(w, "One or more fields failed validation. Check the 'fields' object for details.", validationErrors(err))
		return
	}
	if err := h.svc.ResetPassword(id, req.NewPassword); err != nil {
		response.Error(w, err)
		return
	}
	response.OkWithMessage(w, "Password reset successfully.")
}

// validationErrors converts validator.ValidationErrors into a map of field -> tag.
func validationErrors(err error) map[string]string {
	fields := make(map[string]string)
	if err == nil {
		return fields
	}
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return fields
	}
	for _, e := range errs {
		fields[e.Field()] = e.Tag()
	}
	return fields
}
