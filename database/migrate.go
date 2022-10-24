package database

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations
var migrationFs embed.FS

func NewMigrationInstance(config *config.Configuration) (*migrate.Migrate, error) {

	db, err := sql.Open("postgres", config.POSTGRES_CONNECTION_URI)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	src, err := iofs.New(migrationFs, "migrations")
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"migrations",
		src,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func MigrateUp(m *migrate.Migrate) error {
	err := m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func MigrateDown(m *migrate.Migrate) error {
	err := m.Down()
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
