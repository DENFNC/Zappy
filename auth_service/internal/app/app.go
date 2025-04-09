package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/internal/app/grpc"
)

type App struct {
	App *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	return &App{
		App: grpcapp.New(log, port),
	}
}
