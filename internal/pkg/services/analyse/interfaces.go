package analyse

import (
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
	"github.com/google/uuid"
	"mime/multipart"
)

type Usecase interface {
	Upload(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (*model.FileInfo, error)
}

type Repo interface {
	SaveFileMeta(ctx context.Context, meta *model.FileMeta) (*model.FileInfo, error)
}

type FileStorage interface {
	Save(ctx context.Context, file *filestorage.File) (string, error)
}
