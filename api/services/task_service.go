package services

import (
	"context"
	"time"

	"minitask/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
	// 1. Validar título
	title, err := validateTitle(input.Title)
	if err != nil {
		return nil, err
	}

	// 2. Obtener y validar descripción
	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	description, err = validateDescription(description)
	if err != nil {
		return nil, err
	}

	// 3. Obtener y validar estado
	status := models.StatusPending
	if input.Status != nil {
		status = *input.Status
	}

	if err := validateEnum(status, models.ValidStatuses, "El estado no es válido."); err != nil {
		return nil, err
	}

	// 4. Obtener y validar prioridad
	priority := models.PriorityMedium
	if input.Priority != nil {
		priority = *input.Priority
	}

	if err := validateEnum(priority, models.ValidPriorities, "La prioridad no es válida."); err != nil {
		return nil, err
	}

	// 5. Establecer fechas
	now := time.Now().UTC()

	// 6. Construir tarea
	task := models.Task{
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 7. Guardar mediante el repositorio
	if err := s.repo.Create(ctx, &task); err != nil {
		return nil, err
	}

	return &task, nil
}
