package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/catalog_service/internal/app"
	"github.com/DENFNC/Zappy/catalog_service/internal/config"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/logger"
	psql "github.com/DENFNC/Zappy/catalog_service/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	logger, _ := logger.New(cfg.LogType)

	dbpool, _ := psql.New(cfg.Postgres.URL)

	application := app.New(logger, dbpool, cfg.GRPC.Port)

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
