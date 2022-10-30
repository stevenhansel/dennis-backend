package querier

import "time"

type Episode struct {
	ID          int       `json:"id" db:"episode_id"`
	Episode     int       `json:"episode" db:"episode_episode"`
	EpisodeName *string   `json:"episodeName" db:"episode_name"`
	EpisodeDate time.Time `json:"episodeDate" db:"episode_date"`
	IsCurrent   bool      `json:"isCurrent" db:"episode_is_current"`
}
