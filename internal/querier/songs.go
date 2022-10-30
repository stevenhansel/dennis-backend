package querier

type Song struct {
	ID            int    `json:"id"`
	SongNameJP    string `json:"songNameJp"`
	SongNameEN    string `json:"songNameEn"`
	ArtistNameJP  string `json:"artistNameJp"`
	ArtistNameEN  string `json:"artistNameEn"`
	CoverImageURL string `json:"coverImageUrl"`
}
