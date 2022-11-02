package api

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/container"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/episodes"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server/responseutil"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/socket"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/songs"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/votes"
)

func run(environment config.Environment) {
	log, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	if err := internalRun(environment, log); err != nil {
		log.Fatal("Internal Server Error", zap.String("error", fmt.Sprint(err)))
		os.Exit(1)
	}
}

func internalRun(environment config.Environment, log *zap.Logger) error {
	ctx := context.Background()

	log, err := zap.NewProduction()
	if err != nil {
		return errtrace.Wrap(err)
	}

	config, err := config.New(environment)
	if err != nil {
		return errtrace.Wrap(err)
	}

	responseutil := responseutil.New(log)

	db, err := sqlx.Connect("postgres", config.POSTGRES_CONNECTION_URI)
	if err != nil {
		return errtrace.Wrap(err)
	}
	dbQuerier := database.New(db)

	socketState := socket.NewSocketState(
		rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		512,
	)

	songService := songs.NewService(dbQuerier)
	if err := songService.InitializeSongs(ctx); err != nil {
		return errtrace.Wrap(err)
	}

	episodeService := episodes.NewService(dbQuerier)
	if err := episodeService.SynchronizeThumbnails(ctx); err != nil {
		return errtrace.Wrap(err)
	}

	voteService := votes.NewService(dbQuerier, socketState)

	episodeHttpController := episodes.NewEpisodeHttpController(responseutil, episodeService)
	voteHttpController := votes.NewVoteHttpController(responseutil, voteService)

	container := container.New(
		log,
		environment,
		config,
		responseutil,
		socketState,
		songService,
		episodeService,
		voteService,
		episodeHttpController,
		voteHttpController,
	)

	l, err := net.Listen("tcp", config.LISTEN_ADDR)
	if err != nil {
		return errtrace.Wrap(err)
	}

	log.Info(fmt.Sprintf("Server listening on http://%v", l.Addr()))

	s := server.New(container)

	errc := make(chan error, 1)
	go func() {
		errc <- s.Server.Serve(l)
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

	return s.Server.Shutdown(ctx)
}
