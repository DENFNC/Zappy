package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/review_service/internal/app"
	"github.com/DENFNC/Zappy/review_service/internal/pkg/logger"
	"github.com/DENFNC/Zappy/review_service/utils/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	log := initLogger(cfg.AppLogType)
	db := initStorage(cfg.Storage.URL)

	application, err := app.New(
		context.TODO(),
		log, db,
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

	application.App.Stop()
	db.Stop()
}

func initLogger(logType string) *slog.Logger {
	log, err := logger.New(logType)
	if err != nil {
		panic(err)
	}

	return log
}

func initStorage(url string) *postgres.Storage {
	db, err := postgres.New(url)
	if err != nil {
		panic(err)
	}

	return db
}
