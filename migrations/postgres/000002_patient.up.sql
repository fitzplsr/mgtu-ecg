-- patient
CREATE TABLE IF NOT EXISTS patients
(
    id         serial PRIMARY KEY,
    name       TEXT NOT NULL,
    surname    TEXT,
    bdate      DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- TODO add constraints

CREATE TRIGGER modify_patients_updated_at
    BEFORE UPDATE
    ON patients
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

INSERT INTO patients (name) VALUES ('default patient')