package app

import (
	"context"
	crand "crypto/rand"
	"log/slog"
	"os"
	"time"

	"github.com/DENFNC/Zappy/auth_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/auth_service/internal/adapters/sql/postgres/repo"
	grpcapp "github.com/DENFNC/Zappy/auth_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/auth_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/auth_service/internal/pkg/retry"
	"github.com/DENFNC/Zappy/auth_service/internal/pkg/vault"
	"github.com/DENFNC/Zappy/auth_service/internal/service"
	"github.com/DENFNC/Zappy/auth_service/internal/transport/auth"
	"github.com/DENFNC/Zappy/auth_service/internal/utils/config"
)

type App struct {
	App grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) (*App, error) {
	storage, err := retry.Retry(
		ctx,
		func(ctx context.Context) (*postgres.Storage, error) {
			return initStorage(cfg, log)
		},
		retry.WithMaxRetries(cfg.Postgres.MaxRetries),
		retry.WithInterval(time.Duration(cfg.Postgres.RetryIntervalMS)*time.Millisecond),
	)
	if err != nil {
		log.Error(
			"Couldn't connect to the database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	pgCoder := initPaginateCoder(cfg, log)
	vault, err := retry.Retry(
		ctx,
		func(ctx context.Context) (*vault.Vault, error) {
			return initVault(cfg, log)
		},
		retry.WithMaxRetries(cfg.Vault.MaxRetries),
		retry.WithInterval(time.Duration(cfg.Postgres.RetryIntervalMS)*time.Millisecond),
	)
	if err != nil {
		log.Error(
			"Couldn't connect to the vault",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	authRepo := repo.NewUser(storage, pgCoder)
	authSvc := service.NewAuth(log, authRepo, vault, cfg)
	authHandle := auth.New(authSvc)

	return &App{
			App: *grpcapp.New(
				ctx,
				log,
				cfg.GRPC.Reflection,
				cfg.GRPC.URL,
				cfg.HTTP.URL,
				authHandle,
			),
		},
		nil
}

func initVault(
	cfg *config.Config,
	log *slog.Logger,
) (*vault.Vault, error) {
	vault, err := vault.New(cfg.Vault.URL, cfg.Vault.Token)
	if err != nil {
		log.Error(
			"Error connecting to vault",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return vault, nil
}

func initPaginateCoder(
	cfg *config.Config,
	log *slog.Logger,
) *paginate.Encryptor {
	pgCoder, err := paginate.NewEncryptor(
		[]byte(cfg.PaginateSecret),
		crand.Reader,
	)
	if err != nil {
		log.Error(
			"Error init paginate coder",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	return pgCoder
}

func initStorage(
	cfg *config.Config,
	log *slog.Logger,
) (*postgres.Storage, error) {
	db, err := postgres.NewStorage(cfg.Postgres.URL, log)
	if err != nil {
		log.Error(
			"Error connection to database",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return db, nil
}
