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
		log.Error(
			"Failed to init paginate coder",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	s3Client, objectStore, err := initS3Store(cfg, cfg.ObjectStore.ObjectOrigin, log)
	if err != nil {
		log.Error(
			"Failed to init S3 store",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	kvstore := initKVStorage(cfg, log)
	initObjectStoreNotifyer(s3Client, cfg.ObjectStore.StagingBucket, cfg.ObjectStore.ObjectOrigin, log)

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
		},
		nil
}

func initS3Store(cfg *config.Config, objectOrigin string, log *slog.Logger) (*s3client.Client, *awsstore.Store, error) {
	creds := credentials.NewStaticCredentialsProvider(
		cfg.ObjectStore.AccessKey,
		cfg.ObjectStore.SecretKey,
		"",
	)

	client, err := s3client.NewClient(
		context.TODO(),
		s3client.WithPresignExpiry(time.Minute*15),
		s3client.WithEndpoint(objectOrigin),
		s3client.WithCredentials(creds),
	)
	if err != nil {
		return nil, nil, err
	}
	if err := client.EnsureBucketExists(context.TODO(), cfg.ObjectStore.ImageBucket, cfg.ObjectStore.StagingBucket); err != nil {
		return nil, nil, fmt.Errorf("failed to ensure buckets exist: %w", err)
	}
	store := awsstore.NewStore(client)

	return client, store, nil
}

func initObjectStoreNotifyer(
	client *s3client.Client,
	bucket string,
	objectOrigin string,
	log *slog.Logger,
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
		log.Error(
			"Failed to register new notifyer",
			slog.String("error", err.Error()),
		)
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
	log *slog.Logger,
) *kvstore.Store {
	client := redis.NewClient(
		redis.WithAddr("localhost:6379"),
		redis.WithPassword(""),
		redis.WithDB(0),
	)

	store := kvstore.New(client, log)

	return store
}
