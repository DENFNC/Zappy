package main

import (
	"github.com/DENFNC/Zappy/internal/app"
	"github.com/DENFNC/Zappy/internal/config"
	"github.com/DENFNC/Zappy/internal/pkg/logger"
)

func main() {
	// TODO: Реализовать конфиг
	cfg := config.MustLoad("./config/config.yaml")
	// TODO: Реализовать логирование
	logger, err := logger.New(cfg.LogType)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting Zappy...")
	// TODO: Реализовать обработку сигналов
	// TODO: Запустить сервис
	application := app.New(logger, cfg.GRPC.Port)
	application.App.MustRun()
	// TODO:
}
