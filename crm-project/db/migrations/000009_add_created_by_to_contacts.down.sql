ALTER TABLE contacts DROP CONSTRAINT IF EXISTS fk_contact_created_by;
ALTER TABLE contacts DROP COLUMN IF EXISTS created_by;