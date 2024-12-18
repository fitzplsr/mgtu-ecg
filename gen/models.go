// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Filemeta struct {
	ID          int32
	Format      int16
	Size        int32
	Filename    string
	ContentType string
	Key         string
	UserID      pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type User struct {
	ID           pgtype.UUID
	Role         int32
	Name         string
	Login        string
	PasswordHash []byte
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}
