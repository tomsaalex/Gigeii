-- +goose Up
SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    email text NOT NULL UNIQUE,
    username text NOT NULL UNIQUE,
    pass_hash bytea NOT NULL,
    pass_salt bytea NOT NULL,
    role text NOT NULL DEFAULT 'user'
    );

-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS users;

