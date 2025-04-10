package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/internal/app/grpc"
	authgrpc "github.com/DENFNC/Zappy/internal/grpc/auth"
)

type App struct {
	App *grpcapp.App
}

func New(log *slog.Logger, port int, auth authgrpc.Auth) *App {
	return &App{
		App: grpcapp.New(log, port, auth),
	}
}
