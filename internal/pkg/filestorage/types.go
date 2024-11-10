package filestorage

import (
	"bytes"
	"github.com/google/uuid"
)

type File struct {
	Data        bytes.Buffer
	Filename    string
	Size        int64
	UserID      uuid.UUID
	ContentType string
}
