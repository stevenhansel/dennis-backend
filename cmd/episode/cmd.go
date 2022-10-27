package episode

import (
	"github.com/spf13/cobra"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
)

var subcommands = map[string]func() *cobra.Command{
	"create": func() *cobra.Command {
		environment := config.DEVELOPMENT

		// var episode int
		// var episodeName *string
		// var episodeDate string

		command := &cobra.Command{
			Use:   "create",
			Short: "Create a new episode",
			Run: func(cmd *cobra.Command, args []string) {
				// controller, err := initializeController(environment)
				// if err != nil {
				// 	os.Exit(1)
				// }

				// controller.createEpisode()
			},
		}

		command.PersistentFlags().Var(&environment, "env", "Environment of the application")

		return command
	},
}

func CreateEpisodeCmd() *cobra.Command {
	command := &cobra.Command{
		Use: "api",
	}

	for _, fn := range subcommands {
		command.AddCommand(fn())
	}

	return command

}
