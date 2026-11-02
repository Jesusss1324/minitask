package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"minitask/models"
	"minitask/services"
)

type TaskService interface {
	CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error)
	ListTasks(ctx context.Context, status, search string) ([]models.Task, error)
	GetTask(ctx context.Context, id uint64) (*models.Task, error)
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, &services.ValidationError{
			Message: "El cuerpo de la solicitud no es válido.",
		})
		return
	}

	task, err := h.service.CreateTask(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("error escribiendo respuesta: %v", err)
	}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	tasks, err := h.service.ListTasks(r.Context(), status, search)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("error escribiendo respuesta: %v", err)
	}
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.handleError(w, &services.ValidationError{
			Message: "El identificador no es válido.",
		})
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("error escribiendo respuesta: %v", err)
	}
}

func (h *TaskHandler) handleError(w http.ResponseWriter, err error) {
	var validationErr *services.ValidationError

	if errors.As(err, &validationErr) {
		http.Error(w, validationErr.Message, http.StatusBadRequest)
		return
	}

	if errors.Is(err, models.ErrTaskNotFound) {
		http.Error(w, "La tarea no existe.", http.StatusNotFound)
		return
	}

	log.Printf("error interno: %v", err)
	http.Error(
		w,
		"Error interno del servidor.",
		http.StatusInternalServerError,
	)
}
