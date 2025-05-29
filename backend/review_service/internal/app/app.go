package app

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres"
	grpcapp "github.com/DENFNC/Zappy/review_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/review_service/utils/config"
)

type App struct {
	App grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	db *postgres.Storage,
	cfg *config.Config,
) (*App, error) {
	return &App{
		App: *grpcapp.New(
			ctx,
			log,
			cfg.GRPC.Reflection,
			cfg.GRPC.Port,
			cfg.HTTP.Port,
		),
	}, nil
}
