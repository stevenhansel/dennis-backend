package database

import (
	"context"
	"time"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
)

type InsertEpisodeParams struct {
	Episode            int       `db:"episode"`
	EpisodeName        *string   `db:"episode_name"`
	EpisodeReleaseDate time.Time `db:"episode_date"`
}

func (d *DatabaseQuerier) InsertEpisode(ctx context.Context, params *InsertEpisodeParams) error {
	statement := `
  insert into "episode" ("episode",  "episode_name", "episode_date")
  values (:episode, :episode_name,  :episode_date)
  `

	if _, err := d.db.NamedExec(statement, params); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}
