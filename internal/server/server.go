package server

import (
	"net/http"
	"time"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/container"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server/middleware"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	container *container.Container

	Server *http.Server
}

func New(container *container.Container) *Server {
	s := &Server{
		container: container,
	}

	s.Server = &http.Server{
		Handler:      s.registerHandler(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	return s
}

func (s *Server) registerHandler() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Cors(s.container.Environment))
	s.registerHttpHandler(r)
	s.registerWsHandler(r)

	return r
}
