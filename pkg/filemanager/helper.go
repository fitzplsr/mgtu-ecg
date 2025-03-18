package filemanager

import (
	"errors"
	"fmt"

	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
	"github.com/google/uuid"
)

func GenerateFileID(file *filestorage.File) (string, error) {
	if file.PatientID == 0 || file.Filename == "" {
		return "", errors.New("required fields are empty")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("generate uuid: %w", err)
	}

	return fmt.Sprintf("%s-%d-%s", id.String(), file.PatientID, file.Filename), nil
}
