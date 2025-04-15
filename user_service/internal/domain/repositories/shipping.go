package repositories

import "github.com/DENFNC/Zappy/user_service/internal/domain/models"

type ShippingAddressRepository interface {
	Create(address *models.ShippingAddress) (int, error)
	GetByID(id int) (*models.ShippingAddress, error)
	GetByProfileID(profileID int) ([]*models.ShippingAddress, error)
	SetDefault(addressID int, profileID int) error
	Delete(id int) error
}
