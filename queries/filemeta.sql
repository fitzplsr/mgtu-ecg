-- name: CreateFileMeta :one
INSERT INTO filemetas (format,
                       size,
                       filename,
                       content_type,
                       key,
                       patient_id,
                       data)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFileMetaById :one
SELECT *
FROM filemetas
WHERE id = $1
LIMIT 1;

-- name: GetPatientFileMetas :many
SELECT *
FROM filemetas
WHERE patient_id = $1
ORDER BY created_at
LIMIT $2 OFFSET $3;

-- name: DeleteFileMeta :exec
DELETE
FROM filemetas
WHERE id = $1;
