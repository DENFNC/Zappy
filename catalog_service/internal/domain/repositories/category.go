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
		pageSize int32,
		pageToken string,
	) ([]models.Category, any, error)
	ListByParentID(
		ctx context.Context,
		pageSize int32,
		parentID, pageToken string,
	) (*models.Category, any, error)
	Delete(
		ctx context.Context,
		categoryID string,
	) error
}
