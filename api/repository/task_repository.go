package repository

import (
	"context"
	"database/sql"
	"errors"
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

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*models.Task, error) {

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
		WHERE id = ?
	`

	var task models.Task

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrTaskNotFound
		}

		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, priority = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.UpdatedAt,
		task.ID,
	)

	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) Summary(ctx context.Context) (models.TaskSummary, error) {
	query := `
		SELECT status, COUNT(*)
		FROM tasks
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return models.TaskSummary{}, err
	}
	defer rows.Close()

	var summary models.TaskSummary

	for rows.Next() {
		var status string
		var count int

		if err := rows.Scan(&status, &count); err != nil {
			return models.TaskSummary{}, err
		}

		switch status {
		case models.StatusPending:
			summary.Pending = count

		case models.StatusInProgress:
			summary.InProgress = count

		case models.StatusCompleted:
			summary.Completed = count
		}
	}

	if err := rows.Err(); err != nil {
		return models.TaskSummary{}, err
	}

	summary.Total = summary.Pending +
		summary.InProgress +
		summary.Completed

	return summary, nil
}
