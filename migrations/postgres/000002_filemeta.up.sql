-- filemeta
CREATE TABLE IF NOT EXISTS filemetas
(
    id           SERIAL PRIMARY KEY,
    format       SMALLINT     NOT NULL,
    size         INTEGER     NOT NULL,
    filename     TEXT        NOT NULL,
    content_type TEXT        NOT NULL,
    key          TEXT UNIQUE NOT NULL,
    user_id      UUID        REFERENCES users (id) ON DELETE SET NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);