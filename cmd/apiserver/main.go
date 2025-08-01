package main

import (
	"async_api/apiserver"
	"async_api/config"
	"async_api/store"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf, err := config.New()
	if err != nil {
		return err
	}

	db, err := store.NewPostgresDB(conf)
	if err != nil {
		return err
	}
	dataStore := store.New(db)

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	jwtManager := apiserver.NewJwtManager(conf)
	server := apiserver.New(conf, logger, dataStore, jwtManager)
	if err := server.Start(ctx); err != nil {
		return err
	}

	return nil
}
