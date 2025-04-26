package service

import (
	"context"
	"log/slog"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/DENFNC/Zappy/catalog_service/internal/domain/repositories"
)

type Category struct {
	log  *slog.Logger
	repo repositories.Category
}

func NewCategory(
	log *slog.Logger,
	repo repositories.Category,
) *Category {
	return &Category{
		log:  log,
		repo: repo,
	}
}

func (svc *Category) Create(
	ctx context.Context,
	name, parentID string,
) (string, error) {
	const op = "service.Category.Create"

	log := svc.log.With("op", op)

	categoryID, err := svc.repo.Create(ctx, name, parentID)
	if err != nil {
		log.Error(
			"Critical error",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return categoryID, nil
}

func (svc *Category) List(
	ctx context.Context,
	page, pageSize int32,
) ([]models.Category, int32, error) {
	panic("implement me")
}

func (svc *Category) Delete(
	ctx context.Context,
	categoryID string,
) error {
	panic("implement me")
}
