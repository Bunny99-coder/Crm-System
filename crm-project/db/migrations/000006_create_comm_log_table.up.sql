CREATE TABLE IF NOT EXISTS communication_logs (
    log_id SERIAL PRIMARY KEY,
    contact_id INT NOT NULL,
    user_id INT NOT NULL, -- The user who logged the interaction
    interaction_date TIMESTAMPTZ NOT NULL,
    interaction_type VARCHAR(50) NOT NULL, -- e.g., Call, Email, Meeting
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_comm_log_contact
        FOREIGN KEY(contact_id)
        REFERENCES contacts(contact_id),
    CONSTRAINT fk_comm_log_user
        FOREIGN KEY(user_id)
        REFERENCES users(user_id)
);