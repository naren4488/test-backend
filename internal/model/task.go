package model

import "time"

// Task is the domain model for a task (matches DB row).
// ID and UserID are UUIDs.
type Task struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest is the body for POST /api/v1/users/:id/tasks.
type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,max=500"`
	Description string `json:"description" validate:"max=2000"`
	Status      string `json:"status" validate:"omitempty,oneof=pending in_progress done"`
}

// UpdateTaskRequest is the body for PUT /api/v1/tasks/:id (all optional).
type UpdateTaskRequest struct {
	Title       *string `json:"title" validate:"omitempty,max=500"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	Status      *string `json:"status" validate:"omitempty,oneof=pending in_progress done"`
}

// ListTasksMeta holds pagination info for list-all response.
type ListTasksMeta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// ListTasksResponse is the data shape for GET /api/v1/tasks (paginated list).
type ListTasksResponse struct {
	Tasks []*Task       `json:"tasks"`
	Meta  ListTasksMeta `json:"meta"`
}
