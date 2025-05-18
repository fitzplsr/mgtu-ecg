package model

import (
	"time"
)

//easyjson:skip
type FileMeta struct {
	PatientID   int        `json:"patient_id"`
	Key         string     `json:"-"`
	Filename    string     `json:"filename"`
	Size        int32      `json:"size"`
	Format      FileFormat `json:"format"`
	ContentType string     `json:"content-type"`
	Data        []byte     `json:"data"`
}

//easyjson:json
type FileInfos []*FileInfo

//easyjson:json
type FileInfo struct {
	ID          int64     `json:"id"`
	PatientID   int       `json:"patient_id"`
	Key         string    `json:"-"`
	Filename    string    `json:"filename"`
	Size        int32     `json:"size"`
	Format      string    `json:"format"`
	ContentType string    `json:"content-type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Data        string    `json:"data"`
}

//easyjson:json
type PatientFiles struct {
	PatientID int         `json:"-"`
	Files     []*FileInfo `json:"files"`
}
