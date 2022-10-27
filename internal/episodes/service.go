package episodes

import (
	"context"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

type EpisodeQuerier interface {
	InsertEpisode(ctx context.Context, params *database.InsertEpisodeParams) error
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
	if err := s.querier.InsertEpisode(ctx, params); err != nil {
		return err
	}

	return nil
}
