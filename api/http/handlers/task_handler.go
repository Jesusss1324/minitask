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

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type TaskService interface {
	CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error)
	ListTasks(ctx context.Context, status, search string) ([]models.Task, error)
	GetTask(ctx context.Context, id uint64) (*models.Task, error)
	UpdateTask(ctx context.Context, id uint64, input models.UpdateTaskInput) (*models.Task, error)
	DeleteTask(ctx context.Context, id uint64) error
	GetSummary(ctx context.Context) (models.TaskSummary, error)
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{service: service}
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

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	tasks, err := h.service.ListTasks(r.Context(), status, search)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
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

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.handleError(w, &services.ValidationError{
			Message: "El identificador no es válido.",
		})
		return
	}

	var input models.UpdateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, &services.ValidationError{
			Message: "El cuerpo de la solicitud no es válido.",
		})
		return
	}

	task, err := h.service.UpdateTask(r.Context(), id, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.handleError(w, &services.ValidationError{
			Message: "El identificador no es válido.",
		})
		return
	}

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.GetSummary(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

func (h *TaskHandler) handleError(w http.ResponseWriter, err error) {
	var validationErr *services.ValidationError

	if errors.As(err, &validationErr) {
		writeError(
			w,
			http.StatusBadRequest,
			"validation_error",
			validationErr.Message,
		)
		return
	}

	if errors.Is(err, models.ErrTaskNotFound) {
		writeError(
			w,
			http.StatusNotFound,
			"not_found",
			"La tarea no existe.",
		)
		return
	}

	log.Printf("error interno: %v", err)

	writeError(
		w,
		http.StatusInternalServerError,
		"internal_error",
		"Error interno del servidor.",
	)
}

func writeError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	response := ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	}

	writeJSON(w, status, response)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error escribiendo respuesta JSON: %v", err)
	}
}
