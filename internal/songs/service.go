package songs

import (
	"context"
	"embed"
	"encoding/json"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

//go:embed data.json
var content embed.FS

type SongQuerier interface {
	CountSongs(ctx context.Context) (int, error)
	BulkInsertSong(ctx context.Context, params []*database.InsertSongParams) error
}

type SongService struct {
	querier SongQuerier
}

func NewService(querier SongQuerier) *SongService {
	return &SongService{
		querier: querier,
	}
}

func (s *SongService) InitializeSongs(ctx context.Context) error {
	count, err := s.querier.CountSongs(ctx)
	if err != nil {
		return errtrace.Wrap(err)
	}

	if count > 0 {
		return nil
	}

	var rawSongs []*database.InsertSongParams

	data, _ := content.ReadFile("data.json")
	if err := json.Unmarshal(data, &rawSongs); err != nil {
		return errtrace.Wrap(err)
	}

	if err := s.querier.BulkInsertSong(ctx, rawSongs); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}
