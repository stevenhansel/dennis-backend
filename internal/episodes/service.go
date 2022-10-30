package episodes

import (
	"context"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

type EpisodeQuerier interface {
	InsertEpisode(ctx context.Context, params *database.InsertEpisodeParams) error
	ChangeCurrentEpisode(ctx context.Context, episodeID int) error
}

type EpisodeService struct {
	querier EpisodeQuerier
}

func NewService(querier EpisodeQuerier) *EpisodeService {
	return &EpisodeService{
		querier: querier,
	}
}

func (s *EpisodeService) CreateEpisode(ctx context.Context, params *database.InsertEpisodeParams) error {
	return s.querier.InsertEpisode(ctx, params)
}

func (s *EpisodeService) ChangeCurrentEpisode(ctx context.Context, episodeNumber int) error {
	return s.querier.ChangeCurrentEpisode(ctx, episodeNumber)
}
