package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
    _ "github.com/lib/pq"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/container"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/songs"
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
		return errtrace.Wrap(err)
	}

	config, err := config.New(environment)
	if err != nil {
		return errtrace.Wrap(err)
	}

	db, err := sqlx.Connect("postgres", config.POSTGRES_CONNECTION_URI)
	if err != nil {
    fmt.Println("caught herre")
		return errtrace.Wrap(err)
	}

	dbQuerier := database.New(db)

  songService := songs.NewService(dbQuerier)

  if err := songService.InitializeSongs(ctx); err != nil {
		return errtrace.Wrap(err)
  }

	container := container.New(log, config, songService)

	l, err := net.Listen("tcp", config.LISTEN_ADDR)
	if err != nil {
		return errtrace.Wrap(err)
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
