package services

import (
	"context"
	"strings"
	"time"

	"minitask/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	List(ctx context.Context, status, search string) ([]models.Task, error)
	GetByID(ctx context.Context, id uint64) (*models.Task, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	input models.CreateTaskInput,
) (*models.Task, error) {

	title, err := validateTitle(input.Title)
	if err != nil {
		return nil, err
	}

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	description, err = validateDescription(description)
	if err != nil {
		return nil, err
	}

	status := models.StatusPending
	if input.Status != nil {
		status = *input.Status
	}

	if err := validateEnum(
		status,
		models.ValidStatuses,
		"El estado no es válido.",
	); err != nil {
		return nil, err
	}

	priority := models.PriorityMedium
	if input.Priority != nil {
		priority = *input.Priority
	}

	if err := validateEnum(
		priority,
		models.ValidPriorities,
		"La prioridad no es válida.",
	); err != nil {
		return nil, err
	}

	now := time.Now().UTC().Truncate(time.Millisecond)

	task := models.Task{
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) ListTasks(
	ctx context.Context,
	status string,
	search string,
) ([]models.Task, error) {

	search = strings.TrimSpace(search)

	if status != "" {
		if err := validateEnum(
			status,
			models.ValidStatuses,
			"El estado no es válido.",
		); err != nil {
			return nil, err
		}
	}

	return s.repo.List(ctx, status, search)
}

func (s *TaskService) GetTask(
	ctx context.Context,
	id uint64,
) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}
