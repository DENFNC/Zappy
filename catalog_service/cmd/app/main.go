package main

import (
	"crypto/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/catalog_service/internal/app"
	"github.com/DENFNC/Zappy/catalog_service/internal/config"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/logger"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	psql "github.com/DENFNC/Zappy/catalog_service/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	logger, err := logger.New(cfg.LogType)
	if err != nil {
		panic(err)
	}

	dbpool, err := psql.New(cfg.Postgres.URL)
	if err != nil {
		panic(err)
	}

	paginateCoder, err := paginate.NewEncryptor(
		[]byte(cfg.PaginateSecret),
		rand.Reader,
	)
	if err != nil {
		panic(err)
	}

	application, err := app.New(
		logger,
		dbpool,
		cfg.GRPC.Port,
		paginateCoder,
	)
	if err != nil {
		panic(err)
	}

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
