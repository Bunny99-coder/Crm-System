-- Revert changes
DROP INDEX IF EXISTS idx_comm_logs_contact_id;
DROP INDEX IF EXISTS idx_comm_logs_user_id;
DROP INDEX IF EXISTS idx_comm_logs_lead_id;
DROP INDEX IF EXISTS idx_comm_logs_deal_id;
DROP INDEX IF EXISTS idx_comm_logs_deleted_at;
ALTER TABLE communication_logs DROP COLUMN IF EXISTS lead_id;
ALTER TABLE communication_logs DROP COLUMN IF EXISTS deal_id;
ALTER TABLE communication_logs DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE communication_logs ALTER COLUMN created_at SET NOT NULL;
