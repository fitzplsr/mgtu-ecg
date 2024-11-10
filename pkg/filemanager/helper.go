package filemanager

import (
	"errors"
	"fmt"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
	"github.com/google/uuid"
)

func GenerateFileID(file *filestorage.File) (string, error) {
	if file.UserID.String() == "" || file.Filename == "" {
		return "", errors.New("required fields are empty")
	}

	id, _ := uuid.NewV7()

	return fmt.Sprintf("%s/%s/%s", id.String(), file.UserID.String(), file.Filename), nil
}
