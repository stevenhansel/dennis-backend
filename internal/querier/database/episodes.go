package database

import (
	"context"
	"time"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
)

func toEpisodes(rows ...*EpisodeRow) []*querier.Episode {
	results := make([]*querier.Episode, len(rows))
	for i, r := range rows {
		results[i] = toEpisode(r)
	}

	return results
}

func toEpisode(row *EpisodeRow) *querier.Episode {
	return &querier.Episode{
		ID:          row.ID,
		Episode:     row.Episode,
		EpisodeName: row.EpisodeName,
		EpisodeDate: row.EpisodeDate,
		IsCurrent:   row.IsCurrent,
	}
}

func toEpisodeDetails(rows ...[]*EpisodeDetailRow) []*querier.EpisodeDetail {
	results := make([]*querier.EpisodeDetail, len(rows))
	for i, r := range rows {
		results[i] = toEpisodeDetail(r)
	}

	return results
}

func toEpisodeDetail(row []*EpisodeDetailRow) *querier.EpisodeDetail {
	songs := make([]*querier.EpisodeSong, len(row))
	for i, s := range row {
		songs[i] = &querier.EpisodeSong{
			ID:            s.SongID,
			SongNameJP:    s.SongNameJP,
			SongNameEN:    s.SongNameEN,
			ArtistNameJP:  s.SongArtistNameJP,
			ArtistNameEN:  s.SongArtistNameEN,
			CoverImageURL: s.SongCoverImageURL,
		}
	}

	return &querier.EpisodeDetail{
		ID:          row[0].EpisodeID,
		Episode:     row[0].EpisodeNumber,
		EpisodeName: row[0].EpisodeName,
		EpisodeDate: row[0].EpisodeDate,
		IsCurrent:   row[0].EpisodeIsCurrent,
		Songs:       songs,
	}
}

type InsertEpisodeParams struct {
	Episode            int       `db:"episode"`
	EpisodeName        *string   `db:"episode_name"`
	EpisodeReleaseDate time.Time `db:"episode_date"`
}

func (d *DatabaseQuerier) InsertEpisode(ctx context.Context, params *InsertEpisodeParams) error {
	statement := `
  insert into "episode" ("episode",  "episode_name", "episode_date")
  values ($1, $2,  $3)
  returning "id"
  `

	var episodeID int
	if err := d.db.QueryRow(statement, params.Episode, params.EpisodeName, params.EpisodeReleaseDate).Scan(&episodeID); err != nil {
		return errtrace.Wrap(err)
	}

	songs, err := d.FindAllSongs(ctx)
	if err != nil {
		return errtrace.Wrap(err)
	}

	initializeEpisodeSongParams := make([]*InitializeEpisodeSongParams, len(songs))
	for i, s := range songs {
		initializeEpisodeSongParams[i] = &InitializeEpisodeSongParams{
			EpisodeID: episodeID,
			SongID:    s.ID,
		}
	}

	if err := d.InitializeEpisodeSong(ctx, initializeEpisodeSongParams); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

type InitializeEpisodeSongParams struct {
	EpisodeID int `db:"episode_id"`
	SongID    int `db:"song_id"`
}

func (d *DatabaseQuerier) InitializeEpisodeSong(ctx context.Context, params []*InitializeEpisodeSongParams) error {
	statement := `
  insert into "episode_song"
  ("episode_id", "song_id") 
  values (:episode_id, :song_id)
  `

	if _, err := d.db.NamedExec(statement, params); err != nil {
		return err
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
		return errtrace.Wrap(err)
	}

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

type EpisodeRow struct {
	ID          int       `db:"id"`
	Episode     int       `db:"episode_number"`
	EpisodeName *string   `db:"episode_name"`
	EpisodeDate time.Time `db:"episode_date"`
	IsCurrent   bool      `db:"episode_is_current"`
}

func (d *DatabaseQuerier) FindAllEpisodes(ctx context.Context) ([]*querier.Episode, error) {
	queryStatement := `
  select
    "e"."id" as "episode_id",
    "e"."episode" as "episode_number",
    "e"."episode_name" as "episode_name",
    "e"."episode_date" as "episode_date",
    "e"."is_current" as "episode_is_current"
  from "episode" "e"
  order by "e"."id" asc
  `
	var rows []*EpisodeRow
	if err := d.db.Select(&rows, queryStatement); err != nil {
		return nil, err
	}

	return toEpisodes(rows...), nil
}

type EpisodeDetailRow struct {
	EpisodeID         int       `db:"episode_id"`
	EpisodeNumber     int       `db:"episode_number"`
	EpisodeName       *string   `db:"episode_name"`
	EpisodeDate       time.Time `db:"episode_date"`
	EpisodeIsCurrent  bool      `db:"episode_is_current"`
	SongID            int       `db:"song_id"`
	SongNameJP        string    `db:"song_name_jp"`
	SongNameEN        string    `db:"song_name_en"`
	SongArtistNameJP  string    `db:"song_artist_name_jp"`
	SongArtistNameEN  string    `db:"song_artist_name_en"`
	SongCoverImageURL string    `db:"song_cover_image_url"`
}

func (d *DatabaseQuerier) FindEpisodeDetailByID(ctx context.Context, episodeID int) (*querier.EpisodeDetail, error) {
	queryStatement := `
  select
    "e"."id" as "episode_id",
    "e"."episode" as "episde_number",
    "e"."episode_name" as "episode_name",
    "e"."episode_date" as "episode_date",
    "e"."is_current" as "episode_is_current",
    "s"."id" as "song_id",
    "s"."song_name_jp" as "song_name_jp",
    "s"."song_name_en" as "song_name_en",
    "s"."artist_name_jp" as "song_artist_name_jp",
    "s"."artist_name_en" as "song_artist_name_en",
    "s"."cover_image_url" as "song_cover_image_url"
  from "episode_song" "es"
  join "episode" "e" on "e"."id" = "es"."episode_id"
  join "song" "s" on "s"."id" = "es" ."song_id"
  where "e"."id" = $1
  `

	var row []*EpisodeDetailRow
	if err := d.db.Select(&row, queryStatement); err != nil {
		return nil, err
	}

	return toEpisodeDetail(row), nil
}
