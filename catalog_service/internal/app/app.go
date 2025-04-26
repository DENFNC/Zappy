package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/catalog_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/catalog_service/internal/handler/product"
	repo "github.com/DENFNC/Zappy/catalog_service/internal/repo/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/service"
	psql "github.com/DENFNC/Zappy/catalog_service/internal/storage/postgres"
)

type App struct {
	App grpcapp.App
}

func New(
	log *slog.Logger,
	db *psql.Storage,
	port int,
) *App {

	productRepo := repo.NewProductRepo(db, db.Dial)
	productSvc := service.NewProduct(log, productRepo)
	productHandle := product.New(productSvc)

	return &App{
		App: *grpcapp.New(
			log,
			port,
			productHandle,
		),
	}
}
