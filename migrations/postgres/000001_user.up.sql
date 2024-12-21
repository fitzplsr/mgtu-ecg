CREATE EXTENSION IF NOT EXISTS moddatetime;

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

CREATE TRIGGER modify_users_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
