package handlers

import (
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"
    "strconv"
    "time"

    "crm-project/internal/models"
    "crm-project/internal/service"
    "crm-project/internal/util" // Added this line
    "github.com/go-chi/chi/v5"
)

type TaskHandler struct {
    taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
    return &TaskHandler{taskService: taskService}
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
    TaskName        string     `json:"task_name"`
    TaskDescription *string    `json:"task_description,omitempty"`
    DueDate         string     `json:"due_date"`
    Status          string     `json:"status"`
    AssignedTo      *int       `json:"assigned_to,omitempty"` // Added this line
    LeadID          *int       `json:"lead_id,omitempty"`
    DealID          *int       `json:"deal_id,omitempty"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
    TaskName        string     `json:"task_name"`
    TaskDescription *string    `json:"task_description,omitempty"`
    DueDate         string     `json:"due_date"`
    Status          string     `json:"status"`
    AssignedTo      *int       `json:"assigned_to,omitempty"` // Added this line
}

// TaskResponse represents the response structure for tasks
type TaskResponse struct {
    ID              int        `json:"id"`
    TaskName        string     `json:"task_name"`
    TaskDescription *string    `json:"task_description,omitempty"`
    DueDate         time.Time  `json:"due_date"`
    Status          string     `json:"status"`
    AssignedTo      int        `json:"assigned_to"`
    LeadID          *int       `json:"lead_id,omitempty"`
    DealID          *int       `json:"deal_id,omitempty"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// parseDueDate parses a date string in various formats into a time.Time object.
func parseDueDate(dateStr string) (time.Time, error) {
    // Try parsing MM/DD/YYYY (e.g., 10/01/2025)
    if t, err := time.Parse("01/02/2006", dateStr); err == nil {
        return t, nil
    }
    // Try parsing YYYY-MM-DD (e.g., 2025-09-28)
    if t, err := time.Parse("2006-01-02", dateStr); err == nil {
        return t, nil
    }
    // Try parsing YYYY-MM-DDTHH:MM (e.g., 2025-09-29T20:30)
    if t, err := time.Parse("2006-01-02T15:04", dateStr); err == nil {
        return t, nil
    }
    // Try parsing RFC3339 (e.g., 2025-09-28T00:00:00.000Z)
    if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
        return t, nil
    }
    return time.Time{}, fmt.Errorf("invalid due date format. Expected MM/DD/YYYY, YYYY-MM-DD, YYYY-MM-DDTHH:MM, or RFC3339")
}

// Helper function to convert models.Task to TaskResponse
func convertTaskToResponse(task *models.Task) TaskResponse {
    response := TaskResponse{
        ID:              task.ID,
        TaskName:        task.TaskName,
        TaskDescription: task.TaskDescription,
        DueDate:         task.DueDate,
        Status:          task.Status,
        AssignedTo:      task.AssignedTo,
        LeadID:          task.LeadID,
        DealID:          task.DealID,
        CreatedAt:       task.CreatedAt,
        UpdatedAt:       task.UpdatedAt,
    }
    return response
}

// GetAllTasks handles GET /api/v1/tasks
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    slog.Info("GetAllTasks called", "method", r.Method, "url", r.URL.Path)

    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    var filterUserID *int
    if claims.RoleID == util.RoleSalesAgent {
        filterUserID = &claims.UserID
    }

    tasks, err := h.taskService.GetAllTasks(r.Context(), filterUserID)
    if err != nil {
        slog.Error("Failed to get all tasks", "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get tasks: "+err.Error())
        return
    }

    taskResponses := make([]TaskResponse, len(tasks))
    for i, task := range tasks {
        taskResponses[i] = convertTaskToResponse(&task)
    }

    slog.Info("Successfully retrieved all tasks", "count", len(tasks))
    respondWithJSON(w, http.StatusOK, taskResponses)
}

