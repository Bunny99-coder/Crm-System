CREATE TABLE IF NOT EXISTS notes (
    note_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL, -- The user who created the note
    note_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_note_user
        FOREIGN KEY(user_id)
        REFERENCES users(user_id)
);