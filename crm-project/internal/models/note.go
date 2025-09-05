// File: internal/models/note.go
package models

import "time"

// Note represents a simple text note created by a user.
type Note struct {
	ID        int       `db:"note_id"   json:"id"`
	UserID    int       `db:"user_id"   json:"user_id"`
	NoteDate  time.Time `db:"note_date" json:"note_date"`
	NoteText  string    `db:"note_text" json:"note_text"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}