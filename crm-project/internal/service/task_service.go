package service

import (
	"context"
	"crm-project/internal/config" // Import config
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"crm-project/internal/util" // Import util for claims
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type TaskService struct {
	taskRepo postgres.TaskRepository
	cfg      *config.Config // Add config here
	logger   *slog.Logger
}

func NewTaskService(taskRepo postgres.TaskRepository, cfg *config.Config, logger *slog.Logger) *TaskService {
	return &TaskService{taskRepo: taskRepo, cfg: cfg, logger: logger}
}

// CreateTask creates a new task with role-based assignment
func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) (int, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return 0, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can create tasks.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for CreateTask", "user_id", claims.UserID, "role_id", claims.RoleID)
		return 0, fmt.Errorf("forbidden: only managers can create tasks")
	}

	if task.TaskName == "" {
		return 0, errors.New("task name cannot be empty")
	}
	if task.DueDate.IsZero() {
		return 0, errors.New("due date is required")
	}
	if task.Status == "" {
		task.Status = "Pending" // Default status
	}

	// Set the creator of the task to the currently logged-in user's ID.
	task.CreatedBy = claims.UserID

	// --- ROLE-BASED ASSIGNMENT ---
	// If manager, they can assign to anyone.
	// If sales agent, they can only assign to themselves.
	if claims.RoleID == s.cfg.Roles.SalesAgentID && task.AssignedTo != claims.UserID {
		s.logger.Warn("Sales agent tried to assign task to another user", "user_id", claims.UserID, "assigned_to", task.AssignedTo)
		return 0, errors.New("sales agents can only assign tasks to themselves")
	}
	if task.AssignedTo == 0 { // If not explicitly assigned, assign to creator
		task.AssignedTo = claims.UserID
	}

	err := s.taskRepo.CreateTask(task)
	if err != nil {
		s.logger.Error("failed to create task in repository", "error", err)
		return 0, fmt.Errorf("failed to create task: %w", err)
	}
	return task.ID, nil
}

// GetTaskByID retrieves a task by ID with permission check
func (s *TaskService) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	if id <= 0 {
		return nil, errors.New("invalid task ID")
	}

	task, err := s.taskRepo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}

	// --- PERMISSION CHECK ---
	// A user can view if they are a Receptionist OR if they are the assigned sales agent.
	isAllowed := claims.RoleID == s.cfg.Roles.ReceptionID || (task.AssignedTo == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for GetTaskByID", "user_id", claims.UserID, "role_id", claims.RoleID, "task_id", id, "task_assigned_to", task.AssignedTo)
		return nil, fmt.Errorf("forbidden: you do not have permission to view this task")
	}

	return task, nil
}

// GetTasksByDealID retrieves all tasks for a specific deal with permission check
func (s *TaskService) GetTasksByDealID(ctx context.Context, dealID int) ([]models.Task, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	if dealID <= 0 {
		return nil, errors.New("invalid deal ID")
	}

	// If sales agent, filter by assigned_to. If manager, get all for deal.
	if claims.RoleID == s.cfg.Roles.SalesAgentID {
		return s.taskRepo.GetTasksByDealIDForUser(dealID, claims.UserID)
	}
	return s.taskRepo.GetTasksByDealID(dealID)
}

// GetAllTasks retrieves all tasks with permission check
func (s *TaskService) GetAllTasks(ctx context.Context, assignedToUserID *int) ([]models.Task, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// If the user is a Sales Agent, only show tasks assigned to them.
	if claims.RoleID == s.cfg.Roles.SalesAgentID {
		s.logger.Debug("fetching tasks for single sales agent", "user_id", claims.UserID)
		return s.taskRepo.GetTasksForUser(claims.UserID)
	}

	// If assignedToUserID is provided, filter by that user.
	if assignedToUserID != nil && *assignedToUserID > 0 {
		s.logger.Debug("fetching tasks for specific user", "assigned_to_user_id", *assignedToUserID)
		return s.taskRepo.GetTasksForUser(*assignedToUserID)
	}

	// Otherwise (for Reception/Manager), show all tasks.
	s.logger.Debug("fetching all tasks for manager role", "user_id", claims.UserID)
	return s.taskRepo.GetAllTasks()
}

// UpdateTask updates an existing task with permission check
func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return errors.New("could not retrieve user claims from context")
	}

	if task.ID <= 0 {
		return errors.New("invalid task ID")
	}
	if task.TaskName == "" {
		return errors.New("task name cannot be empty")
	}
	if task.AssignedTo == 0 {
		return errors.New("assigned_to ID is required")
	}

	existingTask, err := s.taskRepo.GetTaskByID(task.ID)
	if err != nil {
		s.logger.Error("Failed to fetch existing task for update", "task_id", task.ID, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with ID %d not found", task.ID)
		}
		return err
	}
	if existingTask == nil {
		s.logger.Warn("Existing task not found", "task_id", task.ID)
		return fmt.Errorf("task with ID %d not found", task.ID)
	}

	// --- PERMISSION CHECK ---
	// A user can update if they are a Receptionist OR if they are the assigned sales agent.
	isAllowed := claims.RoleID == s.cfg.Roles.ReceptionID || (existingTask.AssignedTo == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for UpdateTask", "user_id", claims.UserID, "role_id", claims.RoleID, "task_id", task.ID, "task_assigned_to", existingTask.AssignedTo)
		return fmt.Errorf("forbidden: you do not have permission to update this task")
	}

	// Sales agents can only update tasks assigned to them. Managers can reassign.
	if claims.RoleID == s.cfg.Roles.SalesAgentID && task.AssignedTo != claims.UserID {
		s.logger.Warn("Sales agent tried to reassign task to another user", "user_id", claims.UserID, "assigned_to", task.AssignedTo)
		return errors.New("sales agents cannot reassign tasks")
	}

	return s.taskRepo.UpdateTask(task)
}

