package repository

import (
	"database/sql"
	"strings"

	"test-backend/internal/model"
	apperr "test-backend/pkg/errors"

	"github.com/google/uuid"
)

// UserRepository performs user DB operations.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository returns a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a user with a new UUID and returns it.
func (r *UserRepository) Create(email, name, passwordHash string) (*model.User, error) {
	id := uuid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO users (id, email, name, password_hash) VALUES (?, ?, ?, ?)`,
		id, email, name, passwordHash,
	)
	if err != nil {
		if isSQLiteUnique(err) {
			return nil, apperr.NewConflict("User with this email already exists")
		}
		return nil, err
	}
	return r.GetByID(id)
}

// GetByID returns the user by ID (UUID) or ErrNotFound.
func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(
		`SELECT id, email, name, password_hash, created_at, updated_at FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, apperr.NewNotFound("User not found")
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmail returns the user by email or nil if not found.
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(
		`SELECT id, email, name, password_hash, created_at, updated_at FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// List returns all users ordered by id.
func (r *UserRepository) List() ([]*model.User, error) {
	rows, err := r.db.Query(
		`SELECT id, email, name, password_hash, created_at, updated_at FROM users ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &u)
	}
	return list, rows.Err()
}

// Update updates name and/or email for the user by ID. Returns ErrNotFound if not found.
func (r *UserRepository) Update(id string, name, email *string) (*model.User, error) {
	var set string
	var args []any
	if name != nil {
		set = "name = ?"
		args = append(args, *name)
	}
	if email != nil {
		if set != "" {
			set += ", "
		}
		set += "email = ?"
		args = append(args, *email)
	}
	if set == "" {
		return r.GetByID(id)
	}
	set += ", updated_at = CURRENT_TIMESTAMP"
	args = append(args, id)
	res, err := r.db.Exec(`UPDATE users SET `+set+` WHERE id = ?`, args...)
	if err != nil {
		if isSQLiteUnique(err) {
			return nil, apperr.NewConflict("User with this email already exists")
		}
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, apperr.NewNotFound("User not found")
	}
	return r.GetByID(id)
}

// Delete removes the user by ID. Returns ErrNotFound if not found.
func (r *UserRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.NewNotFound("User not found")
	}
	return nil
}

// UpdatePassword sets the password hash for the user by ID. Returns ErrNotFound if not found.
func (r *UserRepository) UpdatePassword(id string, passwordHash string) error {
	res, err := r.db.Exec(`UPDATE users SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, passwordHash, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.NewNotFound("User not found")
	}
	return nil
}

// isSQLiteUnique returns true if the error is a SQLite unique constraint violation.
func isSQLiteUnique(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE")
}
