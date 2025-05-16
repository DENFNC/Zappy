package hooks

import (
	"context"
	"io"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
)

type ProductImageRepo interface {
	Create(
		ctx context.Context,
		image *models.ProductImage,
	) (string, error)
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
	GetObjectRange(
		ctx context.Context,
		bucket, key string,
		byteRange string, // * Используется строка вида (bytes=0-99)
	) (io.ReadCloser, error)
	DeleteObject(
		ctx context.Context,
		bucket, key string,
	) error
	CopyObject(
		ctx context.Context,
		bucket, key string,
		source string,
	) error
}

type KeyValueStorage interface {
	Set(
		ctx context.Context,
		key, value string,
	) error
	HSet(
		ctx context.Context,
		key string,
		values map[string]any,
	) error
	HGet(
		ctx context.Context,
		key, field string,
	) (string, error)
	HGetAll(
		ctx context.Context,
		key string,
	) (map[string]string, error)
}
