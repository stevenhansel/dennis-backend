package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	rootCmd.AddCommand(apiCmd, databaseCmd, createEpisodeCmd())
}
