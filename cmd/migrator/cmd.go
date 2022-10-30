package migrator

import (
	"github.com/spf13/cobra"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
)

var dbSubCommands = map[string]func() *cobra.Command{
	"up": func() *cobra.Command {
		environment := config.DEVELOPMENT

		command := &cobra.Command{
			Use:   "up",
			Short: "Migrate up the database schema",
			Run: func(cmd *cobra.Command, args []string) {
				run(environment, MigrateUp)
			},
		}

		command.PersistentFlags().Var(&environment, "env", "Environment of the application")

		return command
	},
	"down": func() *cobra.Command {
		environment := config.DEVELOPMENT

		command := &cobra.Command{
			Use:   "down",
			Short: "Migrate down the database schema",
			Run: func(cmd *cobra.Command, args []string) {
				run(environment, MigrateDown)
			},
		}

		command.PersistentFlags().Var(&environment, "env", "Environment of the application")

		return command
	},
  // TODO: add "create" command for migration
}

func CreateMigratorCmd() *cobra.Command {
	command := &cobra.Command{
		Use: "database",
	}

	for _, fn := range dbSubCommands {
		command.AddCommand(fn())
	}

	return command
}
