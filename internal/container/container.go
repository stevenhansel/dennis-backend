package container

import (
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/episodes"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server/responseutil"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/songs"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/votes"

	"go.uber.org/zap"
)

type Container struct {
	Log                   *zap.Logger
	Environment           config.Environment
	Config                *config.Configuration
	Responseutil          *responseutil.Responseutil
	SongService           *songs.SongService
	EpisodeService        *episodes.EpisodeService
	VoteService           *votes.VoteService
	EpisodeHttpController *episodes.EpisodeHttpController
	VoteHttpController    *votes.VoteHttpController
}

func New(
	log *zap.Logger,
	environment config.Environment,
	config *config.Configuration,
	responseutil *responseutil.Responseutil,
	songService *songs.SongService,
	episodeService *episodes.EpisodeService,
	voteService *votes.VoteService,
	episodeHttpController *episodes.EpisodeHttpController,
	voteHttpController *votes.VoteHttpController,
) *Container {
	return &Container{
		Log:                   log,
		Environment:           environment,
		Config:                config,
		Responseutil:          responseutil,
		SongService:           songService,
		EpisodeService:        episodeService,
		VoteService:           voteService,
		EpisodeHttpController: episodeHttpController,
		VoteHttpController:    voteHttpController,
	}
}
