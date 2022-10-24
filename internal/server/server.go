package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/container"
	"golang.org/x/time/rate"
)

type Server struct {
	container *container.Container

	HTTPServer     *http.Server
	InternalServer *InternalServer
}

type InternalServer struct {
	subscriberMessageBuffer int
	publishLimiter          *rate.Limiter
	serveMux                http.ServeMux

	subscribersMutex sync.Mutex
	subscribers      map[*subscriber]struct{}
}

func New(container *container.Container) *Server {
	s := &Server{
		container: container,
		InternalServer: &InternalServer{
			subscriberMessageBuffer: 16,
			subscribers:             make(map[*subscriber]struct{}),
			publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		},
	}

	s.HTTPServer = &http.Server{
		Handler:      s.registerHandler(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	return s
}

func (s *Server) registerHandler() chi.Router {
	r := chi.NewRouter()

	s.registerHttpHandler(r)
	s.registerWsHandler(r)

	return r
}

func (s *Server) registerHttpHandler(r chi.Router) {}

func (s *Server) registerWsHandler(r chi.Router) {}
