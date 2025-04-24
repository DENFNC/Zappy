package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
)

type WishlistRepository interface {
	AddItem(ctx context.Context, item *models.WishlistItem) (string, error)
	GetItemByID(ctx context.Context, itemID string) (*models.WishlistItem, error)
	GetItemsByProfileID(ctx context.Context, profileID string) ([]*models.WishlistItem, error)
	RemoveItem(ctx context.Context, itemID string) error
	UpdateItem(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error)
	Exists(ctx context.Context, profileID string, productID string) (bool, error)
}
