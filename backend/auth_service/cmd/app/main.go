package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/auth_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/auth_service/internal/app"
	"github.com/DENFNC/Zappy/auth_service/internal/pkg/logger"
	"github.com/DENFNC/Zappy/auth_service/internal/utils/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	logger, err := logger.New(cfg.AppLog)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting application...")

	db, err := postgres.NewStorage(cfg.Postgres.URL, logger)
	if err != nil {
		logger.Error(
			"Error connection to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Инициализация сервиса.
	// vault, err := vault.New(cfg.Vault.URL, cfg.Vault.Token)
	// if err != nil {
	// 	logger.Error(
	// 		"Error connecting to vault",
	// 		slog.String("error", err.Error()),
	// 	)
	// 	os.Exit(1)
	// }
	application, err := app.New(context.Background(),
		logger,
		db,
		// vault,
		cfg,
	)
	if err != nil {
		logger.Error(
			"Error starting application",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	go application.App.MustRunGrpc()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	logger.Info(
		"Stopped application with signal",
		"signal", sig.String(),
	)

	application.App.Stop()
}
