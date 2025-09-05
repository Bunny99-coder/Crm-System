// File: internal/repository/postgres/user_repo.go
package postgres

import (
	"crm-project/internal/models"
	"github.com/jmoiron/sqlx"
	    "database/sql"
			"strconv"
			"context"


)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user into the database.
func (r *UserRepo) Create(ctx context.Context, user models.User) (int, error) {
	var newUserID int
	query := `INSERT INTO users (username, password_hash, email, role_id)
			  VALUES ($1, $2, $3, $4)
			  RETURNING user_id`

	err := r.db.QueryRowxContext(ctx, query, user.Username, user.PasswordHash, user.Email, user.RoleID).Scan(&newUserID)
	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// Add to internal/repository/postgres/user_repo.go

// GetAll retrieves all users from the database.
func (r *UserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `SELECT user_id, username, email, role_id, created_at, updated_at FROM users ORDER BY created_at DESC`
	
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Add to internal/repository/postgres/user_repo.go
// You will need to add an import for "database/sql" if it's not there.

// GetByID retrieves a single user by their ID.
func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `SELECT user_id, username, email, role_id, created_at, updated_at FROM users WHERE user_id = $1`
	
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}


// Add to internal/repository/postgres/user_repo.go

// Update modifies an existing user in the database.
func (r *UserRepo) Update(ctx context.Context, user models.User) error {
	query := `UPDATE users SET
				username = $1,
				email = $2,
				role_id = $3,
				updated_at = NOW()`
	
	args := []interface{}{user.Username, user.Email, user.RoleID}

	// Conditionally update the password only if a new one is provided
	if user.PasswordHash != "" {
		query += ", password_hash = $4"
		args = append(args, user.PasswordHash)
	}

	query += " WHERE user_id = $" + strconv.Itoa(len(args)+1)
	args = append(args, user.ID)
	
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}


// Add to internal/repository/postgres/user_repo.go

// Delete removes a user from the database by their ID.
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}


// Add this method to the end of your internal/repository/postgres/user_repo.go file

// GetAllSalesAgents retrieves all users with the 'Sales_Agent' role.
func (r *UserRepo) GetAllSalesAgents(ctx context.Context) ([]models.User, error) {
	var agents []models.User
	// We are assuming role_id 1 is 'Sales_Agent' based on our first migration.
	query := `SELECT user_id, username, email FROM users WHERE role_id = 1 ORDER BY username`
	
	err := r.db.SelectContext(ctx,&agents, query)
	if err != nil {
		return nil, err
	}
	
	return agents, nil
}


// GetByUsername retrieves a single user by their username.
// It's important to select the password_hash for comparison.
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT user_id, username, email, role_id, password_hash FROM users WHERE username = $1`
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil for not found
		}
		return nil, err
	}
	return &user, nil
}