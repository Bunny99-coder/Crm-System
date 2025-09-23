ALTER TABLE contacts ADD COLUMN created_by INT;

ALTER TABLE contacts ADD CONSTRAINT fk_contact_created_by
    FOREIGN KEY (created_by) REFERENCES users(user_id);