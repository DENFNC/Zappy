package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/order_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/order_service/internal/app"
	"github.com/DENFNC/Zappy/order_service/internal/pkg/logger"
	"github.com/DENFNC/Zappy/order_service/internal/utils/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")

	logger, err := logger.New(cfg.LogType)
	if err != nil {
		panic(err)
	}

	dbpool, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		panic(err)
	}

	application, err := app.New(
		context.TODO(),
		logger,
		dbpool, cfg,
	)
	if err != nil {
		panic(err)
	}

	// TODO: Добавить запуск http сервера
	go application.App.MustRunGrpc()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigCh

	logger.Info(
		"Stopped application with signal",
		"signal", sig.String(),
	)

	application.App.Stop()
	logger.Info("Gracefully stopped")
}
