// Replace the contents of internal/service/task_service.go
package service

import (
	"context"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type TaskService struct {
	taskRepo *postgres.TaskRepo
	userRepo *postgres.UserRepo
	logger   *slog.Logger
}

func NewTaskService(tr *postgres.TaskRepo, ur *postgres.UserRepo, logger *slog.Logger) *TaskService {
	return &TaskService{taskRepo: tr, userRepo: ur, logger: logger}
}

func (s *TaskService) CreateTask(ctx context.Context, t models.Task) (int, error) {
	if t.TaskName == "" {
		return 0, errors.New("task name is required")
	}
	// in task_service.go -> CreateTask

	if t.DueDate.IsZero() {
		return 0, errors.New("due date is required")
	}

	// Get the current time
	now := time.Now()
	// Get the beginning of today (Year, Month, Day, at 00:00:00) in the server's location
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Check if the due date is before the start of today
	if t.DueDate.Before(today) {
		return 0, errors.New("due date cannot be in the past")
	}
	if _, err := s.userRepo.GetByID(ctx, t.AssignedTo); err != nil {
		return 0, fmt.Errorf("invalid assigned_to user id: %d", t.AssignedTo)
	}
	// Default status if not provided
	if t.Status == "" {
		t.Status = "Pending"
	}
	return s.taskRepo.Create(ctx, t)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.taskRepo.GetAll(ctx)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task with ID %d not found", id)
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, t models.Task) error {
	_, err := s.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}
	t.ID = id
	// Add validation...
	if t.TaskName == "" {
		return errors.New("task name cannot be empty on update")
	}
	err = s.taskRepo.Update(ctx, t)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with ID %d not found during update", id)
		}
		return err
	}
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	err := s.taskRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with ID %d not found", id)
		}
		return err
	}
	return nil
}
