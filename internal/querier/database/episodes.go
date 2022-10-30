package database

import (
	"context"
	"fmt"
	"time"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
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

func (d *DatabaseQuerier) FindAllEpisodes(ctx context.Context) ([]*querier.Episode, error) {
	queryStatement := `
  select
    "e"."id" as "episode_id",
    "e"."episode" as "episode_episode",
    "e"."episode_name" as "episode_name",
    "e"."episode_date" as "episode_date",
    "e"."is_current" as "episode_is_current"
  from "episode" "e"
  order by "e"."id" asc
  `
	var episodes []*querier.Episode
	if err := d.db.Select(&episodes, queryStatement); err != nil {
		return nil, err
	}

	return episodes, nil
}
