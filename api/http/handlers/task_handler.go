package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"minitask/models"
	"minitask/services"
)

type TaskService interface {
	CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error)
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

	// 1. Decodificar el JSON
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "El cuerpo de la solicitud no es válido.", http.StatusBadRequest)
		return
	}

	// 2. Ejecutar el servicio
	task, err := h.service.CreateTask(r.Context(), input)
	if err != nil {
		var validationErr *services.ValidationError

		if errors.As(err, &validationErr) {
			http.Error(w, validationErr.Message, http.StatusBadRequest)
			return
		}

		log.Printf("error creando tarea: %v", err)
		http.Error(w, "Error interno del servidor.", http.StatusInternalServerError)
		return
	}

	// 3. Configurar respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// 4. Escribir respuesta JSON
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("error escribiendo respuesta: %v", err)
	}
}
