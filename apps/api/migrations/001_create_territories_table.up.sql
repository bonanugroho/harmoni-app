CREATE TABLE IF NOT EXISTS territories (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('rt', 'rw')),
    parent_id VARCHAR(50) REFERENCES territories(id),
    created_at TIMESTAMP DEFAULT NOW()
);
