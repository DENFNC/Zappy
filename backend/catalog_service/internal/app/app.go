package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials"

	s3 "github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3"
	awsstore "github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3/store"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/nosql/redis"
	kvstore "github.com/DENFNC/Zappy/catalog_service/internal/adapters/nosql/redis/store"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres/repo"
	grpcapp "github.com/DENFNC/Zappy/catalog_service/internal/app/grpc"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/retry"
	categoryservice "github.com/DENFNC/Zappy/catalog_service/internal/service/category"
	hookService "github.com/DENFNC/Zappy/catalog_service/internal/service/hooks"
	productservice "github.com/DENFNC/Zappy/catalog_service/internal/service/product"
	productimageservice "github.com/DENFNC/Zappy/catalog_service/internal/service/product_image"
	categoryTransport "github.com/DENFNC/Zappy/catalog_service/internal/transport/category"
	hooksTransport "github.com/DENFNC/Zappy/catalog_service/internal/transport/hooks"
	productTransport "github.com/DENFNC/Zappy/catalog_service/internal/transport/product"
	productImageTransport "github.com/DENFNC/Zappy/catalog_service/internal/transport/product_image"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/config"
)

// S3Pair объединяет S3-клиент и хранилище для передачи в Retry.
type S3Pair struct {
	Client *s3.Client
	Store  *awsstore.Store
}

// App инкапсулирует основное gRPC-приложение.
type App struct {
	App grpcapp.App
}

// New создаёт и инициализирует все необходимые компоненты приложения.
func New(
	ctx context.Context,
	log *slog.Logger,
	db *postgres.Storage,
	cfg *config.Config,
) (*App, error) {
	paginateCoder := initPaginateCoder(cfg, log)
	s3Result, err := retry.Retry(
		ctx,
		func(ctx context.Context) (S3Pair, error) {
			client, store, err := initS3Store(cfg)
			if err != nil {
				return S3Pair{}, err
			}
			return S3Pair{Client: client, Store: store}, nil
		},
		retry.WithMaxRetries(cfg.ObjectStore.MaxRetries),
		retry.WithInterval(time.Duration(cfg.ObjectStore.RetryIntervalMS)*time.Millisecond),
	)
	if err != nil {
		log.Error(
			"Failed to init S3 store",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	kvStore, err := retry.Retry(
		ctx,
		func(ctx context.Context) (*kvstore.Store, error) {
			return initKVStorage(cfg, log), nil
		},
		retry.WithMaxRetries(cfg.Redis.MaxRetries),
		retry.WithInterval(time.Duration(cfg.Redis.RetryIntervalMS)*time.Millisecond),
	)
	if err != nil {
		log.Error(
			"Couldn't connect to Redis",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	initObjectStoreNotifier(s3Result.Client, log, cfg)

	productRepo := repo.NewProductRepo(db, paginateCoder)
	productSvc := productservice.NewProduct(log, productRepo)
	productHandler := productTransport.New(productSvc)

	productImageRepo := repo.NewProductImage(db, paginateCoder)
	productImageSvc := productimageservice.NewProductImage(
		log,
		cfg,
		s3Result.Store,
		kvStore,
		productImageRepo,
	)
	productImageHandler := productImageTransport.New(productImageSvc, cfg.ObjectStore.StagingBucket)

	categoryRepo := repo.NewCategoryRepo(db, paginateCoder)
	categorySvc := categoryservice.NewCategory(log, categoryRepo)
	categoryHandler := categoryTransport.New(categorySvc)

	checkMimeSvc := hookService.New(productImageRepo, log, s3Result.Store, kvStore, cfg)
	checkMimeHandler := hooksTransport.New(checkMimeSvc)

	app := grpcapp.New(
		ctx,
		log,
		cfg.GRPC.Reflection,
		cfg.GRPC.URL,
		cfg.HTTP.URL,
		productHandler,
		categoryHandler,
		productImageHandler,
		checkMimeHandler,
	)

	return &App{App: *app}, nil
}

func initS3Store(cfg *config.Config) (*s3.Client, *awsstore.Store, error) {
	creds := credentials.NewStaticCredentialsProvider(
		cfg.ObjectStore.AccessKey,
		cfg.ObjectStore.SecretKey,
		"",
	)

	client, err := s3.NewClient(
		context.TODO(),
		s3.WithPresignExpiry(15*time.Minute),
		s3.WithEndpoint(cfg.ObjectStore.ObjectOrigin),
		s3.WithCredentials(creds),
	)
	if err != nil {
		return nil, nil, err
	}

	if err := client.EnsureBucketExists(
		context.TODO(),
		cfg.ObjectStore.ImageBucket,
		cfg.ObjectStore.StagingBucket,
	); err != nil {
		return nil, nil, fmt.Errorf("failed to ensure buckets exist: %w", err)
	}

	store := awsstore.NewStore(client)
	return client, store, nil
}

func initObjectStoreNotifier(
	client *s3.Client,
	log *slog.Logger,
	cfg *config.Config,
) {
	notifier := s3.NewNotifyer(client)
	if err := notifier.RegisterNewNotify(
		context.TODO(),
		"MimeValidation",
		"arn:minio:sqs::MIME:webhook",
		cfg.ObjectStore.StagingBucket,
		"PUT",
	); err != nil {
		log.Error(
			"Failed to register object store notifier",
			slog.String("error", err.Error()),
		)
	}
}

func initPaginateCoder(
	cfg *config.Config,
	log *slog.Logger,
) *paginate.Encryptor {
	encryptor, err := paginate.NewEncryptor(
		[]byte(cfg.PaginateSecret),
		rand.Reader,
	)
	if err != nil {
		log.Error(
			"Failed to init paginate coder",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	return encryptor
}

func initKVStorage(
	cfg *config.Config,
	log *slog.Logger,
) *kvstore.Store {
	client := redis.NewClient(
		redis.WithAddr(cfg.Redis.URL),
		redis.WithPassword(cfg.Redis.Password),
		redis.WithDB(cfg.Redis.DB),
	)
	return kvstore.New(client, log)
}
