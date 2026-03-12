package repository

import (
	"database/sql"
	"strings"

	"test-backend/internal/model"
	apperr "test-backend/pkg/errors"

	"github.com/google/uuid"
)

// TaskRepository performs task DB operations.
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository returns a new TaskRepository.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create inserts a task with a new UUID and returns it.
func (r *TaskRepository) Create(userID string, title, description, status string) (*model.Task, error) {
	if status == "" {
		status = "pending"
	}
	id := uuid.New().String()
	_, err := r.db.Exec(
		`INSERT INTO tasks (id, user_id, title, description, status) VALUES (?, ?, ?, ?, ?)`,
		id, userID, title, description, status,
	)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// GetByID returns the task by ID (UUID) or ErrNotFound.
func (r *TaskRepository) GetByID(id string) (*model.Task, error) {
	var t model.Task
	err := r.db.QueryRow(
		`SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, apperr.NewNotFound("Task not found")
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ListByUserID returns all tasks for the given user, ordered by id.
func (r *TaskRepository) ListByUserID(userID string) ([]*model.Task, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks WHERE user_id = ? ORDER BY id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListAll returns tasks with pagination and optional search/user filter.
// page and limit are 1-based and 0 means default (page=1, limit=10).
func (r *TaskRepository) ListAll(page, limit int, search string, userID *string) ([]*model.Task, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var where []string
	var args []any
	if search != "" {
		searchPattern := "%" + strings.TrimSpace(search) + "%"
		where = append(where, `(title LIKE ? OR description LIKE ?)`)
		args = append(args, searchPattern, searchPattern)
	}
	if userID != nil {
		where = append(where, `user_id = ?`)
		args = append(args, *userID)
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE " + strings.Join(where, " AND ")
	}

	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM tasks` + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Select page
	selectArgs := append(args, limit, offset)
	rows, err := r.db.Query(
		`SELECT id, user_id, title, description, status, created_at, updated_at FROM tasks`+whereClause+` ORDER BY id LIMIT ? OFFSET ?`,
		selectArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	list, err := scanTasks(rows)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Update updates task by ID. Returns ErrNotFound if not found.
func (r *TaskRepository) Update(id string, title, description, status *string) (*model.Task, error) {
	u, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	set := []string{"updated_at = CURRENT_TIMESTAMP"}
	var args []any
	if title != nil {
		set = append(set, "title = ?")
		args = append(args, *title)
	}
	if description != nil {
		set = append(set, "description = ?")
		args = append(args, *description)
	}
	if status != nil {
		set = append(set, "status = ?")
		args = append(args, *status)
	}
	if len(args) == 0 {
		return u, nil
	}
	args = append(args, id)
	_, err = r.db.Exec(`UPDATE tasks SET `+strings.Join(set, ", ")+` WHERE id = ?`, args...)
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// Delete removes the task by ID. Returns ErrNotFound if not found.
func (r *TaskRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperr.NewNotFound("Task not found")
	}
	return nil
}

func scanTasks(rows *sql.Rows) ([]*model.Task, error) {
	var list []*model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, rows.Err()
}
