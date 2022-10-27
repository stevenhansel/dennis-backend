package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stevenhansel/csm-ending-prediction-be/cmd/api"
	"github.com/stevenhansel/csm-ending-prediction-be/cmd/episode"
	"github.com/stevenhansel/csm-ending-prediction-be/cmd/migrator"
)

var rootCmd = &cobra.Command{}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(api.CreateApiCmd(), migrator.CreateMigratorCmd(), episode.CreateEpisodeCmd())
}
