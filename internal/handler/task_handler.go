package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/marcosx3/tasker/internal/domain"
	"github.com/marcosx3/tasker/internal/dto"
)

type TaskService interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetAllTasks(ctx context.Context) ([]domain.Task, error)
	GetOneTaskByID(ctx context.Context, id string) (domain.Task, error)
	UpdateTask(ctx context.Context, id string, t domain.Task) (domain.Task, error)
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) PostTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo Inválido", http.StatusBadRequest)
	}
	var request dto.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 1. DTO -> Domain (chama domain.NewTask internamente)
	taskDomain, err := request.ToDomain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Handler -> Service
	createdTask, err := h.service.CreateTask(r.Context(), taskDomain)
	if err != nil {
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}

	// 3. Resposta HTTP 201 Created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo Inválido", http.StatusBadRequest)
	}
	tasks, err := h.service.GetAllTasks(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo Inválido", http.StatusBadRequest)
	}
	id := r.PathValue("id")

	task, err := h.service.GetOneTaskByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	// 1. Validar método HTTP
	if r.Method != http.MethodPut {
		http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	// 2. Recuperar ID da URL
	id := r.PathValue("id")

	// 3. Decodificar JSON
	var request dto.TaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 4. Converter DTO -> Domain
	taskDomain, err := request.ToDomain()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 5. Atualizar tarefa através do service
	updatedTask, err := h.service.UpdateTask(
		r.Context(),
		id,
		taskDomain,
	)

	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// 6. Retornar resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(updatedTask)
}
