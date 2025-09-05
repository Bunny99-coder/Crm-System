CREATE TABLE IF NOT EXISTS sites (
    site_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    number_of_buildings INT
);

CREATE TABLE IF NOT EXISTS property_types (
    property_type_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS properties (
    property_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    site_id INT NOT NULL,
    property_type_id INT NOT NULL,
    unit_no VARCHAR(50),
    size_sqft NUMERIC(10, 2), -- e.g., 1250.50 square feet
    price NUMERIC(12, 2) NOT NULL, -- e.g., 500000.00
    status VARCHAR(50) NOT NULL DEFAULT 'Available', -- e.g., Available, Sold
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_site
        FOREIGN KEY(site_id)
        REFERENCES sites(site_id),
    CONSTRAINT fk_property_type
        FOREIGN KEY(property_type_id)
        REFERENCES property_types(property_type_id)
);

-- Insert some initial property types
INSERT INTO property_types (name) VALUES ('Apartment'), ('Villa'), ('Office'), ('Townhouse');