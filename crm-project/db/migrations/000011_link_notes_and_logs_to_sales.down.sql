-- This migration reverts the changes by dropping the constraints and columns.
-- The order is important: drop constraints before columns.

-- Drop constraints from 'notes'
ALTER TABLE notes DROP CONSTRAINT IF EXISTS fk_note_deal;
ALTER TABLE notes DROP CONSTRAINT IF EXISTS fk_note_lead;
ALTER TABLE notes DROP CONSTRAINT IF EXISTS fk_note_contact;

-- Drop constraints from 'communication_logs'
ALTER TABLE communication_logs DROP CONSTRAINT IF EXISTS fk_comm_log_deal;
ALTER TABLE communication_logs DROP CONSTRAINT IF EXISTS fk_comm_log_lead;

-- Drop columns from 'notes'
ALTER TABLE notes DROP COLUMN IF EXISTS deal_id;
ALTER TABLE notes DROP COLUMN IF EXISTS lead_id;
ALTER TABLE notes DROP COLUMN IF EXISTS contact_id;

-- Drop columns from 'communication_logs'
ALTER TABLE communication_logs DROP COLUMN IF EXISTS deal_id;
ALTER TABLE communication_logs DROP COLUMN IF EXISTS lead_id;