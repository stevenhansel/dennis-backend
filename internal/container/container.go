package container

import (
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"

	"go.uber.org/zap"
)

type Container struct {
	Log     *zap.Logger
	Config  *config.Configuration
	Querier *querier.Querier
}

func New(log *zap.Logger, config *config.Configuration, querier *querier.Querier) *Container {
	return &Container{
		Log:     log,
		Config:  config,
		Querier: querier,
	}
}
