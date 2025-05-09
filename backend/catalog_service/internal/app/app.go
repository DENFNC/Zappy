package app

import (
	"context"
	"crypto/rand"
	"fmt"
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
	"github.com/aws/aws-sdk-go-v2/credentials"
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

	client, objectStore, err := InitStore()
	if err != nil {
		return nil, err
	}

	InitStoreNotifyer(client)

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

func InitStore() (*s3client.Client, *store.Store, error) {
	// TODO: в prod креды должны передаваться через переменные среды
	// TODO: Измени ключи для нормальной работы в локали (в minio вкладка -> Access Keys)
	creds := credentials.NewStaticCredentialsProvider(
		"IODsBpSGHPIdNXFKWltv",
		"YAkkti2wknir6WkUarAxKVZHfWMzGsr4mjE70T2U",
		"",
	)

	client, err := s3client.NewClient(
		context.TODO(),
		s3client.WithPresignExpiry(time.Second*5),
		s3client.WithEndpoint("http://localhost:9000"),
		s3client.WithCredentials(creds),
	)
	if err != nil {
		return nil, nil, err
	}

	store := store.NewStore(client)

	return client, store, nil
}

func InitStoreNotifyer(
	client *s3client.Client,
) {
	notify := s3client.NewNotifyer(client)

	// TODO: Временный хардкод, затем переменные будут передаваться через конфиг
	// TODO: Регистрация сделана для теста AMQP
	err := notify.RegisterNewNotify(
		context.TODO(),
		"MimeValidation",
		"arn:minio:sqs::IMAGE:amqp",
		"test-bucket",
		"PUT",
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
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
