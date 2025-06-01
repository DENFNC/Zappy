package app

import (
	"context"
	"crypto/rand"
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	repo "github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres/repo"
	grpcapp "github.com/DENFNC/Zappy/user_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/paginate"
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
	db *postgres.Storage,
	cfg *config.Config,
) (*App, error) {
	pgCoder := initPaginateCoder(cfg.PaginateSecret)

	profileRepo := repo.NewProfileRepo(db, pgCoder)
	profileSvc := service.NewProfile(log, profileRepo)
	profileHandle := profile.New(profileSvc)

	shippingRepo := repo.NewShippingRepo(db, pgCoder)
	shippingSvc := service.NewShipping(log, shippingRepo)
	shippingHandle := shipping.New(shippingSvc)

	paymentRepo := repo.NewPaymentRepo(db, pgCoder)
	paymentSvc := service.NewPayment(log, paymentRepo)
	paymentHandle := payment.New(paymentSvc)

	return &App{
		App: *grpcapp.New(
			ctx,
			log,
			cfg.GRPC.Reflection,
			cfg.GRPC.Port,
			cfg.HTTP.Port,
			profileHandle,
			shippingHandle,
			paymentHandle,
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
