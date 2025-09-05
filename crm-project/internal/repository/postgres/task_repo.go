// Create new file: internal/repository/postgres/task_repo.go
package postgres

import (
	"context"
	"crm-project/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, t models.Task) (int, error) {
	var newID int
	query := `INSERT INTO tasks (task_name, task_description, due_date, status, assigned_to)
			  VALUES ($1, $2, $3, $4, $5) RETURNING task_id`
	err := r.db.QueryRowxContext(ctx, query, t.TaskName, t.TaskDescription, t.DueDate, t.Status, t.AssignedTo).Scan(&newID)
	return newID, err
}

func (r *TaskRepo) GetAll(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	query := `SELECT * FROM tasks ORDER BY due_date ASC`
	err := r.db.SelectContext(ctx, &tasks, query)
	return tasks, err
}

func (r *TaskRepo) GetByID(ctx context.Context, id int) (*models.Task, error) {
	var task models.Task
	query := `SELECT * FROM tasks WHERE task_id = $1`
	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) Update(ctx context.Context, t models.Task) error {
	query := `UPDATE tasks SET
				task_name = $1, task_description = $2, due_date = $3, status = $4,
				assigned_to = $5, updated_at = NOW()
			  WHERE task_id = $6`
	result, err := r.db.ExecContext(ctx, query, t.TaskName, t.TaskDescription, t.DueDate, t.Status, t.AssignedTo, t.ID)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return sql.ErrNoRows }
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE task_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return sql.ErrNoRows }
	return nil
}