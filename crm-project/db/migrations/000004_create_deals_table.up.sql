CREATE TABLE IF NOT EXISTS deal_stages (
    stage_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS deals (
    deal_id SERIAL PRIMARY KEY,
    lead_id INT NOT NULL,
    property_id INT NOT NULL,
    stage_id INT NOT NULL,
    deal_status VARCHAR(50) NOT NULL DEFAULT 'Pending', -- e.g., Pending, Closed-Won, Closed-Lost
    deal_amount NUMERIC(14, 2) NOT NULL,
    deal_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closing_date TIMESTAMPTZ, -- Nullable, as it's set upon closing
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_lead
        FOREIGN KEY(lead_id)
        REFERENCES leads(lead_id),
    CONSTRAINT fk_deal_property
        FOREIGN KEY(property_id)
        REFERENCES properties(property_id),
    CONSTRAINT fk_deal_stage
        FOREIGN KEY(stage_id)
        REFERENCES deal_stages(stage_id)
);

-- Insert some initial lookup data
INSERT INTO deal_stages (name) VALUES ('Prospecting'), ('Qualification'), ('Negotiation'), ('Closing');