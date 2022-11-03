package server

import "github.com/go-chi/chi/v5"

func (s Server) registerHttpHandler(r chi.Router) {
	socketState := s.container.SocketState

	episodeController := s.container.EpisodeHttpController
	voteController := s.container.VoteHttpController

	// Episodes
	r.Get("/v1/episodes", episodeController.GetAllEpisodes)
	r.Get("/v1/episodes/current", episodeController.GetCurrentEpisode)
	r.Get("/v1/episodes/{episodeId}", episodeController.GetEpisodeByID)
	r.Get("/v1/episodes/{episodeId}/votes", voteController.GetVotesByEpisodeID)
	r.Get("/v1/episodes/{episodeId}/has_voted", voteController.HasVotedByEpisodeID)

	// Votes
	r.Post("/v1/votes", voteController.InsertVote)

	// Subscribers
	r.Get("/v1/subscribers/{episodeId}", socketState.GetNumOfSubscribers)
}
