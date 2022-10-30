package episode

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/episodes"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

var subcommands = map[string]func() *cobra.Command{
	"create": func() *cobra.Command {
		environment := config.DEVELOPMENT
		controller, err := episodes.NewCmdController(environment)
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

				if err := controller.CreateEpisode(&database.InsertEpisodeParams{
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
		controller, err := episodes.NewCmdController(environment)
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

				if err := controller.ChangeCurrentEpisode(episodeNumber); err != nil {
					fmt.Printf("Something when wrong when changing the current episode: %v", err)
					os.Exit(1)
				}

				fmt.Printf("Successfully changed the current episode to %d\n", episodeNumber)
			},
		}

		return command
	},
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type listEpisodeModel struct {
	list list.Model
}

type listEpisodeItem struct {
	title, desc string
}

func (i listEpisodeItem) Title() string       { return i.title }
func (i listEpisodeItem) Description() string { return i.desc }
func (i listEpisodeItem) FilterValue() string { return i.title }

func newListEpisodeModel(episodes []*querier.Episode) listEpisodeModel {
	items := make([]list.Item, len(episodes))
	for i, e := range episodes {
		var desc string
		if e.EpisodeName != nil {
			desc += fmt.Sprintf("%s ", *e.EpisodeName)
		}

    desc += e.EpisodeDate.Format(time.RFC1123)

		items[i] = listEpisodeItem{
			title: fmt.Sprintf("Episode #%d", e.Episode),
			desc:  desc,
		}
	}

	m := listEpisodeModel{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Episode List"

	return m
}

func (m listEpisodeModel) Init() tea.Cmd {
	return nil
}

func (m listEpisodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listEpisodeModel) View() string {
	return docStyle.Render(m.list.View())
}

func CreateEpisodeCmd() *cobra.Command {
	environment := config.DEVELOPMENT

	command := &cobra.Command{
		Use:   "episode",
		Short: "View list of all episodes available",
		Run: func(cmd *cobra.Command, args []string) {
			controller, err := episodes.NewCmdController(environment)
			if err != nil {
				fmt.Printf("Something when wrong when executing the command: %v\n", err)
				os.Exit(1)
			}

			episodes, err := controller.FindAllEpisodes()
			if err != nil {
				fmt.Printf("Something went wrong when getting the episodes: %v\n", err)
				os.Exit(1)
			}

			p := tea.NewProgram(newListEpisodeModel(episodes), tea.WithAltScreen())
			if err := p.Start(); err != nil {
				fmt.Printf("Something when wrong when executing the command: %v\n", err)
				os.Exit(1)
			}
		},
	}

	command.PersistentFlags().Var(&environment, "env", "Environment of the application")

	for _, fn := range subcommands {
		command.AddCommand(fn())
	}

	return command

}
