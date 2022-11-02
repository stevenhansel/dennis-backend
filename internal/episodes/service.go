package episodes

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

//go:embed data.json
var content embed.FS

type EpisodeQuerier interface {
	InsertEpisode(ctx context.Context, params *database.InsertEpisodeParams) error
	ChangeCurrentEpisode(ctx context.Context, episodeID int) error
	FindAllEpisodes(ctx context.Context) ([]*querier.Episode, error)
	FindEpisodeDetailByID(ctx context.Context, episodeID int) (*querier.EpisodeDetail, error)
	FindCurrentEpisode(ctx context.Context) (*querier.EpisodeDetail, error)
	UpdateEpisode(ctx context.Context, params *database.UpdateEpisodeParams) error
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

type ThumbnailData struct {
	Episode int    `json:"episode"`
	URL     string `json:"url"`
}

func (s *EpisodeService) SynchronizeThumbnails(ctx context.Context) error {
	var thumbnailData []*ThumbnailData
	raw, _ := content.ReadFile("data.json")
	if err := json.Unmarshal(raw, &thumbnailData); err != nil {
		return errtrace.Wrap(err)
	}

	episodes, err := s.querier.FindAllEpisodes(ctx)
	if err != nil {
		return errtrace.Wrap(err)
	}

	var filteredEpisodes []*querier.Episode
	for _, d := range thumbnailData {
		for _, e := range episodes {
			if d.Episode == e.Episode {
				e.ThumbnailURL = &d.URL
				filteredEpisodes = append(filteredEpisodes, e)
			}
		}
	}

	for _, e := range filteredEpisodes {
		if err := s.querier.UpdateEpisode(ctx, &database.UpdateEpisodeParams{
			ID:           e.ID,
			Episode:      e.Episode,
			EpisodeName:  e.EpisodeName,
			EpisodeDate:  e.EpisodeDate,
			IsCurrent:    e.IsCurrent,
			ThumbnailURL: e.ThumbnailURL,
		}); err != nil {
			return errtrace.Wrap(err)
		}
	}

	return nil
}
