package querier

type Song struct {
	ID            int    `json:"id" db:"song_id"`
	SongNameJP    string `json:"songNameJp" db:"song_name_jp"`
	SongNameEN    string `json:"songNameEn" db:"song_name_en"`
	ArtistNameJP  string `json:"artistNameJp" db:"artist_name_jp"`
	ArtistNameEN  string `json:"artistNameEn" db:"artist_name_en"`
	CoverImageURL string `json:"coverImageUrl" db:"cover_image_url"`

}
