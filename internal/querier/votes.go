package querier

import "time"

type Vote struct {
	ID            int       `json:"id"`
	IPAddress     string    `json:"ipAddress"`
	EpisodeSongID int       `json:"episodeSongId"`
	CreatedAt     time.Time `json:"createdAt"`
}
