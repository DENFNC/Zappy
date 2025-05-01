package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/user_service/internal/app"
	"github.com/DENFNC/Zappy/user_service/internal/config"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/logger"
	psql "github.com/DENFNC/Zappy/user_service/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	logger, err := logger.New(cfg.LogType)
	if err != nil {
		panic(err)
	}
	dbpool, err := psql.New(cfg.Postgres.URL)
	if err != nil {
		logger.Error(
			"Error connection to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	application := app.New(
		logger,
		dbpool,
		cfg.GRPC.Port,
	)

	go application.App.MustRun()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	logger.Info(
		"Stopped application with signal",
		"signal", sig.String(),
	)

	dbpool.Stop()
	application.App.Stop()
}
