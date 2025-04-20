package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/user_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/user_service/internal/handler/profile"
	"github.com/DENFNC/Zappy/user_service/internal/handler/shipping"
	repo "github.com/DENFNC/Zappy/user_service/internal/repository/postgres"
	"github.com/DENFNC/Zappy/user_service/internal/service"
	psql "github.com/DENFNC/Zappy/user_service/internal/storage/postgres"
)

type App struct {
	App grpcapp.App
}

func New(
	log *slog.Logger,
	db *psql.Storage,
	port int,
) *App {

	profileRepo := repo.NewProfileRepo(db, db.Dial)
	profileSvc := service.NewProfile(log, profileRepo)
	profileHandle := profile.New(profileSvc)

	shippingRepo := repo.NewShippingRepo(db, db.Dial)
	shippingSvc := service.NewShipping(log, shippingRepo)
	shippingHandle := shipping.New(shippingSvc)

	return &App{
		App: *grpcapp.New(
			log,
			port,
			profileHandle,
			shippingHandle,
		),
	}
}
