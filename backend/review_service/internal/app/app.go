package app

import (
	"context"
	"crypto/rand"
	"log/slog"
	"os"
	"time"

	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres/repositories"
	grpcapp "github.com/DENFNC/Zappy/review_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/review_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/review_service/internal/pkg/retry"
	"github.com/DENFNC/Zappy/review_service/internal/service"
	"github.com/DENFNC/Zappy/review_service/internal/transport/review"
	"github.com/DENFNC/Zappy/review_service/internal/utils/config"
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
	db, err := retry.Retry(
		ctx,
		func(ctx context.Context) (*postgres.Storage, error) {
			return initStorage(cfg, log)
		},
		retry.WithMaxRetries(cfg.Storage.MaxRetries),
		retry.WithInterval(time.Duration(cfg.Storage.RetryIntervalMS)*time.Millisecond),
	)
	if err != nil {
		log.Error(
			"Couldn't connect to the database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	pgCoder := initPaginateCoder(cfg.PaginateSecret)

	reviewRepo := repositories.NewReviewRepo(db, pgCoder)
	reviewSvc := service.New(log, reviewRepo)
	reviewHandle := review.New(reviewSvc)

	return &App{
		App: *grpcapp.New(
			ctx,
			log,
			cfg.GRPC.Reflection,
			cfg.GRPC.URL,
			cfg.HTTP.URL,
			reviewHandle,
		),
	}, nil
}

func initStorage(
	cfg *config.Config,
	log *slog.Logger,
) (*postgres.Storage, error) {
	db, err := postgres.NewStorage(cfg.Storage.URL)
	if err != nil {
		log.Error(
			"Error connection to database",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return db, nil
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
