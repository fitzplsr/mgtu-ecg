package model

import (
	"github.com/google/uuid"
	"time"
)

//easyjson:skip
type FileMeta struct {
	UserID      uuid.UUID  `json:"-"`
	Key         string     `json:"-"`
	Filename    string     `json:"filename"`
	Size        int32      `json:"size"`
	Format      FileFormat `json:"format"`
	ContentType string     `json:"content-type"`
}

//easyjson:json
type FileInfo struct {
	ID          int64     `json:"id"`
	UserID      uuid.UUID `json:"-"`
	Key         string    `json:"-"`
	Filename    string    `json:"filename"`
	Size        int32     `json:"size"`
	Format      string    `json:"format"`
	ContentType string    `json:"content-type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
