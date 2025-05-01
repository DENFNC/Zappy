package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
)

type Category interface {
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
