package querier

import "time"

type Episode struct {
	ID          int       `json:"id"`
	Episode     int       `json:"episode"`
	EpisodeName *string   `json:"episodeName"`
	EpisodeDate time.Time `json:"episodeDate"`
	IsCurrent   bool      `json:"isCurrent"`
}

type EpisodeDetail struct {
	ID          int            `json:"id"`
	Episode     int            `json:"episode"`
	EpisodeName *string        `json:"episodeName"`
	EpisodeDate time.Time      `json:"episodeDate"`
	IsCurrent   bool           `json:"isCurrent"`
	Songs       []*EpisodeSong `json:"songs"`
}

type EpisodeSong struct {
	ID            int    `json:"id"`
	EpisodeSongID int    `json:"episodeSongId"`
	SongNameJP    string `json:"songNameJp"`
	SongNameEN    string `json:"songNameEn"`
	ArtistNameJP  string `json:"artistNameJp"`
	ArtistNameEN  string `json:"artistNameEn"`
	CoverImageURL string `json:"coverImageUrl"`
}
