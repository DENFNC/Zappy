package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/app"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/logger"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/config"
)

func main() {
	// Для тестов использовать ./config/config_test.yaml

	// TODO: Чтобы приложение работало нужно прописать переменные окружения
	// TODO: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY и AWS_ENDPOINT_URL

	cfg := config.MustLoad("./config/config.yaml")
	logger, err := logger.New(cfg.LogType)
	if err != nil {
		panic(err)
	}

	dbpool := MustInitDatabasePool(cfg)

	application, err := app.New(
		context.TODO(), logger,
		dbpool, cfg,
	)
	if err != nil {
		panic(err)
	}

	go application.App.MustRunGrpc()
	go application.App.MustRunHttp()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh

	logger.Info(
		"Stopped application with signal",
		"signal", sig.String(),
	)

	application.App.Stop()
}

func MustInitDatabasePool(
	cfg *config.Config,
) *postgres.Storage {
	dbpool, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		panic(err)
	}

	return dbpool
}
