package episodes

type EpisodeHttpController struct{}

func NewEpisodeHttpController() *EpisodeHttpController {
	return &EpisodeHttpController{}
}

func (c *EpisodeHttpController) GetCurrentEpisode() {}
