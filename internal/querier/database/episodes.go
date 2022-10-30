package database

import (
	"context"
	"fmt"
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

func (d *DatabaseQuerier) ChangeCurrentEpisode(ctx context.Context, episodeNumber int) error {
	updateStatement := `
  update "episode"
  set "is_current" = $2
  where "episode" = $1
  `

	var currentEpisodeNumber []int
	queryStatement := `
  select "episode"
  from "episode"
  where "is_current" = true
  `
	err := d.db.Select(&currentEpisodeNumber, queryStatement)
	if err != nil {
		fmt.Println("here")
		return errtrace.Wrap(err)
	}

	fmt.Printf("%v\n", currentEpisodeNumber)

	if len(currentEpisodeNumber) == 0 {
		_, err := d.db.Exec(updateStatement, episodeNumber, true)
		if err != nil {
			return errtrace.Wrap(err)
		}
	} else {
		_, err := d.db.Exec(updateStatement, currentEpisodeNumber[0], false)
		if err != nil {
			return errtrace.Wrap(err)
		}

		_, err = d.db.Exec(updateStatement, episodeNumber, true)
		if err != nil {
			return errtrace.Wrap(err)
		}

	}

	return nil
}
