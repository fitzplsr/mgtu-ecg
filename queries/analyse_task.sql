-- -- name: ListAnalyseTasks :many
-- SELECT *
-- FROM analyse_tasks
-- ORDER BY created_at DESC
-- LIMIT $1 OFFSET $2;

-- name: SearchAnalyseTasks :many
SELECT *
FROM analyse_tasks
WHERE name LIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAnalyseTasksByPatientID :many
SELECT *
FROM analyse_tasks
WHERE patient_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateAnalyseTask :one
INSERT INTO analyse_tasks (name, patient_id, filemeta_id, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SaveAnalyseTaskResult :one
UPDATE analyse_tasks
SET result  = $1,
    predict = $2,
    status  = $3
WHERE id = $4
RETURNING *;

-- name: SetAnalyseTaskStatus :one
UPDATE analyse_tasks
SET status  = $1
WHERE id = $2
RETURNING *;
