package postgres

import (
    "database/sql"
    "fmt"
    "time"

    "crm-project/internal/models"
    "github.com/jmoiron/sqlx"
)

// TaskRepo implements the TaskRepository interface
type TaskRepo struct {
    db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepo {
    return &TaskRepo{db: db}
}

// CreateTask creates a new task
func (r *TaskRepo) CreateTask(task *models.Task) error {
    query := `
        INSERT INTO tasks (task_name, task_description, due_date, status, assigned_to, lead_id, deal_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING task_id, updated_at
    `

    currentTime := time.Now()
    var updatedAt time.Time
    err := r.db.QueryRowx(
        query,
        task.TaskName,
        task.TaskDescription,
        task.DueDate,
        task.Status,
        task.AssignedTo,
        task.LeadID,
        task.DealID,
        task.CreatedAt,
        currentTime,
    ).Scan(&task.ID, &updatedAt)

    if err != nil {
        return fmt.Errorf("failed to create task: %w", err)
    }

    task.UpdatedAt = &updatedAt
    return nil
}

// GetTaskByID retrieves a task by ID
func (r *TaskRepo) GetTaskByID(id int) (*models.Task, error) {
    query := `
        SELECT task_id, task_name, task_description, due_date, status, assigned_to, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM tasks
        WHERE task_id = $1 AND deleted_at IS NULL
    `

    var task models.Task
    var updatedAt sql.NullTime
    var deletedAt sql.NullTime
    err := r.db.Get(&task, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("task not found")
        }
        return nil, fmt.Errorf("failed to get task: %w", err)
    }

    if updatedAt.Valid {
        task.UpdatedAt = &updatedAt.Time
    }
    if deletedAt.Valid {
        task.DeletedAt = &deletedAt.Time
    }

    return &task, nil
}

// GetTasksByDealID retrieves all tasks for a specific deal
func (r *TaskRepo) GetTasksByDealID(dealID int) ([]models.Task, error) {
    query := `
        SELECT task_id, task_name, task_description, due_date, status, assigned_to, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM tasks
        WHERE deal_id = $1 AND deleted_at IS NULL
        ORDER BY due_date ASC
    `

    var tasks []models.Task
    err := r.db.Select(&tasks, query, dealID)
    if err != nil {
        return nil, fmt.Errorf("failed to get tasks by deal ID: %w", err)
    }

    return tasks, nil
}

func (r *TaskRepo) GetTasksByDealIDForUser(dealID int, userID int) ([]models.Task, error) {
    query := `
        SELECT task_id, task_name, task_description, due_date, status, assigned_to, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM tasks
        WHERE deal_id = $1 AND assigned_to = $2 AND deleted_at IS NULL
        ORDER BY due_date ASC
    `

    var tasks []models.Task
    err := r.db.Select(&tasks, query, dealID, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get tasks by deal ID for user: %w", err)
    }

    return tasks, nil
}

// GetAllTasks retrieves all tasks
func (r *TaskRepo) GetAllTasks() ([]models.Task, error) {
    query := `
        SELECT task_id, task_name, task_description, due_date, status, assigned_to, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM tasks
        WHERE deleted_at IS NULL
        ORDER BY due_date ASC
    `

    var tasks []models.Task
    err := r.db.Select(&tasks, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get all tasks: %w", err)
    }

    return tasks, nil
}

// UpdateTask updates an existing task
func (r *TaskRepo) UpdateTask(task *models.Task) error {
    query := `
        UPDATE tasks 
        SET task_name = $1, task_description = $2, due_date = $3, status = $4, assigned_to = $5, lead_id = $6, deal_id = $7, updated_at = $8
        WHERE task_id = $9 AND deleted_at IS NULL
        RETURNING updated_at
    `

    currentTime := time.Now()
    var updatedAt time.Time
    err := r.db.QueryRowx(
        query,
        task.TaskName,
        task.TaskDescription,
        task.DueDate,
        task.Status,
        task.AssignedTo,
        task.LeadID,
        task.DealID,
        currentTime,
        task.ID,
    ).Scan(&updatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("task not found")
        }
        return fmt.Errorf("failed to update task: %w", err)
    }

    task.UpdatedAt = &updatedAt
    return nil
}

// DeleteTask soft deletes a task
func (r *TaskRepo) DeleteTask(id int) error {
    query := `
        UPDATE tasks 
        SET deleted_at = $1
        WHERE task_id = $2 AND deleted_at IS NULL
    `

    result, err := r.db.Exec(query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("task not found")
    }

    return nil
}

// GetTasksForUser retrieves tasks for a specific user
func (r *TaskRepo) GetTasksForUser(userID int) ([]models.Task, error) {
    query := `
        SELECT task_id, task_name, task_description, due_date, status, assigned_to, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM tasks
        WHERE assigned_to = $1 AND deleted_at IS NULL
        ORDER BY due_date ASC
    `

    var tasks []models.Task
    err := r.db.Select(&tasks, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get tasks for user: %w", err)
    }

    return tasks, nil
}