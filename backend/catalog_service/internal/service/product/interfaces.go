package productservice

import (
	"context"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
)

type ProductRepo interface {
	Create(
		ctx context.Context,
		name, desc string,
		categoryIDs []string,
		price int64,
	) (string, error)
	GetByID(
		ctx context.Context,
		uid string,
	) (*models.Product, error)
	List(
		ctx context.Context,
		pageSize uint32,
		pageToken string,
	) ([]models.Product, string, error)
	Update(
		ctx context.Context,
		uid string,
		desc, name string,
		categoryIDs []string,
		price int64,
	) error
	Delete(
		ctx context.Context,
		uid string,
	) error
}
