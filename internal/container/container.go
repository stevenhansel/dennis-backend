package container

import (
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"

	"go.uber.org/zap"
)

type Container struct {
	Log    *zap.Logger
	Config *config.Configuration
}

func New(log *zap.Logger, config *config.Configuration) *Container {
	return &Container{
		Config: config,
	}
}
