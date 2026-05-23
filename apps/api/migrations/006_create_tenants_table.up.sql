CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    block VARCHAR(10) NOT NULL,
    unit_number VARCHAR(10) NOT NULL,
    occupancy VARCHAR(20) NOT NULL CHECK (occupancy IN ('occupied', 'vacant')),
    monthly_fee NUMERIC(12,2) NOT NULL,
    territory_id VARCHAR(50) NOT NULL REFERENCES territories(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(block, unit_number, territory_id)
);

CREATE INDEX idx_tenants_territory ON tenants(territory_id);
