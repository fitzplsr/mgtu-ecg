package analyseusecase

import (
	"bytes"
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/model"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"mime/multipart"
)

var _ analyse.Usecase = (*Analyse)(nil)

type Params struct {
	fx.In

	Log         *zap.Logger
	FileStorage analyse.FileStorage
	Repo        analyse.Repo
}

type Analyse struct {
	log *zap.Logger

	repo        analyse.Repo
	fileStorage analyse.FileStorage
}

func New(p Params) *Analyse {
	return &Analyse{
		log:         p.Log,
		repo:        p.Repo,
		fileStorage: p.FileStorage,
	}
}

func (a *Analyse) Upload(ctx context.Context, fileHeader *multipart.FileHeader, userID uuid.UUID) (*model.FileInfo, error) {
	open, err := fileHeader.Open()
	if err != nil {
		a.log.Error("open file", zap.Error(err))
		return nil, err
	}
	defer open.Close()

	var data []byte
	_, err = open.Read(data)
	if err != nil {
		a.log.Error("read file", zap.Error(err))
		return nil, err
	}

	var buffer bytes.Buffer
	_, err = buffer.Read(data)
	if err != nil {
		a.log.Error("write to buffer", zap.Error(err))
		return nil, err
	}

	file := filestorage.File{
		Data:        buffer,
		Size:        fileHeader.Size,
		Filename:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		UserID:      userID,
	}

	key, err := a.fileStorage.Save(ctx, &file)
	if err != nil {
		a.log.Error("save file", zap.Error(err))
		return nil, err
	}

	meta := model.FileMeta{
		UserID:      userID,
		Key:         key,
		Filename:    fileHeader.Filename,
		Size:        int32(fileHeader.Size),
		Format:      model.EDF,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}

	savedMeta, err := a.repo.SaveFileMeta(ctx, &meta)
	if err != nil {
		a.log.Error("save file meta", zap.Error(err))
		return nil, err
	}

	return savedMeta, nil
}
