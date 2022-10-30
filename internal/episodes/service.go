package episodes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

type EpisodeQuerier interface {
	InsertEpisode(ctx context.Context, params *database.InsertEpisodeParams) error
	ChangeCurrentEpisode(ctx context.Context, episodeID int) error
	FindAllEpisodes(ctx context.Context) ([]*querier.Episode, error)
	FindEpisodeDetailByID(ctx context.Context, episodeID int) (*querier.EpisodeDetail, error)
	FindCurrentEpisode(ctx context.Context) (*querier.EpisodeDetail, error)
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

func (s *EpisodeService) FindAllEpisodes(ctx context.Context) ([]*querier.Episode, error) {
	return s.querier.FindAllEpisodes(ctx)
}

func (s *EpisodeService) FindEpisodeDetailByID(ctx context.Context, episodeID int) (*querier.EpisodeDetail, error) {
	episode, err := s.querier.FindEpisodeDetailByID(ctx, episodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrEpisodeNotFound(fmt.Errorf("Episode not found"))
		}
	}

	return episode, nil
}

func (s *EpisodeService) FindCurrentEpisode(ctx context.Context) (*querier.EpisodeDetail, error) {
	return s.querier.FindCurrentEpisode(ctx)
}
