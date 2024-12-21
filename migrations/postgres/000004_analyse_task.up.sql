-- analyse_results
CREATE TABLE IF NOT EXISTS analyse_tasks
(
    id          SERIAL PRIMARY KEY,
    name        TEXT     not null        default 'Новое исследование',
    patient_id  INTEGER  REFERENCES patients (id) ON DELETE SET NULL,
    filemeta_id INTEGER  REFERENCES filemetas (id) ON DELETE SET NULL,
    status      SMALLINT NOT NULL        DEFAULT 0,
    result      SMALLINT NOT NULL        DEFAULT 0,
    predict     TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER modify_analyse_tasks_updated_at
    BEFORE UPDATE
    ON analyse_tasks
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);
