// Create new file: internal/repository/postgres/note_repo.go
package postgres

import (
	"context"
	"crm-project/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type NoteRepo struct {
	db *sqlx.DB
}

func NewNoteRepo(db *sqlx.DB) *NoteRepo {
	return &NoteRepo{db: db}
}

func (r *NoteRepo) Create(ctx context.Context, n models.Note) (int, error) {
	var newID int
	query := `INSERT INTO notes (user_id, note_text, note_date) VALUES ($1, $2, $3) RETURNING note_id`
	err := r.db.QueryRowxContext(ctx, query, n.UserID, n.NoteText, n.NoteDate).Scan(&newID)
	return newID, err
}

// A function to get all notes by a specific user might be useful.
func (r *NoteRepo) GetAllByUser(ctx context.Context, userID int) ([]models.Note, error) {
	var notes []models.Note
	query := `SELECT * FROM notes WHERE user_id = $1 ORDER BY note_date DESC`
	err := r.db.SelectContext(ctx, &notes, query, userID)
	return notes, err
}



// in internal/repository/postgres/note_repo.go
func (r *NoteRepo) GetAll(ctx context.Context) ([]models.Note, error) {
	var notes []models.Note
	query := `SELECT * FROM notes ORDER BY note_date DESC`
	err := r.db.SelectContext(ctx, &notes, query)
	return notes, err
}


func (r *NoteRepo) GetByID(ctx context.Context, id int) (*models.Note, error) {
	var note models.Note
	query := `SELECT * FROM notes WHERE note_id = $1`
	err := r.db.GetContext(ctx, &note, query, id)
	if err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepo) Update(ctx context.Context, n models.Note) error {
	query := `UPDATE notes SET note_text = $1, note_date = $2 WHERE note_id = $3`
	result, err := r.db.ExecContext(ctx, query, n.NoteText, n.NoteDate, n.ID)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return sql.ErrNoRows }
	return nil
}

func (r *NoteRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM notes WHERE note_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return sql.ErrNoRows }
	return nil
}