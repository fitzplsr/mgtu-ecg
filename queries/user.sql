-- name: Create :one
INSERT INTO users (id, role, name, login, password_hash)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetByLogin :one
SELECT * FROM users WHERE login = $1;

-- name: Update :one
UPDATE users
SET
    role = COALESCE($1, role),
    name = COALESCE($2, name),
    login = COALESCE($3, login),
    password_hash = COALESCE($4, password_hash),
    updated_at = $5
WHERE id = $6
RETURNING *;
