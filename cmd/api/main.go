package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/container"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server"

	"go.uber.org/zap"
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
	ctx := context.Background()

	environment := config.DEVELOPMENT
	flag.Var(
		&environment,
		"env",
		"application environment, could be either (development|staging|production)",
	)
	flag.Parse()

	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	config, err := config.New(environment)
	if err != nil {
		return err
	}

	container := container.New(log, config)

	l, err := net.Listen("tcp", config.LISTEN_ADDR)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Server listening on http://%v", l.Addr()))

	s := server.New(container)

	errc := make(chan error, 1)
	go func() {
		errc <- s.HTTPServer.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Info("Failed to serve server", zap.String("error", fmt.Sprint(err)))
	case sig := <-sigs:
		log.Info("Terminating server", zap.String("signal", fmt.Sprint(sig)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.HTTPServer.Shutdown(ctx)
}
