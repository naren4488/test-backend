package service

import (
	"test-backend/internal/model"
	"test-backend/internal/repository"
)

// TaskService handles task business logic.
type TaskService struct {
	taskRepo *repository.TaskRepository
	userRepo *repository.UserRepository
}

// NewTaskService returns a new TaskService.
func NewTaskService(taskRepo *repository.TaskRepository, userRepo *repository.UserRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo, userRepo: userRepo}
}

// Create creates a task for the given user. Returns error if user not found.
func (s *TaskService) Create(userID string, req *model.CreateTaskRequest) (*model.Task, error) {
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, err
	}
	status := req.Status
	if status == "" {
		status = "pending"
	}
	return s.taskRepo.Create(userID, req.Title, req.Description, status)
}

// GetByID returns the task by ID (UUID) or ErrNotFound.
func (s *TaskService) GetByID(id string) (*model.Task, error) {
	return s.taskRepo.GetByID(id)
}

// ListByUserID returns all tasks for the user or ErrNotFound if user does not exist.
func (s *TaskService) ListByUserID(userID string) ([]*model.Task, error) {
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, err
	}
	return s.taskRepo.ListByUserID(userID)
}

// ListAll returns paginated tasks with optional search and user filter.
func (s *TaskService) ListAll(page, limit int, search string, userID *string) ([]*model.Task, int, error) {
	return s.taskRepo.ListAll(page, limit, search, userID)
}

// Update updates the task by ID. Returns ErrNotFound.
func (s *TaskService) Update(id string, req *model.UpdateTaskRequest) (*model.Task, error) {
	return s.taskRepo.Update(id, req.Title, req.Description, req.Status)
}

// Delete deletes the task by ID. Returns ErrNotFound.
func (s *TaskService) Delete(id string) error {
	return s.taskRepo.Delete(id)
}
