package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
)

type CategoryRepo interface {
	Create(
		ctx context.Context,
		name, parentID string,
	) (string, error)
	GetByID(
		ctx context.Context,
		categoryID string,
	) (*models.Category, error)
	List(
		ctx context.Context,
		pageSize uint,
		pageToken string,
	) ([]models.Category, string, error)
	ListByParentID(
		ctx context.Context,
		pageSize uint,
		parentID, pageToken string,
	) (*models.Category, any, error)
	Delete(
		ctx context.Context,
		categoryID string,
	) error
}

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
		price int64,
	) error
	Delete(
		ctx context.Context,
		uid string,
	) error
}
