-- name: CreateFileMeta :one
INSERT INTO filemetas (format,
                       size,
                       filename,
                       content_type,
                       key,
                       user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFileMetaById :one
SELECT *
FROM filemetas
WHERE id = $1
LIMIT 1;

-- name: GetUserFileMetas :many
SELECT *
FROM filemetas
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: DeleteFileMeta :exec
DELETE
FROM filemetas
WHERE id = $1;
