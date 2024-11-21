package db

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MinioParams struct {
	fx.In

	Logger *zap.Logger
	Config MinioConfig
}

type MinioConfig struct {
	Host            string `env:"MINIO_HOST"`
	Port            string `env:"MINIO_PORT"`
	AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
	UseSSL          bool   `yaml:"useSSL" env-default:"true"`
}

func NewMinioClient(p MinioParams) (*minio.Client, error) {
	endpoint := p.Config.Host + ":" + p.Config.Port
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(p.Config.AccessKeyID, p.Config.SecretAccessKey, ""),
		Secure: p.Config.UseSSL,
	})
	if err != nil {
		p.Logger.Error("connect to minio", zap.Error(err))
		return nil, err
	}
	return client, err
}
