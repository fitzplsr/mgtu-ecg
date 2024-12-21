-- name: CreatePatient :one
INSERT INTO patients (name, surname, bdate)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPatientByID :one
SELECT *
FROM patients
WHERE id = $1;

-- name: ListPatients :many
SELECT *
FROM patients
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchPatient :many
SELECT *
FROM patients
WHERE name like '%' || $1 || '%'
   or surname like '%' || $1 || '%'
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;
