package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/domain/repositories"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
)

type WishlistService struct {
	log  *slog.Logger
	repo repositories.WishlistRepository
}

func NewWishlist(log *slog.Logger, repo repositories.WishlistRepository) *WishlistService {
	return &WishlistService{
		log:  log,
		repo: repo,
	}
}

func (s *WishlistService) CreateItem(ctx context.Context, profileID, productID string) (string, error) {
	const op = "service.WishlistService.CreateItem"
	log := s.log.With("op", op)

	item := &models.WishlistItem{
		ProfileID: profileID,
		ProductID: productID,
		AddedAt:   time.Now(),
		IsActive:  true,
	}

	itemID, err := s.repo.AddItem(ctx, item)
	if err != nil {
		log.Error("Failed to add wishlist item", slog.String("error", err.Error()))
		return "", errpkg.New("CREATE_ERROR", "Failed to create wishlist item", err)
	}

	return itemID, nil
}

func (s *WishlistService) GetItem(ctx context.Context, itemID string) (*models.WishlistItem, error) {
	const op = "service.WishlistService.GetItem"
	log := s.log.With("op", op)

	item, err := s.repo.GetItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			log.Error("Not found", slog.String("itemID", itemID))
			return nil, errpkg.ErrNotFound
		}
		log.Error("Failed to get wishlist item", slog.String("error", err.Error()))
		return nil, errpkg.New("GET_ERROR", "Failed to get wishlist item", err)
	}

	return item, nil
}

func (s *WishlistService) UpdateItem(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error) {
	const op = "service.WishlistService.UpdateItem"
	log := s.log.With("op", op)

	updatedItem, err := s.repo.UpdateItem(ctx, item)
	if err != nil {
		if errors.Is(err, errpkg.ErrNotFound) {
			log.Error("Wishlist item not found", slog.String("itemID", item.ItemID))
			return nil, errpkg.ErrNotFound
		}
		log.Error("Failed to update wishlist item", slog.String("error", err.Error()))
		return nil, errpkg.New("UPDATE_ERROR", "Failed to update wishlist item", err)
	}

	return updatedItem, nil
}

func (s *WishlistService) DeleteItem(ctx context.Context, itemID string) error {
	const op = "service.WishlistService.DeleteItem"
	log := s.log.With("op", op)

	if err := s.repo.RemoveItem(ctx, itemID); err != nil {
		log.Error("Failed to delete wishlist item", slog.String("error", err.Error()))
		return errpkg.New("DELETE_ERROR", "Failed to delete wishlist item", err)
	}

	return nil
}

func (s *WishlistService) ListItems(ctx context.Context, profileID string) ([]*models.WishlistItem, error) {
	const op = "service.WishlistService.ListItems"
	log := s.log.With("op", op)

	items, err := s.repo.GetItemsByProfileID(ctx, profileID)
	if err != nil {
		log.Error("Failed to list wishlist items", slog.String("error", err.Error()))
		return nil, errpkg.New("LIST_ERROR", "Failed to list wishlist items", err)
	}

	return items, nil
}
