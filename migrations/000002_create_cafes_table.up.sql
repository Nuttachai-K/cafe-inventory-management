CREATE TABLE cafes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    latitude DECIMAL(9,6) NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    nearest_station VARCHAR(255),
    opening_time TIME,
    closing_time TIME,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_cafes_updated_at
BEFORE UPDATE ON cafes
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();