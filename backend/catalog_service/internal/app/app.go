package app

import (
	"context"
	"crypto/rand"
	"log/slog"
	"time"

	s3client "github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3/store"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres/repo"
	grpcapp "github.com/DENFNC/Zappy/catalog_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/catalog_service/internal/handler/category"
	"github.com/DENFNC/Zappy/catalog_service/internal/handler/product"
	productimage "github.com/DENFNC/Zappy/catalog_service/internal/handler/product_image"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/catalog_service/internal/service"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/config"
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
	paginateCoder, err := InitPaginateCoder(cfg)
	if err != nil {
		return nil, err
	}

	objectStore, err := InitStore()
	if err != nil {
		return nil, err
	}

	productRepo := repo.NewProductRepo(db, db.Dial, paginateCoder)
	productSvc := service.NewProduct(log, productRepo)
	productHandle := product.New(productSvc)

	productImageSvc := service.NewProductImage(log, objectStore)
	productImageHandle := productimage.New(productImageSvc, cfg.ObjectStore.AwsBucketImage)

	categoryRepo := repo.NewCategoryRepo(db, db.Dial, paginateCoder)
	categorySvc := service.NewCategory(log, categoryRepo)
	categoryHandle := category.New(categorySvc)

	return &App{
		App: *grpcapp.New(
			ctx,
			log,
			cfg.GRPC.Port,
			cfg.HTTP.Port,
			productHandle,
			categoryHandle,
			productImageHandle,
		),
	}, nil
}

func InitStore() (*store.Store, error) {
	s3Store, err := s3client.NewClient(
		context.TODO(),
		s3client.WithPresignExpiry(time.Second*5),
	)
	if err != nil {
		return nil, err
	}

	store := store.NewStore(s3Store)

	return store, nil
}

func InitPaginateCoder(
	cfg *config.Config,
) (*paginate.Encryptor, error) {
	paginateCoder, err := paginate.NewEncryptor(
		[]byte(cfg.PaginateSecret),
		rand.Reader,
	)
	if err != nil {
		return nil, err
	}

	return paginateCoder, nil
}
