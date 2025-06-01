package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
)

type WishlistRepository interface {
	Create(
		ctx context.Context,
		item *models.WishlistItem,
	) (string, error)
	GetItemByID(
		ctx context.Context,
		itemID string,
	) (*models.WishlistItem, error)
	GetItemsByProfileID(
		ctx context.Context,
		profileID string,
	) ([]*models.WishlistItem, error)
	Delete(
		ctx context.Context,
		itemID string,
	) error
	Update(
		ctx context.Context,
		item *models.WishlistItem,
	) (*models.WishlistItem, error)
	Exists(
		ctx context.Context,
		profileID string,
		productID string,
	) (bool, error)
	List(
		ctx context.Context,
		profileID string,
		pageSize uint,
		pageToken string,
	) ([]*models.WishlistItem, string, error)
}

type Wishlist struct {
	*slog.Logger
	WishlistRepository
}

func NewWishlist(
	log *slog.Logger,
	repo WishlistRepository,
) *Wishlist {
	return &Wishlist{
		Logger:             log,
		WishlistRepository: repo,
	}
}

func (svc *Wishlist) CreateItem(
	ctx context.Context,
	profileID string,
	productID string,
) (string, error) {
	const op = "service.Wishlist.CreateItem"

	log := svc.Logger.With("op", op)

	item := &models.WishlistItem{
		ProfileID: profileID,
		ProductID: productID,
		AddedAt:   time.Now(),
		IsActive:  true,
	}

	itemID, err := svc.WishlistRepository.Create(ctx, item)
	if err != nil {
		log.Error(
			"Failed to add wishlist item",
			slog.String("error", err.Error()),
		)
		return "", errpkg.New("CREATE_ERROR", "Failed to create wishlist item", err)
	}

	return itemID, nil
}

func (svc *Wishlist) GetItem(
	ctx context.Context,
	itemID string,
) (*models.WishlistItem, error) {
	const op = "service.Wishlist.GetItem"

	log := svc.Logger.With("op", op)

	item, err := svc.WishlistRepository.GetItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			log.Error("Not found", slog.String("itemID", itemID))
			return nil, errpkg.ErrNotFound
		}
		log.Error(
			"Failed to get wishlist item",
			slog.String("error", err.Error()),
		)
		return nil, errpkg.New("GET_ERROR", "Failed to get wishlist item", err)
	}

	return item, nil
}

func (svc *Wishlist) UpdateItem(
	ctx context.Context,
	item *models.WishlistItem,
) (*models.WishlistItem, error) {
	const op = "service.Wishlist.UpdateItem"

	log := svc.Logger.With("op", op)

	updatedItem, err := svc.WishlistRepository.Update(ctx, item)
	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			log.Error("Wishlist item not found", slog.String("itemID", item.ItemID))
			return nil, errpkg.ErrNotFound
		}
		log.Error(
			"Failed to update wishlist item",
			slog.String("error", err.Error()),
		)
		return nil, errpkg.New("UPDATE_ERROR", "Failed to update wishlist item", err)
	}

	return updatedItem, nil
}

func (svc *Wishlist) DeleteItem(
	ctx context.Context,
	itemID string,
) error {
	const op = "service.Wishlist.DeleteItem"

	log := svc.Logger.With("op", op)

	if err := svc.WishlistRepository.Delete(ctx, itemID); err != nil {
		log.Error(
			"Failed to delete wishlist item",
			slog.String("error", err.Error()),
		)
		return errpkg.New("DELETE_ERROR", "Failed to delete wishlist item", err)
	}

	return nil
}

func (svc *Wishlist) ListItems(
	ctx context.Context,
	profileID string,
	pageSize uint,
	pageToken string,
) ([]*models.WishlistItem, string, error) {
	const op = "service.Wishlist.ListItems"

	log := svc.Logger.With("op", op)

	items, nextPageToken, err := svc.WishlistRepository.List(
		ctx,
		profileID,
		pageSize,
		pageToken,
	)
	if err != nil {
		log.Error(
			"Failed to list wishlist items",
			slog.String("error", err.Error()),
		)
		return nil, "", errpkg.New("LIST_ERROR", "Failed to list wishlist items", err)
	}

	return items, nextPageToken, nil
}
