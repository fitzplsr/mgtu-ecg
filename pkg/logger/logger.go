package logger

import "go.uber.org/zap"

type Config struct {
	Environment string `yaml:"environment" env-default:"prod" yaml-description:"available: local, dev, prod"`
	LogFilePath string `yaml:"logFilePath" env-default:"ecg.log"`
}

func Provide(cfg Config) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.FunctionKey = "method"
	return zap.Must(config.Build())
}
