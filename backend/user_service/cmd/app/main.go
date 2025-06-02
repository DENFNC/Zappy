package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/Zappy/user_service/internal/app"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/logger"
	"github.com/DENFNC/Zappy/user_service/internal/utils/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	log, err := logger.New(cfg.LogType)
	if err != nil {
		panic(err)
	}

	application, err := app.New(
		context.TODO(),
		log, cfg,
	)
	if err != nil {
		log.Error(
			"Error when launching the application",
			slog.String("error", err.Error()),
		)
	}

	go application.App.MustRunGrpc()
	go application.App.MustRunHttp()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	log.Info(
		"Stopped application with signal",
		"signal", sig.String(),
	)

	application.App.Stop()
}
