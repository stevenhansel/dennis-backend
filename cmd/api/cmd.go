package api

import (
	"github.com/spf13/cobra"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
)

var apiSubCommands = map[string]func() *cobra.Command{
	"run": func() *cobra.Command {
		environment := config.DEVELOPMENT

		command := &cobra.Command{
			Use:   "run",
			Short: "Run the HTTP and WebSocket Server that contains the APIs",
			Run: func(cmd *cobra.Command, args []string) {
				run(environment)
			},
		}

		command.PersistentFlags().Var(&environment, "env", "Environment of the application")

		return command
	},
}

func CreateApiCmd() *cobra.Command {
	command := &cobra.Command{
		Use: "api",
	}

	for _, fn := range apiSubCommands {
		command.AddCommand(fn())
	}

	return command
}
