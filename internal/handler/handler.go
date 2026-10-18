package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-checklist-pet/internal/service"
	"go-checklist-pet/internal/task"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service  *service.TaskService
	validate *validator.Validate
}

func NewHandler(service *service.TaskService) *Handler {
	return &Handler{service: service, validate: validator.New()}
}

func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskDTO task.TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(taskDTO); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	newTask := h.service.CreateTask(taskDTO)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		fmt.Println("failed to encode response:", err)
	}
}

func (h *Handler) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := h.service.ListTasks()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		fmt.Println("failed to encode response:", err)
	}
}

func (h *Handler) MarkTaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := taskIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	updatedTask, ok := h.service.MarkTaskDone(id)
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		fmt.Println("failed to encode response:", err)
	}
}

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := taskIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	ok := h.service.DeleteTask(id)
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func taskIDFromPath(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		return 0, errors.New("task ID must be a positive integer")
	}
	return id, nil
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /tasks", h.CreateTaskHandler)
	router.HandleFunc("GET /tasks", h.ListTasksHandler)
	router.HandleFunc("DELETE /tasks/{id}", h.DeleteTaskHandler)
	router.HandleFunc("PUT /tasks/{id}/done", h.MarkTaskDoneHandler)
	router.HandleFunc("GET /health", h.HealthHandler)
}
