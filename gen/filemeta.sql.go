// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: filemeta.sql

package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFileMeta = `-- name: CreateFileMeta :one
INSERT INTO filemetas (format,
                       size,
                       filename,
                       content_type,
                       key,
                       user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, format, size, filename, content_type, key, user_id, created_at, updated_at
`

type CreateFileMetaParams struct {
	Format      int16
	Size        int32
	Filename    string
	ContentType string
	Key         string
	UserID      pgtype.UUID
}

func (q *Queries) CreateFileMeta(ctx context.Context, arg CreateFileMetaParams) (Filemeta, error) {
	row := q.db.QueryRow(ctx, createFileMeta,
		arg.Format,
		arg.Size,
		arg.Filename,
		arg.ContentType,
		arg.Key,
		arg.UserID,
	)
	var i Filemeta
	err := row.Scan(
		&i.ID,
		&i.Format,
		&i.Size,
		&i.Filename,
		&i.ContentType,
		&i.Key,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFileMeta = `-- name: DeleteFileMeta :exec
DELETE
FROM filemetas
WHERE id = $1
`

func (q *Queries) DeleteFileMeta(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteFileMeta, id)
	return err
}

const getFileMetaById = `-- name: GetFileMetaById :one
SELECT id, format, size, filename, content_type, key, user_id, created_at, updated_at
FROM filemetas
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetFileMetaById(ctx context.Context, id int32) (Filemeta, error) {
	row := q.db.QueryRow(ctx, getFileMetaById, id)
	var i Filemeta
	err := row.Scan(
		&i.ID,
		&i.Format,
		&i.Size,
		&i.Filename,
		&i.ContentType,
		&i.Key,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserFileMetas = `-- name: GetUserFileMetas :many
SELECT id, format, size, filename, content_type, key, user_id, created_at, updated_at
FROM filemetas
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3
`

type GetUserFileMetasParams struct {
	UserID pgtype.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) GetUserFileMetas(ctx context.Context, arg GetUserFileMetasParams) ([]Filemeta, error) {
	rows, err := q.db.Query(ctx, getUserFileMetas, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Filemeta
	for rows.Next() {
		var i Filemeta
		if err := rows.Scan(
			&i.ID,
			&i.Format,
			&i.Size,
			&i.Filename,
			&i.ContentType,
			&i.Key,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}