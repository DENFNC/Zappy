package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	application, err := app.New(context.Background(), logger, cfg)
	if err != nil {
		logger.Error(
			"Error starting application",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
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
