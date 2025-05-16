package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	s3client "github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3"
	awsstore "github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3/store"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/nosql/redis"
	kvstore "github.com/DENFNC/Zappy/catalog_service/internal/adapters/nosql/redis/store"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres/repo"
	grpcapp "github.com/DENFNC/Zappy/catalog_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	categoryservice "github.com/DENFNC/Zappy/catalog_service/internal/service/category"
	hookService "github.com/DENFNC/Zappy/catalog_service/internal/service/hooks"
	productservice "github.com/DENFNC/Zappy/catalog_service/internal/service/product"
	productimageservice "github.com/DENFNC/Zappy/catalog_service/internal/service/product_image"
	"github.com/DENFNC/Zappy/catalog_service/internal/transport/category"
	"github.com/DENFNC/Zappy/catalog_service/internal/transport/hooks"
	"github.com/DENFNC/Zappy/catalog_service/internal/transport/product"
	productimage "github.com/DENFNC/Zappy/catalog_service/internal/transport/product_image"
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
	paginateCoder, err := initPaginateCoder(cfg)
	if err != nil {
		return nil, err
	}

	s3Client, objectStore, err := initS3Store()
	if err != nil {
		return nil, err
	}
	kvstore := initKVStorage(cfg)
	initObjectStoreNotifyer(s3Client, cfg.ObjectStore.StagingBucket)

	productRepo := repo.NewProductRepo(db, paginateCoder)
	productSvc := productservice.NewProduct(log, productRepo)
	productHandle := product.New(productSvc)

	productImageRepo := repo.NewProductImage(db, paginateCoder)
	productImageSvc := productimageservice.NewProductImage(log, cfg, objectStore, kvstore, productImageRepo)
	productImageHandle := productimage.New(productImageSvc, cfg.ObjectStore.StagingBucket)

	categoryRepo := repo.NewCategoryRepo(db, paginateCoder)
	categorySvc := categoryservice.NewCategory(log, categoryRepo)
	categoryHandle := category.New(categorySvc)

	checkMimeSvcHook := hookService.New(productImageRepo, log, objectStore, kvstore, cfg)
	checkMimeHandleHook := hooks.New(checkMimeSvcHook)

	return &App{
		App: *grpcapp.New(
			ctx,
			log,
			cfg.GRPC.Reflection,
			cfg.GRPC.Port,
			cfg.HTTP.Port,
			productHandle,
			categoryHandle,
			productImageHandle,
			checkMimeHandleHook,
		),
	}, nil
}

func initS3Store() (*s3client.Client, *awsstore.Store, error) {
	// TODO: в prod креды должны передаваться через переменные среды
	// TODO: Измени ключи для нормальной работы в локали (в minio вкладка -> Access Keys)
	creds := credentials.NewStaticCredentialsProvider(
		"VtK99Smfp4o5KCzw1LBu",                     // * Access key
		"LeSncjLNQZoCTncz1nTJw0XL28eUy7cGEwAIuu8X", // * Secret key
		"",
	)

	client, err := s3client.NewClient(
		context.TODO(),
		s3client.WithPresignExpiry(time.Minute*15),
		s3client.WithEndpoint("http://localhost:9000"),
		s3client.WithCredentials(creds),
	)
	if err != nil {
		return nil, nil, err
	}

	store := awsstore.NewStore(client)

	return client, store, nil
}

func initObjectStoreNotifyer(
	client *s3client.Client,
	bucket string,
) {
	notify := s3client.NewNotifyer(client)
	// TODO: Временный хардкод, затем переменные будут передаваться через конфиг
	// TODO: Регистрация сделана для теста AMQP
	err := notify.RegisterNewNotify(
		context.TODO(),
		"MimeValidation",
		"arn:minio:sqs::MIME:webhook",
		bucket,
		"PUT",
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func initPaginateCoder(
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

func initKVStorage(
	cfg *config.Config,
) *kvstore.Store {
	client := redis.NewClient(
		redis.WithAddr("localhost:6379"),
		redis.WithPassword(""),
		redis.WithDB(0),
	)

	store := kvstore.New(client)

	return store
}
