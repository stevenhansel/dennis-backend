package querier

import "time"

type Vote struct {
	ID            int       `json:"id"`
	IPAddress     string    `json:"-"`
	EpisodeSongID int       `json:"episodeSongId"`
	CreatedAt     time.Time `json:"createdAt"`
}

type EpisodeVote struct {
	EpisodeSongID int `json:"episodeSongId"`
	NumOfVotes    int `json:"numOfVotes"`
	Rank          int `json:"rank,omitempty"`
}
