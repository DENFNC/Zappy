package app

import (
	"log/slog"

	grpcapp "github.com/DENFNC/Zappy/catalog_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/catalog_service/internal/handler/category"
	"github.com/DENFNC/Zappy/catalog_service/internal/handler/product"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/catalog_service/internal/repo"
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
	coder paginate.TokenCoder,
) (*App, error) {

	productRepo := repo.NewProductRepo(db, db.Dial)
	productSvc := service.NewProduct(log, productRepo)
	productHandle := product.New(productSvc)

	categoryRepo, err := repo.NewCategoryRepo(db, db.Dial, coder)
	if err != nil {
		return nil, err
	}
	categorySvc := service.NewCategory(log, categoryRepo)
	categoryHandle := category.New(categorySvc)

	return &App{
		App: *grpcapp.New(
			log,
			port,
			productHandle,
			categoryHandle,
		),
	}, nil
}
