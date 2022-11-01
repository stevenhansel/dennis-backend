package querier

type Song struct {
	ID                  int     `json:"id"`
	ReleasedAtEpisodeID *int    `json:"releasedAtEpisodeId"`
	SongNameJP          string  `json:"songNameJp"`
	SongNameEN          string  `json:"songNameEn"`
	ArtistNameJP        string  `json:"artistNameJp"`
	ArtistNameEN        string  `json:"artistNameEn"`
	CoverImageURL       string  `json:"coverImageUrl"`
	YoutubeURL          *string `json:"youtubeUrl"`
	SpotifyURL          *string `json:"spotifyUrl"`
}
