package service

import (
	"context"
	"log/slog"

	s3store "github.com/DENFNC/Zappy/catalog_service/internal/adapters/aws/s3/store"
)

type ProductImage struct {
	log   *slog.Logger
	store *s3store.Store
}

func NewProductImage(
	log *slog.Logger,
	store *s3store.Store,
) *ProductImage {
	return &ProductImage{
		log:   log,
		store: store,
	}
}

func (svc *ProductImage) GetUploadURL(
	ctx context.Context,
	bucket, key, contentType string,
) (string, error) {
	const op = "service.ProductImage.GetUploadURL"

	log := svc.log.With("op", op)

	url, err := svc.store.PresignPut(ctx, bucket, key, contentType)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return url, nil
}