// CreateTask handles POST /api/v1/tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    slog.Info("CreateTask called", "method", r.Method, "url", r.URL.Path)
    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    // ... (previous code)

    var req CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    parsedDueDate, err := parseDueDate(req.DueDate)
    if err != nil {
        slog.Error("Failed to parse due date", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid due date format. Expected YYYY-MM-DD")
        return
    }

    assignedToUserID := claims.UserID // Default to the user creating the task

    if claims.RoleID == util.RoleReception {
        if req.AssignedTo != nil {
            assignedToUserID = *req.AssignedTo
        }
    } else if claims.RoleID == util.RoleSalesAgent {
        // Sales agents can only assign tasks to themselves
        assignedToUserID = claims.UserID
    }

    task := &models.Task{
        TaskName:        req.TaskName,
        TaskDescription: req.TaskDescription,
        DueDate:         parsedDueDate,
        Status:          req.Status,
        AssignedTo:      assignedToUserID,
        LeadID:          req.LeadID,
        DealID:          req.DealID,
        CreatedAt:       time.Now(),
    }

    // ... (rest of the code)

    if _, err := h.taskService.CreateTask(r.Context(), task); err != nil {
        slog.Error("Failed to create task", "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to create task: "+err.Error())
        return
    }

    slog.Info("Successfully created task", "taskID", task.ID)
    respondWithJSON(w, http.StatusCreated, convertTaskToResponse(task))
}

// GetTaskByID handles GET /api/v1/tasks/{id}
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
    slog.Info("GetTaskByID called", "method", r.Method, "url", r.URL.Path)
    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        slog.Error("Invalid task ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid task ID")
        return
    }

    task, err := h.taskService.GetTaskByID(r.Context(), taskID)
    if err != nil {
        slog.Error("Failed to get task", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusNotFound, "Task not found: "+err.Error())
        return
    }

    if claims.RoleID == util.RoleSalesAgent && task.AssignedTo != claims.UserID {
        slog.Warn("Sales agent attempted to view unassigned task", "userID", claims.UserID, "taskID", taskID, "taskAssignedTo", task.AssignedTo)
        respondWithError(w, http.StatusForbidden, "You can only view tasks assigned to you")
        return
    }

    slog.Info("Successfully retrieved task", "taskID", task.ID)
    respondWithJSON(w, http.StatusOK, convertTaskToResponse(task))
}

