package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed postgres/*.sql
var migrationFiles embed.FS

type Params struct {
	fx.In

	Log *zap.Logger
	DB  *sql.DB
}

func RunMigrations(p Params) error {
	// Initialize the source driver with embedded migration files
	sourceDriver, err := iofs.New(migrationFiles, "postgres")
	if err != nil {
		return fmt.Errorf("failed to initialize migrations source driver: %w", err)
	}

	// Initialize the database driver for postgres
	dbDriver, err := postgres.WithInstance(p.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to initialize postgres driver: %w", err)
	}

	// Create the migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	err = sourceDriver.Close()
	if err != nil {
		p.Log.Error("Failed to close sourceDriver", zap.Error(err))
	}

	err = dbDriver.Close()
	if err != nil {
		p.Log.Error("Failed to close dbDriver", zap.Error(err))
	}

	return nil
}
