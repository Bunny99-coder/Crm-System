ALTER TABLE deals ADD COLUMN created_by INT;

ALTER TABLE deals ADD CONSTRAINT fk_deal_created_by
    FOREIGN KEY (created_by) REFERENCES users(user_id);