package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/internal/app"
	"github.com/DENFNC/Zappy/internal/config"
	"github.com/DENFNC/Zappy/internal/pkg/logger"
	"github.com/DENFNC/Zappy/internal/pkg/vault"
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
	vault, err := vault.New(cfg.Vault.URL, cfg.Vault.Token)
	if err != nil {
		logger.Error(
			"Error connecting to vault",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	application, err := app.New(logger, db, vault, cfg.GRPC.Port)
	if err != nil {
		logger.Error(
			"Error starting application",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
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
