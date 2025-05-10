package fsstorage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/services/analyse"
	"github.com/fitzplsr/mgtu-ecg/pkg/filemanager"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

	fullPath := filepath.Join(m.path, key)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		m.log.Error("failed to create directories", zap.Error(err))
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	defer f.Close()

	_, err = f.ReadFrom(file.Data)
	if err != nil {
		return "", fmt.Errorf("failed to ReadFrom: %w", err)
	}

	return key, nil
}
