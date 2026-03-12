package service

import (
	"test-backend/internal/model"
	"test-backend/internal/repository"
	apperr "test-backend/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

// UserService handles user business logic.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService returns a new UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create creates a user (hashes password). Returns conflict if email exists.
func (s *UserService) Create(req *model.CreateUserRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(req.Email, req.Name, string(hash))
}

// Login validates email and password and returns the user or ErrUnauthorized.
func (s *UserService) Login(email, password string) (*model.User, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil || u == nil {
		return nil, apperr.NewUnauthorized("Invalid email or password.")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, apperr.NewUnauthorized("Invalid email or password.")
	}
	return u, nil
}

// GetByID returns the user by ID (UUID) or ErrNotFound.
func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.repo.GetByID(id)
}

// List returns all users.
func (s *UserService) List() ([]*model.User, error) {
	return s.repo.List()
}

// Update updates the user by ID. Returns ErrNotFound or Conflict.
func (s *UserService) Update(id string, req *model.UpdateUserRequest) (*model.User, error) {
	return s.repo.Update(id, req.Name, req.Email)
}

// Delete deletes the user by ID. Returns ErrNotFound.
func (s *UserService) Delete(id string) error {
	return s.repo.Delete(id)
}

// ResetPassword sets a new password for the user. Returns ErrNotFound.
func (s *UserService) ResetPassword(id string, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(id, string(hash))
}
