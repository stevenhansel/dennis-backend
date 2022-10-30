package container

import (
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/episodes"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server/responseutil"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/songs"

	"go.uber.org/zap"
)

type Container struct {
	Log                   *zap.Logger
	Environment           config.Environment
	Config                *config.Configuration
	Responseutil          *responseutil.Responseutil
	SongService           *songs.SongService
	EpisodeService        *episodes.EpisodeService
	EpisodeHttpController *episodes.EpisodeHttpController
}

func New(
	log *zap.Logger,
	environment config.Environment,
	config *config.Configuration,
	responseutil *responseutil.Responseutil,
	songService *songs.SongService,
	episodeService *episodes.EpisodeService,
	episodeHttpController *episodes.EpisodeHttpController,
) *Container {
	return &Container{
		Log:                   log,
		Environment:           environment,
		Config:                config,
		Responseutil:          responseutil,
		SongService:           songService,
		EpisodeService:        episodeService,
		EpisodeHttpController: episodeHttpController,
	}
}
