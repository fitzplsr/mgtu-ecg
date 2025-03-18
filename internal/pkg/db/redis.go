package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

type RedisParams struct {
	fx.In

	Config RedisConfig
}

// RedisConfig represents the configuration for a Redis client
type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

// NewRedisClient creates a new Redis client
func NewRedisClient(p RedisParams) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", p.Config.Host, p.Config.Port),
		Password: p.Config.Password,
		DB:       p.Config.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, err
}
