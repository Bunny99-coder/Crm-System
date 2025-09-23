-- Revert changes
DROP INDEX IF EXISTS idx_tasks_assigned_to;
DROP INDEX IF EXISTS idx_tasks_lead_id;
DROP INDEX IF EXISTS idx_tasks_deal_id;
DROP INDEX IF EXISTS idx_tasks_deleted_at;
ALTER TABLE tasks DROP COLUMN IF EXISTS lead_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS deal_id;
ALTER TABLE tasks DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE tasks ALTER COLUMN updated_at SET NOT NULL;
