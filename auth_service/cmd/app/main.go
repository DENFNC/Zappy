package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/internal/app"
	"github.com/DENFNC/Zappy/internal/config"
	"github.com/DENFNC/Zappy/internal/pkg/logger"
)

func main() {
	// Загрузка конфигурации (реализация зависит от вашего пакета config).
	cfg := config.MustLoad("./config/config.yaml")

	// Инициализация логгера.
	logger, err := logger.New(cfg.LogType)
	if err != nil {
		// Если инициализация не удалась, аварийно завершаем программу.
		panic(err)
	}

	logger.Info("Starting application...")

	// Инициализация сервиса.
	application := app.New(logger, cfg.GRPC.Port)
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
