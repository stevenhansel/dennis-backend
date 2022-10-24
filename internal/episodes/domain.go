package episodes

type Episode struct {
	ID          int
	Episode     int
	EpisodeName string
}

type EpisodeSong struct {
	ID            int
	Rank          int
	NumberOfVotes string
	SongID        int
	EpisodeID     int
}

type Song struct {
	ID           int
	SongJpName   string
	SongEnName   string
	ArtistJpName string
	ArtistEnName string
	ImageURL     string
}

type Vote struct {
	ID            int
	IPAddress     string
	EpisodeSongID int
}
