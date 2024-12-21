package config

import (
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/analyser"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/auther"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/db"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/filestorage/fsstorage"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/refresh"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/server"
	"github.com/fitzplsr/mgtu-ecg/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
	"log"
	"os"
)

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config/config.yaml"`

	Logger logger.Config `yaml:"logger"`

	HTTPServer    server.Config    `yaml:"httpServer"`
	Auth          auther.Config    `yaml:"authJWT"`
	Refresh       refresh.Config   `yaml:"refresh"`
	DB            db.Config        `yaml:"db"`
	Redis         db.RedisConfig   `yaml:"redis"`
	Minio         db.MinioConfig   `yaml:"minio"`
	FS            fsstorage.Config `yaml:"fs"`
	AnalyseClient analyser.Config  `yaml:"analyse-client"`
}

type Out struct {
	fx.Out

	Logger        logger.Config
	HTTPServer    server.Config
	Auth          auther.Config
	DB            db.Config
	Redis         db.RedisConfig
	Refresh       refresh.Config
	Minio         db.MinioConfig
	FS            fsstorage.Config
	AnalyseClient analyser.Config
}

func MustLoad() Out {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("cannot read .env file: %s\n (fix: you need to put .env file in main dir)", err)
		os.Exit(1)
	}

	if _, err := os.Stat(cfg.ConfigPath); os.IsNotExist(err) {
		log.Printf("config file does not exist: %s", cfg.ConfigPath)
		os.Exit(1)
	}

	if err := cleanenv.ReadConfig(cfg.ConfigPath, &cfg); err != nil {
		log.Printf("cannot read %s: %v", cfg.ConfigPath, err)
		os.Exit(1)
	}

	return Out{
		Logger:        cfg.Logger,
		HTTPServer:    cfg.HTTPServer,
		Auth:          cfg.Auth,
		DB:            cfg.DB,
		Redis:         cfg.Redis,
		Refresh:       cfg.Refresh,
		Minio:         cfg.Minio,
		FS:            cfg.FS,
		AnalyseClient: cfg.AnalyseClient,
	}
}
