-- Add columns for lead/deal association and soft deletes
ALTER TABLE communication_logs ADD COLUMN IF NOT EXISTS lead_id INT REFERENCES leads(lead_id);
ALTER TABLE communication_logs ADD COLUMN IF NOT EXISTS deal_id INT REFERENCES deals(deal_id);
ALTER TABLE communication_logs ALTER COLUMN created_at DROP NOT NULL;
ALTER TABLE communication_logs ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_comm_logs_contact_id ON communication_logs(contact_id);
CREATE INDEX IF NOT EXISTS idx_comm_logs_user_id ON communication_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_comm_logs_lead_id ON communication_logs(lead_id);
CREATE INDEX IF NOT EXISTS idx_comm_logs_deal_id ON communication_logs(deal_id);
CREATE INDEX IF NOT EXISTS idx_comm_logs_deleted_at ON communication_logs(deleted_at);
