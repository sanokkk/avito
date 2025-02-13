CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  
    username text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    salt BYTEA NOT NULL
);

CREATE INDEX idx_users_username ON users (username);