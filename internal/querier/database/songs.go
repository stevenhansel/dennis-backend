package database

import (
	"context"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
)

func toSongs(rows ...*SongRow) []*querier.Song {
	results := make([]*querier.Song, len(rows))
	for i, r := range rows {
		results[i] = toSong(r)
	}

	return results
}

func toSong(row *SongRow) *querier.Song {
	return &querier.Song{
		ID:            row.ID,
		SongNameJP:    row.SongNameJP,
		SongNameEN:    row.SongNameEN,
		ArtistNameJP:  row.ArtistNameJP,
		ArtistNameEN:  row.ArtistNameEN,
		CoverImageURL: row.CoverImageURL,
	}
}

func (d *DatabaseQuerier) CountSongs(ctx context.Context) (int, error) {
	statement := `
  select count(*) as "count"
  from "song"
  `

	var count int
	if err := d.db.Get(&count, statement); err != nil {
		return 0, err
	}

	return count, nil
}

type SongRow struct {
	ID                int    `db:"id"`
	ReleasedAtEpisode *int   `db:"released_at_episode"`
	SongNameJP        string `db:"song_name_jp"`
	SongNameEN        string `db:"song_name_en"`
	ArtistNameJP      string `db:"artist_name_jp"`
	ArtistNameEN      string `db:"artist_name_en"`
	CoverImageURL     string `db:"cover_image_url"`
	YoutubeURL        string `db:"youtube_url"`
	SpotifyURL        string `db:"spotify_url"`
}

func (d *DatabaseQuerier) FindAllSongs(ctx context.Context) ([]*querier.Song, error) {
	queryStatement := `
  select
    "s"."id" as "id",
		"s"."released_at_episode" as "released_at_episode",
    "s"."song_name_jp" as "song_name_jp",
    "s"."song_name_en" as "song_name_en",
    "s"."artist_name_jp" as "artist_name_jp",
    "s"."artist_name_en" as "artist_name_en",
    "s"."cover_image_url" as "cover_image_url"
  from "song" "s"
  order by "s"."id" asc
  `
	var rows []*SongRow
	if err := d.db.SelectContext(ctx, &rows, queryStatement); err != nil {
		return nil, err
	}

	return toSongs(rows...), nil

}

type InsertSongParams struct {
	SongNameJP    string `json:"songNameJp" db:"song_name_jp"`
	SongNameEN    string `json:"songNameEn" db:"song_name_en"`
	ArtistNameJP  string `json:"artistNameJp" db:"artist_name_jp"`
	ArtistNameEN  string `json:"artistNameEn" db:"artist_name_en"`
	CoverImageURL string `json:"coverImageUrl" db:"cover_image_url"`
}

func (d *DatabaseQuerier) BulkInsertSong(ctx context.Context, params []*InsertSongParams) error {
	statement := `
  insert into "song" ("song_name_jp", "song_name_en", "artist_name_jp", "artist_name_en", "cover_image_url")
  values (:song_name_jp, :song_name_en, :artist_name_jp, :artist_name_en, :cover_image_url)
  `
	if _, err := d.db.NamedExecContext(ctx, statement, params); err != nil {
		return err
	}

	return nil
}

type UpdateSongParams struct {
	ID            int    `json:"id" db:"id"`
	SongNameJP    string `json:"songNameJp" db:"song_name_jp"`
	SongNameEN    string `json:"songNameEn" db:"song_name_en"`
	ArtistNameJP  string `json:"artistNameJp" db:"artist_name_jp"`
	ArtistNameEN  string `json:"artistNameEn" db:"artist_name_en"`
	CoverImageURL string `json:"coverImageUrl" db:"cover_image_url"`
}

func (d *DatabaseQuerier) BulkUpdateSong(ctx context.Context, params []*UpdateSongParams) error {
	for _, p := range params {
		statement := `
		update "song"
		set
			song_name_jp = :song_name_jp, 
			song_name_en = :song_name_en,
			artist_name_jp = :artist_name_jp,
			artist_name_en = :artist_name_en,
			cover_image_url = :cover_image_url
		where id = :id
		`
		if _, err := d.db.NamedExecContext(ctx, statement, p); err != nil {
			return err
		}
	}

	return nil
}
