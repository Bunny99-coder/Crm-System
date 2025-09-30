-- Add columns for lead/deal association and soft deletes
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS lead_id INT REFERENCES leads(lead_id);
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS deal_id INT REFERENCES deals(deal_id);
ALTER TABLE tasks ALTER COLUMN updated_at DROP NOT NULL;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX IF NOT EXISTS idx_tasks_lead_id ON tasks(lead_id);
CREATE INDEX IF NOT EXISTS idx_tasks_deal_id ON tasks(deal_id);
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);
