package container

import (
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/songs"

	"go.uber.org/zap"
)

type Container struct {
	Log         *zap.Logger
	Config      *config.Configuration
	SongService *songs.SongService
}

func New(log *zap.Logger, config *config.Configuration, songService *songs.SongService) *Container {
	return &Container{
		Log:         log,
		Config:      config,
		SongService: songService,
	}
}
