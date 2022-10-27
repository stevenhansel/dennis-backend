package episode

import "github.com/spf13/cobra"

func CreateEpisodeCmd() *cobra.Command {
  return &cobra.Command{
    Use: "episode",
  }
}

