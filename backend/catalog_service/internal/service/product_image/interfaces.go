package productimageservice

import (
	"context"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
)

type ProductImageRepo interface {
	GetByID(
		ctx context.Context,
		uid string,
	) (*models.ProductImage, error)
	List(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
		productID string,
	) ([]models.ProductImage, string, error)
	Delete(
		ctx context.Context,
		imageID string,
	) (string, error)
}

type ObjectStorage interface {
	PresignPut(
		ctx context.Context,
		bucket, key, contentType string,
	) (string, error)
	DeleteObject(
		ctx context.Context,
		bucket, key string,
	) error
}

type KeyValueStorage interface {
	SetEx(
		ctx context.Context,
		key string,
		value any,
		expire time.Duration,
	) error
	Get(
		ctx context.Context,
		key string,
	) (string, error)
	HSet(
		ctx context.Context,
		key string,
		values map[string]any,
	) error
	HGet(
		ctx context.Context,
		key, field string,
	) (string, error)
}
