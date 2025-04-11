package app

import (
	"fmt"
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/internal/app/grpc"
	"github.com/DENFNC/Zappy/internal/config"
	"github.com/DENFNC/Zappy/internal/infrastructure/repo"
	"github.com/DENFNC/Zappy/internal/pkg/vault"
	authservice "github.com/DENFNC/Zappy/internal/service/auth"
	psql "github.com/DENFNC/Zappy/internal/storage/postgres"
)

type App struct {
	App *grpcapp.App
}

func New(
	log *slog.Logger,
	db *psql.Storage,
	cfg *config.Config,
) *App {
	userRepo := repo.NewUser(db)

	vault, err := vault.New(cfg.Vault.URL, cfg.Vault.Token)
	if err != nil {
		fmt.Println(err)
	}
	authService := authservice.NewAuth(log, userRepo, vault)

	return &App{
		App: grpcapp.New(
			log,
			cfg.GRPC.Port,
			authService,
		),
	}
}
