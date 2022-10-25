package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	if err := run(log); err != nil {
		log.Fatal("Internal Server Error", zap.String("error", fmt.Sprint(err)))
		os.Exit(1)
	}
}

func run(log *zap.Logger) error {
	log, err := zap.NewProduction()
	if err != nil {
		return errtrace.Wrap(err)
	}

	config, err := config.New(config.DEVELOPMENT)
	if err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}
