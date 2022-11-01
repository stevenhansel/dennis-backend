package server

import (
	"github.com/go-chi/chi/v5"
)

func (s *Server) registerWsHandler(r chi.Router) {
	socketState := s.container.SocketState

  r.HandleFunc("/v1/subscribe/{episodeId}", socketState.SubscribeHandler)
}

