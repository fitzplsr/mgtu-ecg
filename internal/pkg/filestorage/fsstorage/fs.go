package fsstorage

import (
	"context"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/fitzplsr/mgtu-ecg/pkg/filemanager"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

var _ analyse.FileStorage = &FSStorage{}

type Config struct {
	Path string `yaml:"files_path" env-default:"files/inputs"`
}

type Params struct {
	fx.In

	Logger *zap.Logger
	Config Config
}

type FSStorage struct {
	log  *zap.Logger
	path string
}

func New(p Params) *FSStorage {
	return &FSStorage{
		log:  p.Logger,
		path: p.Config.Path,
	}
}

func (m *FSStorage) Save(_ context.Context, file *filestorage.File) (string, error) {
	key, err := filemanager.GenerateFileID(file)
	if err != nil {
		m.log.Error("generate file id", zap.Error(err))
		return "", err
	}

	f, err := os.Create(filepath.Join(m.path, key))
	if err != nil {
		return "", err
	}

	_, err = f.ReadFrom(file.Data)
	if err != nil {
		return "", err
	}

	return key, nil
}
