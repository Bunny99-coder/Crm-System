-- This migration adds optional lead_id and deal_id columns to the notes
-- and communication_logs tables to create a unified activity timeline.

-- Add columns to the 'notes' table
ALTER TABLE notes ADD COLUMN contact_id INT;
ALTER TABLE notes ADD COLUMN lead_id INT;
ALTER TABLE notes ADD COLUMN deal_id INT;

-- Add columns to the 'communication_logs' table
ALTER TABLE communication_logs ADD COLUMN lead_id INT;
ALTER TABLE communication_logs ADD COLUMN deal_id INT;


-- Add foreign key constraints for the new columns.
-- We use ON DELETE SET NULL to preserve the activity history even if a lead/deal is deleted.

ALTER TABLE notes
ADD CONSTRAINT fk_note_contact
FOREIGN KEY (contact_id) REFERENCES contacts(contact_id) ON DELETE SET NULL;

ALTER TABLE notes
ADD CONSTRAINT fk_note_lead
FOREIGN KEY (lead_id) REFERENCES leads(lead_id) ON DELETE SET NULL;

ALTER TABLE notes
ADD CONSTRAINT fk_note_deal
FOREIGN KEY (deal_id) REFERENCES deals(deal_id) ON DELETE SET NULL;

ALTER TABLE communication_logs
ADD CONSTRAINT fk_comm_log_lead
FOREIGN KEY (lead_id) REFERENCES leads(lead_id) ON DELETE SET NULL;

ALTER TABLE communication_logs
ADD CONSTRAINT fk_comm_log_deal
FOREIGN KEY (deal_id) REFERENCES deals(deal_id) ON DELETE SET NULL;