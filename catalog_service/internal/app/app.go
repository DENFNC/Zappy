package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/catalog_service/internal/app/grpc"
	psql "github.com/DENFNC/Zappy/catalog_service/internal/storage/postgres"
)

type App struct {
	App grpcapp.App
}

func New(
	log *slog.Logger,
	db *psql.Storage,
	port int,
) *App {

	return &App{
		App: *grpcapp.New(
			log,
			port,
		),
	}
}
