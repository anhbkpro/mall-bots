CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    product VARCHAR(100),
    quantity INTEGER,
    status VARCHAR(20), -- PENDING, CONFIRMED, CANCELLED
    saga_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE inventory (
    product VARCHAR(100) PRIMARY KEY,
    quantity INTEGER
);

CREATE TABLE inventory_reservations (
    id SERIAL PRIMARY KEY,
    saga_id VARCHAR(100),
    product VARCHAR(100),
    quantity INTEGER,
    status VARCHAR(20), -- RESERVED, CONFIRMED, CANCELLED
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert test data
INSERT INTO inventory (product, quantity) VALUES ('laptop', 100);
