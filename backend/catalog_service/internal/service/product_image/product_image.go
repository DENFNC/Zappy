package productimageservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/config"
)

type ProductImage struct {
	log *slog.Logger
	cfg *config.Config
	ObjectStorage
	KeyValueStorage
	ProductImageRepo
}

func NewProductImage(
	log *slog.Logger,
	cfg *config.Config,
	store ObjectStorage,
	kvstorage KeyValueStorage,
	repo ProductImageRepo,
) *ProductImage {
	return &ProductImage{
		log:              log,
		cfg:              cfg,
		ObjectStorage:    store,
		KeyValueStorage:  kvstorage,
		ProductImageRepo: repo,
	}
}

func (svc *ProductImage) GetByID(
	ctx context.Context,
	imageID string,
) (*models.ProductImage, error) {
	const op = "service.ProductImage.GetByID"

	log := svc.log.With("op", op)

	item, err := svc.ProductImageRepo.GetByID(ctx, imageID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return item, nil
}

func (svc *ProductImage) ListProductImages(
	ctx context.Context,
	pageSize uint32,
	pageToken, productID string,
) ([]models.ProductImage, string, error) {
	const op = "service.ProductImage.ListProductImages"

	log := svc.log.With("op", op)

	items, nextPageToken, err := svc.ProductImageRepo.List(
		ctx,
		pageSize,
		pageToken,
		productID,
	)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, "", err
	}

	return items, nextPageToken, nil
}

func (svc *ProductImage) DeleteProductImage(
	ctx context.Context,
	imageID string,
) error {
	const op = "service.ProductImage.DeleteProductImage"

	log := svc.log.With("op", op)

	objectKey, err := svc.ProductImageRepo.Delete(ctx, imageID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return err
	}
	if err := svc.ObjectStorage.DeleteObject(ctx,
		svc.cfg.ObjectStore.ImageBucket,
		objectKey,
	); err != nil {
		return err
	}

	return nil
}

func (svc *ProductImage) GetUploadURL(
	ctx context.Context,
	bucket, key, contentType string,
	productID, alt string,
) (string, error) {
	const op = "service.ProductImage.GetUploadURL"

	log := svc.log.With("op", op)

	fmt.Println(key)

	if err := svc.KeyValueStorage.HSet(ctx, key,
		map[string]any{
			"status":     "pending",
			"bucket":     svc.cfg.ObjectStore.ImageBucket,
			"product_id": productID,
			"alt":        alt,
		}); err != nil {
		return "", err
	}

	url, err := svc.ObjectStorage.PresignPut(ctx, bucket, key, contentType)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return url, nil
}

func (svc *ProductImage) GetUploadStatus(
	ctx context.Context,
	key string,
) (string, error) {
	status, err := svc.KeyValueStorage.HGet(
		ctx,
		key, "status",
	)
	if err != nil {
		return "", err
	}

	return status, nil
}
