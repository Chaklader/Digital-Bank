package main

import (
	"context"
	"github.com/Chaklader/DigitalBank/api"
	"github.com/Chaklader/DigitalBank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/Chaklader/DigitalBank/db/sqlc"
	_ "github.com/lib/pq"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)
	runGinServer(config, store)
}

// TODO: run the DB migration
func runDBMigration(url string, source string) {

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("The GIN server is unable to create ...")
	}

	serverChannel := make(chan os.Signal, 1)
	signal.Notify(serverChannel, interruptSignals...)

	go func() {
		if err := server.Start(config.HTTPServerAddress); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("The GIN server is unable to run ....")
		}
	}()

	sig := <-serverChannel
	log.Info().Msgf("Received signal %s, gracefully shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Stop(ctx, config.HTTPServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Failed to gracefully stop server")
	}
}
