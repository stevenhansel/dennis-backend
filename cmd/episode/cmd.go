package episode

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

var subcommands = map[string]func() *cobra.Command{
	"create": func() *cobra.Command {
		environment := config.DEVELOPMENT
		controller, err := initializeController(environment)
		if err != nil {
			fmt.Printf("Something when wrong when executing the command: %v", err)
			os.Exit(1)
		}

		command := &cobra.Command{
			Use:   "create",
			Short: "Create a new episode",
			Run: func(cmd *cobra.Command, args []string) {
				reader := bufio.NewReader(os.Stdin)

				fmt.Printf("Episode Number: ")
				episodeNumberStr, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Invalid input")
					os.Exit(1)
				}

				episodeNumberStr = strings.Trim(episodeNumberStr, "\n")
				episodeNumber, err := strconv.Atoi(episodeNumberStr)
				if err != nil {
					fmt.Println("Episode number should be a valid integer")
					os.Exit(1)
				}

				var episodeName *string
				fmt.Printf("Episode Name (optional): ")
				episodeNameStr, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Invalid input")
					os.Exit(1)
				}

				episodeNameStr = strings.Trim(episodeNameStr, "\n")
				if episodeNameStr != "" {
					episodeName = &episodeNameStr
				}

				fmt.Printf("Episode Release Date: ")
				episodeReleaseDateStr, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Invalid input")
					os.Exit(1)
				}

				episodeReleaseDateStr = strings.Trim(episodeReleaseDateStr, "\n")
				episodeReleaseDate, err := time.Parse(time.RFC3339, episodeReleaseDateStr)
				if err != nil {
					fmt.Println("Date should be formatted in RFC3339")
					os.Exit(1)
				}

				if err := controller.createEpisode(&database.InsertEpisodeParams{
					Episode:            episodeNumber,
					EpisodeName:        episodeName,
					EpisodeReleaseDate: episodeReleaseDate,
				}); err != nil {
					fmt.Println("Something when wrong when creating the episode")
					os.Exit(1)
				}

				fmt.Printf("Episode #%d created successfully", episodeNumber)
			},
		}

		command.PersistentFlags().Var(&environment, "env", "Environment of the application")

		return command
	},
	"change": func() *cobra.Command {
		environment := config.DEVELOPMENT
		controller, err := initializeController(environment)
		if err != nil {
			fmt.Printf("Something when wrong when executing the command: %v", err)
			os.Exit(1)
		}

		command := &cobra.Command{
			Use:   "change",
			Short: "Change the current active episode",
			Run: func(cmd *cobra.Command, args []string) {
				reader := bufio.NewReader(os.Stdin)

				fmt.Printf("Episode Number: ")
				episodeNumberStr, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Invalid input")
					os.Exit(1)
				}

				episodeNumberStr = strings.Trim(episodeNumberStr, "\n")
				episodeNumber, err := strconv.Atoi(episodeNumberStr)
				if err != nil {
					fmt.Println("Episode number should be a valid integer")
					os.Exit(1)
				}

				if err := controller.changeCurrentEpisode(episodeNumber); err != nil {
					fmt.Printf("Something when wrong when changing the current episode: %v", err)
					os.Exit(1)
				}

				fmt.Printf("Successfully changed the current episode to %d\n", episodeNumber)
			},
		}

		return command
	},
}

func CreateEpisodeCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "episode",
		Short: "View list of all episodes available",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	for _, fn := range subcommands {
		command.AddCommand(fn())
	}

	return command

}
