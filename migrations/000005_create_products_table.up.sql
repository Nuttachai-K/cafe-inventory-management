CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    cafe_id INT REFERENCES cafes(id) ON DELETE RESTRICT NOT NULL,
    category_id INT REFERENCES categories(id) ON DELETE RESTRICT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();