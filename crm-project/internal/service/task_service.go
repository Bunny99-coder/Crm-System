package service

import (
    "errors"
    "fmt"

    "crm-project/internal/models"
    "crm-project/internal/repository/postgres"
)

type TaskService struct {
    taskRepo postgres.TaskRepository
}

func NewTaskService(taskRepo postgres.TaskRepository) *TaskService {
    return &TaskService{taskRepo: taskRepo}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(task *models.Task) error {
    if task.TaskName == "" {
        return errors.New("task name cannot be empty")
    }
    if task.AssignedTo == 0 {
        return errors.New("assigned_to ID is required")
    }
    if task.DueDate.IsZero() {
        return errors.New("due date is required")
    }
    if task.Status == "" {
        task.Status = "Pending" // Default status
    }

    return s.taskRepo.CreateTask(task)
}

// GetTaskByID retrieves a task by ID
func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
    if id <= 0 {
        return nil, errors.New("invalid task ID")
    }

    return s.taskRepo.GetTaskByID(id)
}

// GetTasksByDealID retrieves all tasks for a specific deal
func (s *TaskService) GetTasksByDealID(dealID int) ([]models.Task, error) {
    if dealID <= 0 {
        return nil, errors.New("invalid deal ID")
    }

    return s.taskRepo.GetTasksByDealID(dealID)
}

// GetAllTasks retrieves all tasks
func (s *TaskService) GetAllTasks() ([]models.Task, error) {
    return s.taskRepo.GetAllTasks()
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(task *models.Task) error {
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
        return fmt.Errorf("failed to verify task existence: %w", err)
    }

    if existingTask.AssignedTo != task.AssignedTo {
        return errors.New("unauthorized to update this task")
    }

    return s.taskRepo.UpdateTask(task)
}

// DeleteTask soft deletes a task
func (s *TaskService) DeleteTask(id int) error {
    if id <= 0 {
        return errors.New("invalid task ID")
    }

    return s.taskRepo.DeleteTask(id)
}

// GetTasksForUser retrieves tasks for a specific user
func (s *TaskService) GetTasksForUser(userID int) ([]models.Task, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user ID")
    }

    return s.taskRepo.GetTasksForUser(userID)
}

// CreateDealTask creates a new task for a deal (nested route)
func (s *TaskService) CreateDealTask(task *models.Task) error {
    if task.TaskName == "" {
        return errors.New("task name cannot be empty")
    }
    if task.AssignedTo == 0 {
        return errors.New("assigned_to ID is required")
    }
    if task.DealID == nil || *task.DealID <= 0 {
        return errors.New("deal ID is required")
    }
    if task.DueDate.IsZero() {
        return errors.New("due date is required")
    }
    if task.Status == "" {
        task.Status = "Pending"
    }

    return s.taskRepo.CreateTask(task)
}

// UpdateDealTask updates a task for a deal
func (s *TaskService) UpdateDealTask(task *models.Task) error {
    if task.ID <= 0 {
        return errors.New("invalid task ID")
    }
    if task.TaskName == "" {
        return errors.New("task name cannot be empty")
    }

    existingTask, err := s.taskRepo.GetTaskByID(task.ID)
    if err != nil {
        return fmt.Errorf("failed to verify task existence: %w", err)
    }

    if existingTask.DealID == nil || *existingTask.DealID != *task.DealID {
        return errors.New("task not found for this deal")
    }

    return s.taskRepo.UpdateTask(task)
}

// DeleteDealTask deletes a task for a deal
func (s *TaskService) DeleteDealTask(id int) error {
    if id <= 0 {
        return errors.New("invalid task ID")
    }

    return s.taskRepo.DeleteTask(id)
}