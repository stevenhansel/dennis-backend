package migrator

import (
	"fmt"

	db "github.com/stevenhansel/csm-ending-prediction-be/database"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"go.uber.org/zap"
)

type MigrationCommand string

const (
	MigrateUp   MigrationCommand = "up"
	MigrateDown MigrationCommand = "down"
)

func run(environment config.Environment, migrationCommand MigrationCommand) {
	command := MigrationCommand(migrationCommand)

	app, err := NewMigrationApp(environment)
	if err != nil {
		return
	}

	if command == MigrateUp {
		err = db.MigrateUp(app.Instance)
		if err != nil {
			app.Log.Error("Something went wrong when migrating up the database", zap.String("error", fmt.Sprint(err)))
			return
		}

		app.Log.Info("Migrated up the database successfully!")
	} else if command == MigrateDown {
		err = db.MigrateDown(app.Instance)
		if err != nil {
			app.Log.Error("Something went wrong when migrating down the database", zap.String("error", fmt.Sprint(err)))
			return
		}

		app.Log.Info("Migrated down the database successfully!")
	} else {
		app.Log.Error(fmt.Sprintf(`Migration command is invalid, should be either %s or %s`, MigrateUp, MigrateDown))
		return
	}
}
