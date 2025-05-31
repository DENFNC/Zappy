package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/user_service/internal/app"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/logger"
	"github.com/DENFNC/Zappy/user_service/internal/utils/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	log, err := logger.New(cfg.LogType)
	if err != nil {
		panic(err)
	}
	dbpool, err := postgres.NewStorage(cfg.Postgres.URL)
	if err != nil {
		log.Error(
			"Error connection to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	application, err := app.New(
		context.TODO(),
		log, dbpool,
		cfg,
	)
	if err != nil {
		log.Error(
			"Error when launching the application",
			slog.String("error", err.Error()),
		)
	}

	go application.App.MustRunGrpc()
	go application.App.MustRunHttp()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	log.Info(
		"Stopped application with signal",
		"signal", sig.String(),
	)

	dbpool.Stop()
	application.App.Stop()
}
