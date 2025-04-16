package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/user_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/user_service/internal/handler/profile"
	psqlrepoprofile "github.com/DENFNC/Zappy/user_service/internal/repository/postgres"
	profileservice "github.com/DENFNC/Zappy/user_service/internal/service/profile"
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

	profileRepository := psqlrepoprofile.New(db, db.Dial)
	profileService := profileservice.New(profileRepository)
	profileHandle := profile.New(profileService)

	return &App{
		App: *grpcapp.New(
			log,
			port,
			profileHandle,
		),
	}
}