// DeleteTask soft deletes a task with permission check
func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return errors.New("could not retrieve user claims from context")
	}

	if id <= 0 {
		return errors.New("invalid task ID")
	}

	existingTask, err := s.taskRepo.GetTaskByID(id)
	if err != nil {
		s.logger.Error("Failed to fetch existing task for deletion", "task_id", id, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with ID %d not found", id)
		}
		return err
	}
	if existingTask == nil {
		s.logger.Warn("Existing task not found", "task_id", id)
		return fmt.Errorf("task with ID %d not found", id)
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can delete tasks.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for DeleteTask", "user_id", claims.UserID, "role_id", claims.RoleID, "task_id", id)
		return fmt.Errorf("forbidden: only managers can delete tasks")
	}

	return s.taskRepo.DeleteTask(id)
}

// GetTasksForUser retrieves tasks for a specific user (used internally or by manager)
func (s *TaskService) GetTasksForUser(ctx context.Context, userID int) ([]models.Task, error) {
	return s.taskRepo.GetTasksForUser(userID)
}

// CreateDealTask creates a new task for a deal (nested route) with role-based assignment
func (s *TaskService) CreateDealTask(ctx context.Context, task *models.Task) (int, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return 0, errors.New("could not retrieve user claims from context")
	}

	if task.TaskName == "" {
		return 0, errors.New("task name cannot be empty")
	}
	if task.DealID == nil || *task.DealID <= 0 {
		return 0, errors.New("deal ID is required")
	}
	if task.DueDate.IsZero() {
		return 0, errors.New("due date is required")
	}
	if task.Status == "" {
		task.Status = "Pending"
	}

	// Set the creator of the task to the currently logged-in user's ID.
	task.CreatedBy = claims.UserID

	// --- ROLE-BASED ASSIGNMENT ---
	// If manager, they can assign to anyone.
	// If sales agent, they can only assign to themselves.
	if claims.RoleID == s.cfg.Roles.SalesAgentID && task.AssignedTo != claims.UserID {
		s.logger.Warn("Sales agent tried to assign deal task to another user", "user_id", claims.UserID, "assigned_to", task.AssignedTo)
		return 0, errors.New("sales agents can only assign deal tasks to themselves")
	}
	if task.AssignedTo == 0 { // If not explicitly assigned, assign to creator
		task.AssignedTo = claims.UserID
	}

	err := s.taskRepo.CreateTask(task)
	if err != nil {
		s.logger.Error("Failed to create deal task in repository", "error", err)
		return 0, fmt.Errorf("failed to create deal task: %w", err)
	}
	return task.ID, nil
}

// UpdateDealTask updates a task for a deal with permission check
func (s *TaskService) UpdateDealTask(ctx context.Context, task *models.Task) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return errors.New("could not retrieve user claims from context")
	}

	if task.ID <= 0 {
		return errors.New("invalid task ID")
	}
	if task.TaskName == "" {
		return errors.New("task name cannot be empty")
	}

	existingTask, err := s.taskRepo.GetTaskByID(task.ID)
	if err != nil {
		s.logger.Error("Failed to fetch existing deal task for update", "task_id", task.ID, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("deal task with ID %d not found", task.ID)
		}
		return err
	}
	if existingTask == nil {
		s.logger.Warn("Existing deal task not found", "task_id", task.ID)
		return fmt.Errorf("deal task with ID %d not found", task.ID)
	}

	if existingTask.DealID == nil || *existingTask.DealID != *task.DealID {
		return errors.New("task not found for this deal")
	}

	// --- PERMISSION CHECK ---
	// A user can update if they are a Receptionist OR if they are the assigned sales agent.
	isAllowed := claims.RoleID == s.cfg.Roles.ReceptionID || (existingTask.AssignedTo == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for UpdateDealTask", "user_id", claims.UserID, "role_id", claims.RoleID, "task_id", task.ID, "task_assigned_to", existingTask.AssignedTo)
		return fmt.Errorf("forbidden: you do not have permission to update this deal task")
	}

	// Sales agents can only update tasks assigned to them. Managers can reassign.
	if claims.RoleID == s.cfg.Roles.SalesAgentID && task.AssignedTo != claims.UserID {
		s.logger.Warn("Sales agent tried to reassign deal task to another user", "user_id", claims.UserID, "assigned_to", task.AssignedTo)
		return errors.New("sales agents cannot reassign deal tasks")
	}

	return s.taskRepo.UpdateTask(task)
}

// DeleteDealTask deletes a task for a deal with permission check
func (s *TaskService) DeleteDealTask(ctx context.Context, id int) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return errors.New("could not retrieve user claims from context")
	}

	if id <= 0 {
		return errors.New("invalid task ID")
	}

	existingTask, err := s.taskRepo.GetTaskByID(id)
	if err != nil {
		s.logger.Error("Failed to fetch existing deal task for deletion", "task_id", id, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("deal task with ID %d not found", id)
		}
		return err
	}
	if existingTask == nil {
		s.logger.Warn("Existing deal task not found", "task_id", id)
		return fmt.Errorf("deal task with ID %d not found", id)
	}

	// --- PERMISSION CHECK ---
	// A user can delete if they are a Receptionist OR if they are the assigned sales agent.
	isAllowed := claims.RoleID == s.cfg.Roles.ReceptionID || (existingTask.AssignedTo == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for DeleteDealTask", "user_id", claims.UserID, "role_id", claims.RoleID, "task_id", id, "task_assigned_to", existingTask.AssignedTo)
		return fmt.Errorf("forbidden: you do not have permission to delete this deal task")
	}

	return s.taskRepo.DeleteTask(id)
}
