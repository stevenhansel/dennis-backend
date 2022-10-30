package server

import "github.com/go-chi/chi/v5"

func (s Server) registerHttpHandler(r chi.Router) {
	episodeController := s.container.EpisodeHttpController

	r.Get("/v1/episodes/current", episodeController.GetCurrentEpisode)
	r.Get("/v1/episodes/{episodeId}", episodeController.GetEpisodeByID)
}
