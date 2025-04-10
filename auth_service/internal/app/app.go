package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/internal/app/grpc"
	"github.com/DENFNC/Zappy/internal/infrastructure/repo"
	authservice "github.com/DENFNC/Zappy/internal/service/auth"
	psql "github.com/DENFNC/Zappy/internal/storage/postgres"
)

type App struct {
	App *grpcapp.App
}

func New(
	log *slog.Logger,
	port int,
	db *psql.Storage,
) *App {
	userRepo := repo.NewUser(db)

	authService := authservice.NewAuth(log, userRepo)

	return &App{
		App: grpcapp.New(log, port, authService),
	}
}
