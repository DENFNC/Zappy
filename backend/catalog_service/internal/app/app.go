package app

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres/repo"
	grpcapp "github.com/DENFNC/Zappy/catalog_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/catalog_service/internal/handler/category"
	"github.com/DENFNC/Zappy/catalog_service/internal/handler/product"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/catalog_service/internal/service"
)

type App struct {
	App grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	db *postgres.Storage,
	grpcPort int,
	httpPort int,
	coder paginate.TokenCoder,
) (*App, error) {

	productRepo := repo.NewProductRepo(db, db.Dial, coder)
	productSvc := service.NewProduct(log, productRepo)
	productHandle := product.New(productSvc)

	categoryRepo := repo.NewCategoryRepo(db, db.Dial, coder)
	categorySvc := service.NewCategory(log, categoryRepo)
	categoryHandle := category.New(categorySvc)

	return &App{
		App: *grpcapp.New(
			ctx,
			log,
			grpcPort,
			httpPort,
			productHandle,
			categoryHandle,
		),
	}, nil
}
