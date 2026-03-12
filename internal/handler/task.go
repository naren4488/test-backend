package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"test-backend/internal/middleware"
	"test-backend/internal/model"
	"test-backend/internal/response"
	"test-backend/internal/service"
	apperr "test-backend/pkg/errors"

	"github.com/go-playground/validator/v10"
)

// TaskHandler handles task HTTP requests.
type TaskHandler struct {
	svc      *service.TaskService
	validate *validator.Validate
}

// NewTaskHandler returns a new TaskHandler.
func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc, validate: validator.New()}
}

// Create handles POST /api/v1/users/:id/tasks (id = user ID).
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDFromPath, ok := parseUUIDParam(r, "id")
	if !ok {
		response.Error(w, apperr.NewBadRequest("Invalid user ID. Must be a valid UUID."))
		return
	}
	if !requireSelf(w, r, userIDFromPath) {
		return
	}
	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperr.NewBadRequest("Invalid request body. Expected valid JSON."))
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		response.ValidationError(w, "One or more fields failed validation. Check the 'fields' object for details.", validationErrors(err))
		return
	}
	task, err := h.svc.Create(userIDFromPath, &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, "Task created successfully.", task)
}

// ListByUser handles GET /api/v1/users/:id/tasks (id = user ID).
func (h *TaskHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userIDFromPath, ok := parseUUIDParam(r, "id")
	if !ok {
		response.Error(w, apperr.NewBadRequest("Invalid user ID. Must be a valid UUID."))
		return
	}
	if !requireSelf(w, r, userIDFromPath) {
		return
	}
	tasks, err := h.svc.ListByUserID(userIDFromPath)
	if err != nil {
		response.Error(w, err)
		return
	}
	if tasks == nil {
		tasks = []*model.Task{}
	}
	response.Success(w, "Tasks retrieved successfully.", tasks)
}

// ListAll handles GET /api/v1/tasks?page=1&limit=10&search=...&user_id=...
// When authenticated, returns only the current user's tasks (user_id filter ignored for ownership).
func (h *TaskHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, apperr.NewUnauthorized("Authentication required."))
		return
	}
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	search := q.Get("search")
	// Only return current user's tasks
	userIDFilter := &userID
	tasks, total, err := h.svc.ListAll(page, limit, search, userIDFilter)
	if err != nil {
		response.Error(w, err)
		return
	}
	if tasks == nil {
		tasks = []*model.Task{}
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	data := model.ListTasksResponse{
		Tasks: tasks,
		Meta:  model.ListTasksMeta{Total: total, Page: page, Limit: limit},
	}
	response.Success(w, "Tasks retrieved successfully.", data)
}

// GetByID handles GET /api/v1/tasks/:id.
func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(r, "id")
	if !ok {
		response.Error(w, apperr.NewBadRequest("Invalid task ID. Must be a valid UUID."))
		return
	}
	task, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(w, err)
		return
	}
	if !requireTaskOwner(w, r, task.UserID) {
		return
	}
	response.Success(w, "Task retrieved successfully.", task)
}

// Update handles PUT /api/v1/tasks/:id.
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(r, "id")
	if !ok {
		response.Error(w, apperr.NewBadRequest("Invalid task ID. Must be a valid UUID."))
		return
	}
	task, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(w, err)
		return
	}
	if !requireTaskOwner(w, r, task.UserID) {
		return
	}
	var req model.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperr.NewBadRequest("Invalid request body. Expected valid JSON."))
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		response.ValidationError(w, "One or more fields failed validation. Check the 'fields' object for details.", validationErrors(err))
		return
	}
	task, err = h.svc.Update(id, &req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, "Task updated successfully.", task)
}

// Delete handles DELETE /api/v1/tasks/:id.
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(r, "id")
	if !ok {
		response.Error(w, apperr.NewBadRequest("Invalid task ID. Must be a valid UUID."))
		return
	}
	task, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(w, err)
		return
	}
	if !requireTaskOwner(w, r, task.UserID) {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.Error(w, err)
		return
	}
	response.OkWithMessage(w, "Task deleted successfully.")
}

// requireTaskOwner checks task's user_id matches authenticated user. Writes 403 and returns false if not.
func requireTaskOwner(w http.ResponseWriter, r *http.Request, taskUserID string) bool {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, apperr.NewUnauthorized("Authentication required."))
		return false
	}
	if taskUserID != userID {
		response.Error(w, apperr.NewForbidden("You can only access your own tasks."))
		return false
	}
	return true
}