// UpdateTask handles PUT /api/v1/tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    slog.Info("UpdateTask called", "method", r.Method, "url", r.URL.Path)
    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        slog.Error("Invalid task ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid task ID")
        return
    }

    var req UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    parsedDueDate, err := parseDueDate(req.DueDate)
    if err != nil {
        slog.Error("Failed to parse due date", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid due date format. Expected YYYY-MM-DD")
        return
    }

    existingTask, err := h.taskService.GetTaskByID(r.Context(), taskID)
    if err != nil {
        slog.Error("Task not found", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusNotFound, "Task not found: "+err.Error())
        return
    }

    assignedToUserID := existingTask.AssignedTo // Default to existing assigned user

    if claims.RoleID == util.RoleSalesAgent {
        if existingTask.AssignedTo != claims.UserID {
            slog.Warn("Sales agent attempted to update unassigned task", "userID", claims.UserID, "taskAssignedTo", existingTask.AssignedTo)
            respondWithError(w, http.StatusForbidden, "You can only update tasks assigned to you")
            return
        }
        // Sales agents cannot change who the task is assigned to
        if req.AssignedTo != nil && *req.AssignedTo != existingTask.AssignedTo {
            slog.Warn("Sales agent attempted to reassign task", "userID", claims.UserID, "taskAssignedTo", existingTask.AssignedTo, "requestedAssignedTo", *req.AssignedTo)
            respondWithError(w, http.StatusForbidden, "Sales agents cannot reassign tasks")
            return
        }
    } else if claims.RoleID == util.RoleReception {
        if req.AssignedTo != nil {
            assignedToUserID = *req.AssignedTo
        }
    }

    task := &models.Task{
        ID:              taskID,
        TaskName:        req.TaskName,
        TaskDescription: req.TaskDescription,
        DueDate:         parsedDueDate,
        Status:          req.Status,
        AssignedTo:      assignedToUserID,
    }

    if err := h.taskService.UpdateTask(r.Context(), task); err != nil {
        slog.Error("Failed to update task", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to update task: "+err.Error())
        return
    }

    updatedTask, err := h.taskService.GetTaskByID(r.Context(), taskID)
    if err != nil {
        slog.Error("Failed to fetch updated task", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get updated task: "+err.Error())
        return
    }

    slog.Info("Successfully updated task", "taskID", task.ID)
    respondWithJSON(w, http.StatusOK, convertTaskToResponse(updatedTask))
}

// DeleteTask handles DELETE /api/v1/tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    slog.Info("DeleteTask called", "method", r.Method, "url", r.URL.Path)
    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    if claims.RoleID == util.RoleSalesAgent {
        slog.Warn("Sales agent attempted to delete task", "userID", claims.UserID)
        respondWithError(w, http.StatusForbidden, "Sales agents are not allowed to delete tasks")
        return
    }

    taskID, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        slog.Error("Invalid task ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid task ID")
        return
    }

    // Check if the task exists before attempting to delete
    _, err = h.taskService.GetTaskByID(r.Context(), taskID)
    if err != nil {
        slog.Error("Task not found", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusNotFound, "Task not found: "+err.Error())
        return
    }

    if err := h.taskService.DeleteTask(r.Context(), taskID); err != nil {
        slog.Error("Failed to delete task", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to delete task: "+err.Error())
        return
    }

    slog.Info("Successfully deleted task", "taskID", taskID)
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

// GetTasksForUser handles GET /api/v1/users/{userId}/tasks
func (h *TaskHandler) GetTasksForUser(w http.ResponseWriter, r *http.Request) {
    slog.Info("GetTasksForUser called", "method", r.Method, "url", r.URL.Path)
    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }

    if claims.RoleID == util.RoleSalesAgent && claims.UserID != userID {
        slog.Warn("Sales agent attempted to view tasks for another user", "userID", claims.UserID, "requestedUserID", userID)
        respondWithError(w, http.StatusForbidden, "You can only view your own tasks")
        return
    }

    tasks, err := h.taskService.GetTasksForUser(r.Context(), userID)
    if err != nil {
        slog.Error("Failed to get user tasks", "userID", userID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get user tasks: "+err.Error())
        return
    }

    taskResponses := make([]TaskResponse, len(tasks))
    for i, task := range tasks {
        taskResponses[i] = convertTaskToResponse(&task)
    }

    slog.Info("Successfully retrieved user tasks", "userID", userID, "count", len(tasks))
    respondWithJSON(w, http.StatusOK, taskResponses)
}

// GetDealTasks handles GET /api/v1/deals/{dealId}/tasks
func (h *TaskHandler) GetDealTasks(w http.ResponseWriter, r *http.Request) {
    slog.Info("GetDealTasks called", "method", r.Method, "url", r.URL.Path)
    dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
    if err != nil {
        slog.Error("Invalid deal ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
        return
    }

    tasks, err := h.taskService.GetTasksByDealID(r.Context(), dealID)
    if err != nil {
        slog.Error("Failed to get deal tasks", "dealID", dealID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get deal tasks: "+err.Error())
        return
    }

    taskResponses := make([]TaskResponse, len(tasks))
    for i, task := range tasks {
        taskResponses[i] = convertTaskToResponse(&task)
    }

    slog.Info("Successfully retrieved deal tasks", "dealID", dealID, "count", len(tasks))
    respondWithJSON(w, http.StatusOK, taskResponses)
}

// CreateDealTask handles POST /api/v1/deals/{dealId}/tasks
func (h *TaskHandler) CreateDealTask(w http.ResponseWriter, r *http.Request) {
    slog.Info("CreateDealTask called", "method", r.Method, "url", r.URL.Path)
    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
    if err != nil {
        slog.Error("Invalid deal ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
        return
    }

    var req CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    parsedDueDate, err := parseDueDate(req.DueDate)
    if err != nil {
        slog.Error("Failed to parse due date", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid due date format. Expected YYYY-MM-DD")
        return
    }

    assignedToUserID := claims.UserID // Default to the user creating the task

    if claims.RoleID == util.RoleReception {
        if req.AssignedTo != nil {
            assignedToUserID = *req.AssignedTo
        }
    } else if claims.RoleID == util.RoleSalesAgent {
        // Sales agents can only assign tasks to themselves
        assignedToUserID = claims.UserID
    }

    task := &models.Task{
        TaskName:        req.TaskName,
        TaskDescription: req.TaskDescription,
        DueDate:         parsedDueDate,
        Status:          req.Status,
        AssignedTo:      assignedToUserID,
        DealID:          &dealID,
        CreatedAt:       time.Now(),
    }

    if _, err := h.taskService.CreateDealTask(r.Context(), task); err != nil {
        slog.Error("Failed to create deal task", "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to create deal task: "+err.Error())
        return
    }

    slog.Info("Successfully created deal task", "taskID", task.ID, "dealID", dealID)
    respondWithJSON(w, http.StatusCreated, convertTaskToResponse(task))
}

// UpdateDealTask handles PUT /api/v1/deals/{dealId}/tasks/{taskId}
func (h *TaskHandler) UpdateDealTask(w http.ResponseWriter, r *http.Request) {
    slog.Info("UpdateDealTask called", "method", r.Method, "url", r.URL.Path)
    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
    if err != nil {
        slog.Error("Invalid deal ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
        return
    }

    taskID, err := strconv.Atoi(chi.URLParam(r, "taskId"))
    if err != nil {
        slog.Error("Invalid task ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid task ID")
        return
    }

    var req UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    parsedDueDate, err := parseDueDate(req.DueDate)
    if err != nil {
        slog.Error("Failed to parse due date", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid due date format. Expected YYYY-MM-DD")
        return
    }

    existingTask, err := h.taskService.GetTaskByID(r.Context(), taskID)
    if err != nil {
        slog.Error("Task not found", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusNotFound, "Task not found: "+err.Error())
        return
    }

    assignedToUserID := existingTask.AssignedTo // Default to existing assigned user

    if claims.RoleID == util.RoleSalesAgent {
        if existingTask.AssignedTo != claims.UserID {
            slog.Warn("Sales agent attempted to update unassigned deal task", "userID", claims.UserID, "taskAssignedTo", existingTask.AssignedTo)
            respondWithError(w, http.StatusForbidden, "You can only update deal tasks assigned to you")
            return
        }
        // Sales agents cannot change who the task is assigned to
        if req.AssignedTo != nil && *req.AssignedTo != existingTask.AssignedTo {
            slog.Warn("Sales agent attempted to reassign deal task", "userID", claims.UserID, "taskAssignedTo", existingTask.AssignedTo, "requestedAssignedTo", *req.AssignedTo)
            respondWithError(w, http.StatusForbidden, "Sales agents cannot reassign deal tasks")
            return
        }
    } else if claims.RoleID == util.RoleReception {
        if req.AssignedTo != nil {
            assignedToUserID = *req.AssignedTo
        }
    }

    task := &models.Task{
        ID:              taskID,
        TaskName:        req.TaskName,
        TaskDescription: req.TaskDescription,
        DueDate:         parsedDueDate,
        Status:          req.Status,
        AssignedTo:      assignedToUserID,
        DealID:          &dealID,
    }

    if err := h.taskService.UpdateDealTask(r.Context(), task); err != nil {
        slog.Error("Failed to update deal task", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to update deal task: "+err.Error())
        return
    }

    updatedTask, err := h.taskService.GetTaskByID(r.Context(), taskID)
    if err != nil {
        slog.Error("Failed to fetch updated task", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get updated task: "+err.Error())
        return
    }

    slog.Info("Successfully updated deal task", "taskID", taskID, "dealID", dealID)
    respondWithJSON(w, http.StatusOK, convertTaskToResponse(updatedTask))
}

// DeleteDealTask handles DELETE /api/v1/deals/{dealId}/tasks/{taskId}
func (h *TaskHandler) DeleteDealTask(w http.ResponseWriter, r *http.Request) {
    slog.Info("DeleteDealTask called", "method", r.Method, "url", r.URL.Path)
    claims, ok := util.GetClaimsFromContext(r.Context())
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
    if err != nil {
        slog.Error("Invalid deal ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
        return
    }

    taskID, err := strconv.Atoi(chi.URLParam(r, "taskId"))
    if err != nil {
        slog.Error("Invalid task ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid task ID")
        return
    }

    existingTask, err := h.taskService.GetTaskByID(r.Context(), taskID)
    if err != nil {
        slog.Error("Task not found", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusNotFound, "Task not found: "+err.Error())
        return
    }

    if existingTask.DealID == nil || *existingTask.DealID != dealID {
        slog.Warn("Task does not belong to deal", "taskID", taskID, "dealID", dealID)
        respondWithError(w, http.StatusNotFound, "Task not found for this deal")
        return
    }

    if claims.RoleID == util.RoleSalesAgent {
        slog.Warn("Sales agent attempted to delete deal task", "userID", claims.UserID)
        respondWithError(w, http.StatusForbidden, "Sales agents are not allowed to delete deal tasks")
        return
    }

    // Reception can delete any deal task
    // No additional checks needed for reception as they have full control

    if err := h.taskService.DeleteDealTask(r.Context(), taskID); err != nil {
        slog.Error("Failed to delete deal task", "taskID", taskID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to delete deal task: "+err.Error())
        return
    }

    slog.Info("Successfully deleted deal task", "taskID", taskID, "dealID", dealID)
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Deal task deleted successfully"})
}