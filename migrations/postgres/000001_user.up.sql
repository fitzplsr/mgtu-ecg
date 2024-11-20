-- user
CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY,
    role          INTEGER     NOT NULL,
    name          TEXT        NOT NULL,
    login         TEXT UNIQUE NOT NULL,
    password_hash BYTEA,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- TODO add constraints
