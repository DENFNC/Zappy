package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	repo "github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres/repo"
	grpcapp "github.com/DENFNC/Zappy/user_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/retry"
	"github.com/DENFNC/Zappy/user_service/internal/service"
	"github.com/DENFNC/Zappy/user_service/internal/transport/payment"
	"github.com/DENFNC/Zappy/user_service/internal/transport/profile"
	"github.com/DENFNC/Zappy/user_service/internal/transport/shipping"
	"github.com/DENFNC/Zappy/user_service/internal/utils/config"
	"github.com/jackc/pgx/v5/pgtype"
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
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	pgCoder := initPaginateCoder(cfg.PaginateSecret)

	profileRepo := repo.NewProfileRepo(storage, pgCoder)
	profileSvc := service.NewProfile(log, profileRepo)
	profileHandle := profile.New(profileSvc)

	shippingRepo := repo.NewShippingRepo(storage, pgCoder)
	shippingSvc := service.NewShipping(log, shippingRepo)
	shippingHandle := shipping.New(shippingSvc)

	paymentRepo := repo.NewPaymentRepo(storage, pgCoder)
	paymentSvc := service.NewPayment(log, paymentRepo)
	paymentHandle := payment.New(paymentSvc)

	return &App{
		App: *grpcapp.New(
			ctx,
			log,
			cfg.GRPC.Reflection,
			cfg.GRPC.URL,
			cfg.HTTP.URL,
			profileHandle,
			shippingHandle,
			paymentHandle,
		),
	}, nil
}

func initStorage(cfg *config.Config, log *slog.Logger) (*postgres.Storage, error) {
	return postgres.NewStorage(cfg.Postgres.URL)
}

func initPaginateCoder(key string) *paginate.Encryptor {
	paginateCoder, err := paginate.NewEncryptor(
		[]byte(key), rand.Reader,
	)
	if err != nil {
		panic(err)
	}
	paginate.PaginateTypeRegister(
		pgtype.Timestamp{},
		pgtype.UUID{},
	)

	return paginateCoder
}
