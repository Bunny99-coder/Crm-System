-- Add optional columns to link tasks and events directly to leads and deals
ALTER TABLE tasks ADD COLUMN lead_id INT;
ALTER TABLE tasks ADD COLUMN deal_id INT;
ALTER TABLE events ADD COLUMN lead_id INT;
ALTER TABLE events ADD COLUMN deal_id INT;

-- Add the foreign key constraints
ALTER TABLE tasks ADD CONSTRAINT fk_task_lead FOREIGN KEY (lead_id) REFERENCES leads(lead_id) ON DELETE SET NULL;
ALTER TABLE tasks ADD CONSTRAINT fk_task_deal FOREIGN KEY (deal_id) REFERENCES deals(deal_id) ON DELETE SET NULL;
ALTER TABLE events ADD CONSTRAINT fk_event_lead FOREIGN KEY (lead_id) REFERENCES leads(lead_id) ON DELETE SET NULL;
ALTER TABLE events ADD CONSTRAINT fk_event_deal FOREIGN KEY (deal_id) REFERENCES deals(deal_id) ON DELETE SET NULL;