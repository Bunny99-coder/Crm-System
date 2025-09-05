// Replace the contents of internal/api/handlers/task_handler.go
package handlers

import (
	"crm-project/internal/models"
	"crm-project/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	service *service.TaskService
	logger  *slog.Logger
}

func NewTaskHandler(s *service.TaskService, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{service: s, logger: logger}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		h.logger.Warn("invalid create task request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newID, err := h.service.CreateTask(ctx, newTask)
	if err != nil {
		h.logger.Error("failed to create task", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Info("task created successfully", "task_id", newID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := h.service.GetAllTasks(ctx)
	if err != nil {
		h.logger.Error("failed to get all tasks", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Debug("retrieved all tasks", "count", len(tasks))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	task, err := h.service.GetTaskByID(ctx, id)
	if err != nil {
		h.logger.Warn("task not found", "task_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved task by id", "task_id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	var taskToUpdate models.Task
	if err := json.NewDecoder(r.Body).Decode(&taskToUpdate); err != nil {
		h.logger.Warn("invalid update task request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateTask(ctx, id, taskToUpdate)
	if err != nil {
		h.logger.Error("failed to update task", "task_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("task updated successfully", "task_id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteTask(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete task", "task_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("task deleted successfully", "task_id", id)
	w.WriteHeader(http.StatusNoContent)
}