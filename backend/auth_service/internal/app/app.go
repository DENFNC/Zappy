package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/auth_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/auth_service/internal/config"
	"github.com/DENFNC/Zappy/auth_service/internal/infrastructure/repo"
	"github.com/DENFNC/Zappy/auth_service/internal/pkg/authjwt/signature"
	authservice "github.com/DENFNC/Zappy/auth_service/internal/service/auth"
	psql "github.com/DENFNC/Zappy/auth_service/internal/storage/postgres"
)

type App struct {
	App *grpcapp.App
}

func New(
	log *slog.Logger,
	db *psql.Storage,
	vault signature.VaultKMS,
	cfgVault config.ConfigVault,
	port int,
) (*App, error) {
	userRepo := repo.NewUser(db)

	authService := authservice.NewAuth(log, userRepo, vault, cfgVault)

	return &App{
		App: grpcapp.New(
			log,
			port,
			authService,
		),
	}, nil
}
