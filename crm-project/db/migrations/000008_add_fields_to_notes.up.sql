-- File: 20250919_add_fields_to_notes.up.sql

ALTER TABLE notes
ADD COLUMN contact_id INT NULL,
ADD COLUMN lead_id INT NULL,
ADD COLUMN deal_id INT NULL,
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Optional: if you want foreign keys
ALTER TABLE notes
ADD CONSTRAINT fk_note_contact FOREIGN KEY(contact_id) REFERENCES contacts(contact_id),
ADD CONSTRAINT fk_note_lead FOREIGN KEY(lead_id) REFERENCES leads(lead_id),
ADD CONSTRAINT fk_note_deal FOREIGN KEY(deal_id) REFERENCES deals(deal_id);
