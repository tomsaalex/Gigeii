-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    timezone TEXT NOT NULL,  -- IANA timezone (e.g., 'Europe/Bucharest')
    created_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS products;
