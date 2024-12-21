package analyse

import "errors"

var (
	ErrFileNotExist    = errors.New("file not exist")
	ErrPatientNotExist = errors.New("patient not exist")
)
