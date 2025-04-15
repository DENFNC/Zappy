package repositories

import "github.com/DENFNC/Zappy/user_service/internal/domain/models"

type WishlistRepository interface {
	AddItem(item *models.WishlistItem) (int, error)
	GetItemsByProfileID(profileID int) ([]*models.WishlistItem, error)
	RemoveItem(itemID int) error
	Exists(profileID int, productID int) (bool, error)
}
