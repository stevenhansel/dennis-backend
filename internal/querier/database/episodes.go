package database

import (
	"context"
	"database/sql"
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
	releasedSong := &querier.EpisodeSong{
		ID:                  row.SongID,
		ReleasedAtEpisodeID: row.SongReleasedAtEpisodeID,
		EpisodeSongID:       row.EpisodeSongID,
		SongNameJP:          row.SongNameJP,
		SongNameEN:          row.SongNameEN,
		ArtistNameJP:        row.SongArtistNameJP,
		ArtistNameEN:        row.SongArtistNameEN,
		CoverImageURL:       row.SongCoverImageURL,
		YoutubeURL:          row.SongYoutubeURL,
		SpotifyURL:          row.SongSpotifyURL,
	}

	return &querier.Episode{
		ID:               row.ID,
		Episode:          row.Episode,
		EpisodeName:      row.EpisodeName,
		EpisodeDate:      row.EpisodeDate,
		IsCurrent:        row.IsCurrent,
		IsVotingOpen:     row.IsVotingOpen,
		ThumbnailURL:     row.ThumbnailURL,
		NumOfVotesCasted: row.NumOfVotesCasted,
		ReleasedSong:     releasedSong,
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
			ID:                  s.SongID,
			ReleasedAtEpisodeID: s.SongReleasedAtEpisodeID,
			EpisodeSongID:       s.EpisodeSongID,
			SongNameJP:          s.SongNameJP,
			SongNameEN:          s.SongNameEN,
			ArtistNameJP:        s.SongArtistNameJP,
			ArtistNameEN:        s.SongArtistNameEN,
			CoverImageURL:       s.SongCoverImageURL,
			YoutubeURL:          s.SongYoutubeURL,
			SpotifyURL:          s.SongSpotifyURL,
		}
	}

	return &querier.EpisodeDetail{
		ID:               row[0].EpisodeID,
		Episode:          row[0].EpisodeNumber,
		EpisodeName:      row[0].EpisodeName,
		EpisodeDate:      row[0].EpisodeDate,
		IsCurrent:        row[0].EpisodeIsCurrent,
		ThumbnailURL:     row[0].EpisodeThumbnailURL,
		NumOfVotesCasted: row[0].NumOfVotesCasted,
		Songs:            songs,
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
	if err := d.db.QueryRowContext(ctx, statement, params.Episode, params.EpisodeName, params.EpisodeReleaseDate).Scan(&episodeID); err != nil {
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

	if err := d.ChangeCurrentEpisode(ctx, params.Episode); err != nil {
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

	if _, err := d.db.NamedExecContext(ctx, statement, params); err != nil {
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
	err := d.db.SelectContext(ctx, &currentEpisodeNumber, queryStatement)
	if err != nil {
		return errtrace.Wrap(err)
	}

	if len(currentEpisodeNumber) == 0 {
		_, err := d.db.ExecContext(ctx, updateStatement, episodeNumber, true)
		if err != nil {
			return errtrace.Wrap(err)
		}
	} else {
		_, err := d.db.ExecContext(ctx, updateStatement, currentEpisodeNumber[0], false)
		if err != nil {
			return errtrace.Wrap(err)
		}

		_, err = d.db.ExecContext(ctx, updateStatement, episodeNumber, true)
		if err != nil {
			return errtrace.Wrap(err)
		}

	}

	return nil
}

type EpisodeRow struct {
	ID                      int       `db:"episode_id"`
	Episode                 int       `db:"episode_number"`
	EpisodeName             *string   `db:"episode_name"`
	EpisodeDate             time.Time `db:"episode_date"`
	IsCurrent               bool      `db:"episode_is_current"`
	IsVotingOpen            bool      `db:"episode_is_voting_open"`
	ThumbnailURL            *string   `db:"episode_thumbnail_url"`
	NumOfVotesCasted        int       `db:"num_of_votes_casted"`
	SongID                  int       `db:"song_id"`
	SongReleasedAtEpisodeID *int      `db:"song_released_at_episode_id"`
	SongNameJP              string    `db:"song_name_jp"`
	SongNameEN              string    `db:"song_name_en"`
	SongArtistNameJP        string    `db:"song_artist_name_jp"`
	SongArtistNameEN        string    `db:"song_artist_name_en"`
	SongCoverImageURL       string    `db:"song_cover_image_url"`
	SongYoutubeURL          *string   `db:"song_youtube_url"`
	SongSpotifyURL          *string   `db:"song_spotify_url"`
	EpisodeSongID           int       `db:"episode_song_id"`
}

func (d *DatabaseQuerier) FindAllEpisodes(ctx context.Context) ([]*querier.Episode, error) {
	prevEpisodes, err := d.FindPreviousEpisodes(ctx)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	currentEpisode, err := d.FindCurrentEpisode(ctx)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	allEpisodes := append(prevEpisodes, &querier.Episode{
		ID:               currentEpisode.ID,
		Episode:          currentEpisode.Episode,
		EpisodeName:      currentEpisode.EpisodeName,
		EpisodeDate:      currentEpisode.EpisodeDate,
		IsCurrent:        currentEpisode.IsCurrent,
		IsVotingOpen:     currentEpisode.IsVotingOpen,
		ThumbnailURL:     currentEpisode.ThumbnailURL,
		NumOfVotesCasted: currentEpisode.NumOfVotesCasted,
	})

	return allEpisodes, nil
}

func (d *DatabaseQuerier) FindPreviousEpisodes(ctx context.Context) ([]*querier.Episode, error) {
	queryStatement := `
  select
    "e"."id" as "episode_id",
    "e"."episode" as "episode_number",
    "e"."episode_name" as "episode_name",
    "e"."episode_date" as "episode_date",
    "e"."is_current" as "episode_is_current",
    "e"."is_voting_open" as "episode_is_voting_open",
		"e"."thumbnail_url" as "episode_thumbnail_url",
		"c"."num_of_votes" as "num_of_votes_casted",
    "s"."id" as "song_id",
		"s"."released_at_episode" as "song_released_at_episode_id",
    "s"."song_name_jp" as "song_name_jp",
    "s"."song_name_en" as "song_name_en",
    "s"."artist_name_jp" as "song_artist_name_jp",
    "s"."artist_name_en" as "song_artist_name_en",
    "s"."cover_image_url" as "song_cover_image_url",
		"s"."youtube_url" as "song_youtube_url",
		"s"."spotify_url" as "song_spotify_url",
    "es"."id" as "episode_song_id"
		from "episode" "e"
		join "episode_song" "es" on "es"."episode_id" = "e"."id"
		join "song" "s" on "s"."id" = "es" ."song_id" and "s"."released_at_episode" = "e"."id"
    left join lateral (
        select count(*) as "num_of_votes"
        from "vote"
        join "episode_song" on "episode_song"."id" = "vote"."episode_song_id"
        join "song" on "song"."id" = "episode_song"."song_id"
        where 
					"episode_song"."episode_id" = "e"."id" and 
					("song"."released_at_episode" is null or "song"."released_at_episode" >= "e"."id")
    ) "c" on true
		order by "e"."id" asc
	`

	var rows []*EpisodeRow
	if err := d.db.SelectContext(ctx, &rows, queryStatement); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return toEpisodes(rows...), nil
}

type EpisodeDetailRow struct {
	EpisodeID               int       `db:"episode_id"`
	EpisodeNumber           int       `db:"episode_number"`
	EpisodeName             *string   `db:"episode_name"`
	EpisodeDate             time.Time `db:"episode_date"`
	EpisodeIsCurrent        bool      `db:"episode_is_current"`
	EpisodeIsVotingOpen     bool      `db:"episode_is_voting_open"`
	EpisodeThumbnailURL     *string   `db:"episode_thumbnail_url"`
	NumOfVotesCasted        int       `db:"num_of_votes_casted"`
	SongID                  int       `db:"song_id"`
	SongReleasedAtEpisodeID *int      `db:"song_released_at_episode_id"`
	SongNameJP              string    `db:"song_name_jp"`
	SongNameEN              string    `db:"song_name_en"`
	SongArtistNameJP        string    `db:"song_artist_name_jp"`
	SongArtistNameEN        string    `db:"song_artist_name_en"`
	SongCoverImageURL       string    `db:"song_cover_image_url"`
	SongYoutubeURL          *string   `db:"song_youtube_url"`
	SongSpotifyURL          *string   `db:"song_spotify_url"`
	EpisodeSongID           int       `db:"episode_song_id"`
}

type FindEpisodeDetailParams struct {
	EpisodeID int   `db:"episode_id"`
	IsCurrent *bool `db:"is_current"`
}

func (d *DatabaseQuerier) FindEpisodeDetailByID(ctx context.Context, episodeID int) (*querier.EpisodeDetail, error) {
	queryStatement := `
  select
    "e"."id" as "episode_id",
    "e"."episode" as "episode_number",
    "e"."episode_name" as "episode_name",
    "e"."episode_date" as "episode_date",
    "e"."is_current" as "episode_is_current",
		"e"."is_voting_open" as "episode_is_voting_open",
		"e"."thumbnail_url" as "episode_thumbnail_url",
		"c"."num_of_votes" as "num_of_votes_casted",
    "s"."id" as "song_id",
		"s"."released_at_episode" as "song_released_at_episode_id",
    "s"."song_name_jp" as "song_name_jp",
    "s"."song_name_en" as "song_name_en",
    "s"."artist_name_jp" as "song_artist_name_jp",
    "s"."artist_name_en" as "song_artist_name_en",
    "s"."cover_image_url" as "song_cover_image_url",
		"s"."youtube_url" as "song_youtube_url",
		"s"."spotify_url" as "song_spotify_url",
    "es"."id" as "episode_song_id"
  from "episode_song" "es"
  join "episode" "e" on "e"."id" = "es"."episode_id"
  join "song" "s" on "s"."id" = "es" ."song_id"
	left join lateral (
			select count(*) as "num_of_votes"
			from "vote"
			join "episode_song" on "episode_song"."id" = "vote"."episode_song_id"
			join "song" on "song"."id" = "episode_song"."song_id"
			where 
				"episode_song"."episode_id" = "e"."id" and 
				("song"."released_at_episode" is null or "song"."released_at_episode" >= "e"."id")
	) "c" on true
  where "e"."id" = $1 and ("s"."released_at_episode" is null or "s"."released_at_episode" >= "e"."id")
  `

	var row []*EpisodeDetailRow
	if err := d.db.SelectContext(ctx, &row, queryStatement, episodeID); err != nil {
		return nil, errtrace.Wrap(err)
	}

	if len(row) == 0 {
		return nil, errtrace.Wrap(sql.ErrNoRows)
	}

	return toEpisodeDetail(row), nil
}

func (d *DatabaseQuerier) FindCurrentEpisode(ctx context.Context) (*querier.EpisodeDetail, error) {
	queryStatement := `
  select
    "e"."id" as "episode_id",
    "e"."episode" as "episode_number",
    "e"."episode_name" as "episode_name",
    "e"."episode_date" as "episode_date",
    "e"."is_current" as "episode_is_current",
		"e"."is_voting_open" as "episode_is_voting_open",
		"e"."thumbnail_url" as "episode_thumbnail_url",
		"c"."num_of_votes" as "num_of_votes_casted",
    "s"."id" as "song_id",
		"s"."released_at_episode" as "song_released_at_episode_id",
    "s"."song_name_jp" as "song_name_jp",
    "s"."song_name_en" as "song_name_en",
    "s"."artist_name_jp" as "song_artist_name_jp",
    "s"."artist_name_en" as "song_artist_name_en",
    "s"."cover_image_url" as "song_cover_image_url",
		"s"."youtube_url" as "song_youtube_url",
		"s"."spotify_url" as "song_spotify_url",
    "es"."id" as "episode_song_id"
  from "episode_song" "es"
  join "episode" "e" on "e"."id" = "es"."episode_id"
  join "song" "s" on "s"."id" = "es" ."song_id"
	left join lateral (
			select count(*) as "num_of_votes"
			from "vote"
			join "episode_song" on "episode_song"."id" = "vote"."episode_song_id"
			join "song" on "song"."id" = "episode_song"."song_id"
			where 
				"episode_song"."episode_id" = "e"."id" and 
				("song"."released_at_episode" is null or "song"."released_at_episode" >= "e"."id")
	) "c" on true
  where "e"."is_current" = true and ("s"."released_at_episode" is null or "s"."released_at_episode" >= "e"."id")
  `

	var row []*EpisodeDetailRow
	if err := d.db.SelectContext(ctx, &row, queryStatement); err != nil {
		return nil, errtrace.Wrap(err)
	}

	return toEpisodeDetail(row), nil
}

func (d *DatabaseQuerier) FindEpisodeDetailByEpisodeSongID(ctx context.Context, episodeSongID int) (*querier.EpisodeDetail, error) {
	queryStatement := `
  select
    "e"."id" as "episode_id",
    "e"."episode" as "episode_number",
    "e"."episode_name" as "episode_name",
    "e"."episode_date" as "episode_date",
    "e"."is_current" as "episode_is_current",
		"e"."is_voting_open" as "episode_is_voting_open",
		"e"."thumbnail_url" as "episode_thumbnail_url",
		"c"."num_of_votes" as "num_of_votes_casted",
    "s"."id" as "song_id",
		"s"."released_at_episode" as "song_released_at_episode_id",
    "s"."song_name_jp" as "song_name_jp",
    "s"."song_name_en" as "song_name_en",
    "s"."artist_name_jp" as "song_artist_name_jp",
    "s"."artist_name_en" as "song_artist_name_en",
    "s"."cover_image_url" as "song_cover_image_url",
		"s"."youtube_url" as "song_youtube_url",
		"s"."spotify_url" as "song_spotify_url",
    "es"."id" as "episode_song_id"
  from "episode_song" "es"
  join "episode" "e" on "e"."id" = "es"."episode_id"
  join "song" "s" on "s"."id" = "es" ."song_id"
	left join lateral (
			select count(*) as "num_of_votes"
			from "vote"
			join "episode_song" on "episode_song"."id" = "vote"."episode_song_id"
			join "song" on "song"."id" = "episode_song"."song_id"
			where 
				"episode_song"."episode_id" = "e"."id" and 
				("song"."released_at_episode" is null or "song"."released_at_episode" >= "e"."id")
	) "c" on true
  where "e"."id" = (select "episode_id" from "episode_song" where "id" = $1) and ("s"."released_at_episode" is null or "s"."released_at_episode" >= "e"."id")
  `

	var row []*EpisodeDetailRow
	if err := d.db.SelectContext(ctx, &row, queryStatement, episodeSongID); err != nil {
		return nil, errtrace.Wrap(err)
	}

	if len(row) == 0 {
		return nil, errtrace.Wrap(sql.ErrNoRows)
	}

	return toEpisodeDetail(row), nil
}

type UpdateEpisodeParams struct {
	ID           int       `db:"id"`
	Episode      int       `db:"episode"`
	EpisodeName  *string   `db:"episode_name"`
	EpisodeDate  time.Time `db:"episode_date"`
	IsCurrent    bool      `db:"is_current"`
	ThumbnailURL *string   `db:"thumbnail_url"`
}

func (d *DatabaseQuerier) UpdateEpisode(ctx context.Context, params *UpdateEpisodeParams) error {
	statement := `
	update "episode"
	set
		episode = :episode,
		episode_name = :episode_name,
		episode_date = :episode_date,
		is_current = :is_current,
		thumbnail_url = :thumbnail_url
	where id = :id
	`

	if _, err := d.db.NamedExecContext(ctx, statement, params); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}
