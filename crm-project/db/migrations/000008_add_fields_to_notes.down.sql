-- File: 20250919_add_fields_to_notes.down.sql

ALTER TABLE notes
DROP CONSTRAINT IF EXISTS fk_note_contact,
DROP CONSTRAINT IF EXISTS fk_note_lead,
DROP CONSTRAINT IF EXISTS fk_note_deal;

ALTER TABLE notes
DROP COLUMN IF EXISTS contact_id,
DROP COLUMN IF EXISTS lead_id,
DROP COLUMN IF EXISTS deal_id,
DROP COLUMN IF EXISTS updated_at;
