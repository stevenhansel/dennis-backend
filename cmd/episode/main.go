package episode

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/episodes"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

type EpisodeController struct {
	service *episodes.EpisodeService
}

func initializeController(environment config.Environment) (*EpisodeController, error) {
	config, err := config.New(environment)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	db, err := sqlx.Connect("postgres", config.POSTGRES_CONNECTION_URI)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	service := episodes.NewService(database.New(db))

	return &EpisodeController{
		service: service,
	}, nil
}

func (c *EpisodeController) createEpisode(params *database.InsertEpisodeParams) error {
	return c.service.CreateEpisode(context.Background(), params)
}
