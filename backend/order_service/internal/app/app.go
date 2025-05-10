package app

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/order_service/internal/adapters/sql/postgres"
	grpcapp "github.com/DENFNC/Zappy/order_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/order_service/internal/utils/config"
)

type App struct {
	App *grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	db *postgres.Storage,
	cfg *config.Config,
) (*App, error) {

	// orderRepo := repo.NewOrderRepo(db)
	// orderSvc := serivce.NewOrder(log, orderRepo)
	// orderHandler := product.New(orderSvc)

	grpcApp := grpcapp.New(ctx, log, cfg.GRPC.Port)

	return &App{
		App: grpcApp,
	}, nil

	// storage, err := InitStore()
	// if err != nil {
	// 	return nil, err
	// }

	// 	productRepo := repo.NewProductRepo(db, db.Dial, coder)
	// productSvc := service.NewProduct(log, productRepo)
	// productHandle := product.New(productSvc)

	// categoryRepo := repo.NewCategoryRepo(db, db.Dial, coder)
	// categorySvc := service.NewCategory(log, categoryRepo)
	// categoryHandle := category.New(categorySvc)

}

// func InitStore() (*store)
