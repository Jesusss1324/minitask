package repository

import (
	"context"
	"database/sql"
	"strings"

	"minitask/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks (
			title,
			description,
			status,
			priority,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	task.ID = uint64(id)

	return nil
}

func (r *TaskRepository) List(
	ctx context.Context,
	status string,
	search string,
) ([]models.Task, error) {

	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	if search != "" {
		escapedSearch := strings.NewReplacer(
			"\\", "\\\\",
			"%", "\\%",
			"_", "\\_",
		).Replace(search)

		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+escapedSearch+"%")
	}

	query := `
		SELECT
			id,
			title,
			description,
			status,
			priority,
			created_at,
			updated_at
		FROM tasks
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC, id DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
