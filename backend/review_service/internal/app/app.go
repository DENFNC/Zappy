package app

import (
	"context"
	"crypto/rand"
	"log/slog"

	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres/repositories"
	grpcapp "github.com/DENFNC/Zappy/review_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/review_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/review_service/internal/service"
	"github.com/DENFNC/Zappy/review_service/internal/transport/review"
	"github.com/DENFNC/Zappy/review_service/utils/config"
	"github.com/jackc/pgx/v5/pgtype"
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
	pgCoder := initPaginateCoder(cfg.PaginateSecret)

	reviewRepo := repositories.NewReviewRepo(db, pgCoder)
	reviewSvc := service.New(log, reviewRepo)
	reviewHandle := review.New(reviewSvc)

	return &App{
		App: *grpcapp.New(
			ctx,
			log,
			cfg.GRPC.Reflection,
			cfg.GRPC.Port,
			cfg.HTTP.Port,
			reviewHandle,
		),
	}, nil
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
