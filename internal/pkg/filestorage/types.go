package filestorage

import (
	"bytes"
)

type File struct {
	Data        *bytes.Buffer
	Filename    string
	Size        int64
	PatientID   int
	ContentType string
}
