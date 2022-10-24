package database

import "context"

type InsertSongParams struct {
	SongNameJP    string `json:"songNameJp" db:"song_name_jp"`
	SongNameEN    string `json:"songNameEn" db:"song_name_en"`
	ArtistNameJP  string `json:"artistNameJp" db:"artist_name_jp"`
	ArtistNameEN  string `json:"artistNameEn" db:"artist_name_en"`
	CoverImageURL string `json:"coverImageUrl" db:"cover_image_url"`
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

func (d *DatabaseQuerier) BulkInsertSong(ctx context.Context, params []*InsertSongParams) error {
	statement := `
  insert into "song" ("song_name_jp", "song_name_en", "artist_name_jp", "artist_name_en", "cover_image_url")
  values (:song_name_jp, :song_name_en, :artist_name_jp, :artist_name_en, :cover_image_url)
  `
	if _, err := d.db.NamedExec(statement, params); err != nil {
		return err
	}

	return nil
}

