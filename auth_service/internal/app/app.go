package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/internal/app/grpc"
	"github.com/DENFNC/Zappy/internal/infrastructure/repo"
	vaulttoken "github.com/DENFNC/Zappy/internal/pkg/authjwt"
	authservice "github.com/DENFNC/Zappy/internal/service/auth"
	psql "github.com/DENFNC/Zappy/internal/storage/postgres"
)

type App struct {
	App *grpcapp.App
}

func New(
	log *slog.Logger,
	db *psql.Storage,
	vault vaulttoken.VaultKMS,
	port int,
) (*App, error) {
	userRepo := repo.NewUser(db)

	authService := authservice.NewAuth(log, userRepo, vault)

	return &App{
		App: grpcapp.New(
			log,
			port,
			authService,
		),
	}, nil
}
