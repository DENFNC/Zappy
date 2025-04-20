package repositories

import "github.com/DENFNC/Zappy/user_service/internal/domain/models"

type ShippingRepository interface {
	Create(address *models.Shipping) (uint32, error)
	GetByID(id int) (*models.Shipping, error)
	GetByProfileID(profileID int) ([]*models.Shipping, error)
	SetDefault(addressID int, profileID int) (uint32, error)
	Delete(id int) (uint32, error)
}
