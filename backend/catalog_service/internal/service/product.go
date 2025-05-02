package service

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/DENFNC/Zappy/catalog_service/internal/domain/repositories"
)

type Product struct {
	log  *slog.Logger
	repo repositories.Product
}

func NewProduct(
	log *slog.Logger,
	repo repositories.Product,
) *Product {
	return &Product{
		log:  log,
		repo: repo,
	}
}

func (svc *Product) Create(
	ctx context.Context,
	name, desc string,
	categoryIDs []string,
	price int64,
) (string, error) {
	const op = "service.Product.Create"

	log := svc.log.With("op", op)

	productID, err := svc.repo.Create(ctx, name, desc, categoryIDs, price)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return productID, nil
}

func (svc *Product) Get(
	ctx context.Context,
	productID string,
) (*models.Product, error) {
	const op = "service.Product.Get"

	log := svc.log.With("op", op)

	product, err := svc.repo.GetByID(ctx, productID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return product, nil
}

func (svc *Product) List(
	ctx context.Context,
	pageSize uint32,
	pageToken string,
) ([]models.Product, string, error) {
	const op = "service.Product.List"

	log := svc.log.With("op", op)

	products, nextPageToken, err := svc.repo.List(ctx, pageSize, pageToken)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return nil, "", err
	}

	return products, nextPageToken, nil
}

func (svc *Product) Update(
	ctx context.Context,
	productID, name, description string,
	price int64,
	currency string,
	categoryIDs []string,
	isActive *bool,
) (*models.Product, error) {
	panic("implement me")
}

func (svc *Product) Delete(
	ctx context.Context,
	productID string,
) error {
	const op = "service.Product.Delete"

	log := svc.log.With("op", op)

	if err := svc.repo.Delete(ctx, productID); err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}
