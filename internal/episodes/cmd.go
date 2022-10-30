package episodes

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

type EpisodeCmdController struct {
	service *EpisodeService
}

func NewCmdController(environment config.Environment) (*EpisodeCmdController, error) {
	config, err := config.New(environment)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	db, err := sqlx.Connect("postgres", config.POSTGRES_CONNECTION_URI)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	service := NewService(database.New(db))

	return &EpisodeCmdController{
		service: service,
	}, nil
}

func (c *EpisodeCmdController) CreateEpisode(params *database.InsertEpisodeParams) error {
	return c.service.CreateEpisode(context.Background(), params)
}

func (c *EpisodeCmdController) ChangeCurrentEpisode(episodeNumber int) error {
	return c.service.ChangeCurrentEpisode(context.Background(), episodeNumber)
}

func (c *EpisodeCmdController) FindAllEpisodes() ([]*querier.Episode, error) {
	return c.service.FindAllEpisodes(context.Background())
}
