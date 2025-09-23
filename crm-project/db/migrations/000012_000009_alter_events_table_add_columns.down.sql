-- Remove the added columns
ALTER TABLE events 
    DROP COLUMN IF EXISTS lead_id,
    DROP COLUMN IF EXISTS deal_id,
    DROP COLUMN IF EXISTS deleted_at;

-- Revert updated_at to NOT NULL
ALTER TABLE events 
    ALTER COLUMN updated_at SET NOT NULL;

-- Drop the indexes
DROP INDEX IF EXISTS idx_events_organizer_id;
DROP INDEX IF EXISTS idx_events_lead_id;
DROP INDEX IF EXISTS idx_events_deal_id;
DROP INDEX IF EXISTS idx_events_deleted_at;
