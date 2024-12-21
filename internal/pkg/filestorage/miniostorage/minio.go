package miniostorage

import (
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/fitzplsr/mgtu-ecg/pkg/filemanager"
	"github.com/minio/minio-go/v7"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var _ analyse.FileStorage = (*MinioStorage)(nil)

type Params struct {
	fx.In

	Client *minio.Client
	Logger *zap.Logger
}

type MinioStorage struct {
	log    *zap.Logger
	client *minio.Client
}

func New(p Params) *MinioStorage {
	return &MinioStorage{
		log:    p.Logger,
		client: p.Client,
	}
}

func (m *MinioStorage) Save(ctx context.Context, file *filestorage.File) (string, error) {
	key, err := filemanager.GenerateFileID(file)
	if err != nil {
		m.log.Error("generate file id", zap.Error(err))
		return "", err
	}

	_, err = m.client.PutObject(
		ctx,
		string(filestorage.EDFBucket),
		key,
		file.Data,
		file.Size,
		minio.PutObjectOptions{ContentType: file.ContentType},
	)

	return key, nil
}
