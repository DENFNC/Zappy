package app

import (
	"log/slog"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	repo "github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres/repo"
	grpcapp "github.com/DENFNC/Zappy/user_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/user_service/internal/service"
	"github.com/DENFNC/Zappy/user_service/internal/transport/payment"
	"github.com/DENFNC/Zappy/user_service/internal/transport/profile"
	"github.com/DENFNC/Zappy/user_service/internal/transport/shipping"
)

type App struct {
	App grpcapp.App
}

func New(
	log *slog.Logger,
	db *postgres.Storage,
	port int,
) *App {

	profileRepo := repo.NewProfileRepo(db)
	profileSvc := service.NewProfile(log, profileRepo)
	profileHandle := profile.New(profileSvc)

	shippingRepo := repo.NewShippingRepo(db)
	shippingSvc := service.NewShipping(log, shippingRepo)
	shippingHandle := shipping.New(shippingSvc)

	paymentRepo := repo.NewPaymentRepo(db)
	paymentSvc := service.NewPayment(log, paymentRepo)
	paymentHandle := payment.New(paymentSvc)

	return &App{
		App: *grpcapp.New(
			log,
			port,
			profileHandle,
			shippingHandle,
			paymentHandle,
		),
	}
}
