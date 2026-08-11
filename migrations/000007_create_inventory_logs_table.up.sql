CREATE TABLE inventory_logs (
    id SERIAL PRIMARY KEY,
    inventory_id INT REFERENCES inventory(id) ON DELETE RESTRICT NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE RESTRICT NOT NULL,
    change_quantity INT NOT NULL,
    operation VARCHAR(6) NOT NULL CHECK (operation IN ('IN', 'OUT', 'ADJUST')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);