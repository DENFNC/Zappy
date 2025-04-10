package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/internal/app"
	"github.com/DENFNC/Zappy/internal/config"
	"github.com/DENFNC/Zappy/internal/pkg/logger"
	psql "github.com/DENFNC/Zappy/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	logger, err := logger.New(cfg.LogType)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting application...")

	db := psql.New(cfg.Postgres.URL)

	// Инициализация сервиса.

	application := app.New(logger, cfg.GRPC.Port, db)
	go application.App.MustRun()

	// Создаем канал для получения системных сигналов.
	sigCh := make(chan os.Signal, 1)
	// Регистрируем канал для получения сигналов SIGINT и SIGTERM.
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала (блокирующее чтение из канала).
	sig := <-sigCh
	logger.Info(
		"Stopped application with signal",
		"signal", sig.String(),
	)

	// Корректное завершение работы сервиса.
	application.App.Stop()
}
