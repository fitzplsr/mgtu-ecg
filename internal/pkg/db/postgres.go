package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func getConnStr(cfg *Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
}

type PostgresParams struct {
	fx.In

	Cfg Config
	Log *zap.Logger
}

func NewPostgresPool(p PostgresParams) (*pgxpool.Pool, error) {
	connStr := getConnStr(&p.Cfg)
	p.Log.Debug(connStr)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		p.Log.Error("parse config", zap.Error(err))
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.Cfg.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		p.Log.Error("pool connect", zap.Error(err))
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		p.Log.Error("ping pool", zap.Error(err))
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pool, nil
}

func NewPostgresConn(p PostgresParams) (*sql.DB, error) {
	connStr := getConnStr(&p.Cfg)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return db, nil
}
