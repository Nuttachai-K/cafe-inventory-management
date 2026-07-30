CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    user_role VARCHAR(10) NOT NULL CHECK (user_role IN ('ADMIN', 'STAFF')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

INSERT INTO users(username, email, password_hash, user_role)
VALUES (
    'admin',
    'admin@cafe.local',
    '$2a$10$1xFOSNaYaOiqfT9KXAEHped.W8TRWlr6sU6tfIdKDGglbbcm/q9zq',
    'ADMIN'
);