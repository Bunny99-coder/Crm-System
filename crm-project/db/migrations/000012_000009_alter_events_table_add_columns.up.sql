-- Add columns for lead/deal association and soft deletes
ALTER TABLE events 
    ADD COLUMN IF NOT EXISTS lead_id INT REFERENCES leads(lead_id),
    ADD COLUMN IF NOT EXISTS deal_id INT REFERENCES deals(deal_id),
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Make updated_at nullable
ALTER TABLE events 
    ALTER COLUMN updated_at DROP NOT NULL;

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_events_organizer_id ON events(organizer_id);
CREATE INDEX IF NOT EXISTS idx_events_lead_id ON events(lead_id);
CREATE INDEX IF NOT EXISTS idx_events_deal_id ON events(deal_id);
CREATE INDEX IF NOT EXISTS idx_events_deleted_at ON events(deleted_at);
