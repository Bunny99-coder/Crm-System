CREATE TABLE IF NOT EXISTS lead_sources (
    source_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS lead_statuses (
    status_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS leads (
    lead_id SERIAL PRIMARY KEY,
    contact_id INT NOT NULL,
    property_id INT, -- Can be nullable if the lead is not for a specific property yet
    source_id INT NOT NULL,
    status_id INT NOT NULL,
    assigned_to INT NOT NULL, -- The user_id of the salesperson
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_contact
        FOREIGN KEY(contact_id)
        REFERENCES contacts(contact_id),
    CONSTRAINT fk_property
        FOREIGN KEY(property_id)
        REFERENCES properties(property_id),
    CONSTRAINT fk_source
        FOREIGN KEY(source_id)
        REFERENCES lead_sources(source_id),
    CONSTRAINT fk_status
        FOREIGN KEY(status_id)
        REFERENCES lead_statuses(status_id),
    CONSTRAINT fk_assigned_to
        FOREIGN KEY(assigned_to)
        REFERENCES users(user_id)
);

-- Insert some initial lookup data
INSERT INTO lead_sources (name) VALUES ('Website Inquiry'), ('Phone Call'), ('Social Media'), ('Referral');
INSERT INTO lead_statuses (name) VALUES ('New'), ('Contacted'), ('Qualified'), ('Converted'), ('Lost');