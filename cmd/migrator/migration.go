package migrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	"github.com/stevenhansel/csm-ending-prediction-be/database"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
)

type MigrationApplication struct {
	Instance *migrate.Migrate
	Log      *zap.Logger
	Config   *config.Configuration
}

func NewMigrationApp(environment config.Environment) (*MigrationApplication, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	config, err := config.New(environment)
	if err != nil {
		log.Fatal("Something went wrong when initializing the configuration", zap.String("error", fmt.Sprint(err)))
		return nil, err
	}

	m, err := database.NewMigrationInstance(config)
	if err != nil {
		log.Fatal("Something went wrong when initializing the migration instance", zap.String("error", fmt.Sprint(err)))
		return nil, err
	}

	return &MigrationApplication{
		Instance: m,
		Log:      log,
		Config:   config,
	}, nil
}
