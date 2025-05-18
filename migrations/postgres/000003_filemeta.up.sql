-- filemeta
CREATE TABLE IF NOT EXISTS filemetas
(
    id           SERIAL PRIMARY KEY,
    format       SMALLINT    NOT NULL,
    size         INTEGER     NOT NULL,
    filename     TEXT        NOT NULL,
    content_type TEXT        NOT NULL,
    key          TEXT UNIQUE NOT NULL,
    data         JSONB,
    patient_id   INTEGER     REFERENCES patients (id) ON DELETE SET NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER modify_filemetas_updated_at
    BEFORE UPDATE
    ON filemetas
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);