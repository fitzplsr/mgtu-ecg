package config

import (
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/auther"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/db"
	"github.com/fitzplsr/mgtu-ecg/internal/pkg/refresh"
	"github.com/fitzplsr/mgtu-ecg/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/fx"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config/config.yaml"`

	Logger logger.Config `yaml:"logger"`

	HTTPServer HTTPServer     `yaml:"httpServer"`
	Auth       auther.Config  `yaml:"authJWT"`
	Refresh    refresh.Config `yaml:"refresh"`
	DB         db.Config      `yaml:"db"`
	Redis      db.RedisConfig `yaml:"redis"`
}

type HTTPServer struct {
	Address           string        `yaml:"address" env-default:"localhost:8080"`
	Timeout           time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout       time.Duration `yaml:"idleTimeout" env-default:"60s"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" env-defualt:"10s"`
}

type Out struct {
	fx.Out

	Logger     logger.Config
	HTTPServer HTTPServer
	Auth       auther.Config
	DB         db.Config
	Redis      db.RedisConfig
	Refresh    refresh.Config
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
		Logger:     cfg.Logger,
		HTTPServer: cfg.HTTPServer,
		Auth:       cfg.Auth,
		DB:         cfg.DB,
		Redis:      cfg.Redis,
		Refresh:    cfg.Refresh,
	}
}
