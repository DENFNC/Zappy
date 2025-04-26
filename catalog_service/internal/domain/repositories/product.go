package repositories

import (
	"context"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type Product interface {
	Create(
		ctx context.Context,
		desc, name string,
		price pgtype.Numeric,
	) (string, error)
	GetByID(
		ctx context.Context,
		uid string,
	) (*models.Product, error)
	List(
		ctx context.Context,
	) ([]models.Product, error)
	Update(
		ctx context.Context,
		uid string,
		desc, name string,
		price pgtype.Numeric,
	) error
	Delete(
		ctx context.Context,
		uid string,
	) error
}
