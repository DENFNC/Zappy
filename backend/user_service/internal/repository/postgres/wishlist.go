package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	psql "github.com/DENFNC/Zappy/user_service/internal/storage/postgres"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
)

type WishlistRepo struct {
	*psql.Storage
	goqu *goqu.DialectWrapper
}

func NewWishlistRepo(
	db *psql.Storage,
	goqu *goqu.DialectWrapper,
) *WishlistRepo {
	return &WishlistRepo{
		Storage: db,
		goqu:    goqu,
	}
}

func (r *WishlistRepo) AddItem(ctx context.Context, item *models.WishlistItem) (string, error) {
	id := r.NewV7().String()
	stmt, args, err := r.goqu.Insert("wishlist_item").
		Rows(goqu.Record{
			"item_id":    id,
			"profile_id": item.ProfileID,
			"product_id": item.ProductID,
			"added_at":   item.AddedAt,
			"is_active":  item.IsActive,
		}).
		Returning("item_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", fmt.Errorf("build AddItem query: %w", err)
	}

	var itemID string
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(&itemID); err != nil {
		return "", fmt.Errorf("exec AddItem query: %w", err)
	}

	return itemID, nil
}

func (r *WishlistRepo) GetItemByID(ctx context.Context, itemID string) (*models.WishlistItem, error) {
	stmt, args, err := r.goqu.Select(
		"item_id",
		"profile_id",
		"product_id",
		"added_at",
		"is_active",
	).From("wishlist_item").
		Where(goqu.Ex{"item_id": itemID}).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build GetItemByID query: %w", err)
	}

	var item models.WishlistItem
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(
		&item.ItemID,
		&item.ProfileID,
		&item.ProductID,
		&item.AddedAt,
		&item.IsActive,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrNotFound
		}
		return nil, fmt.Errorf("exec GetItemByID query: %w", err)
	}

	return &item, nil
}

func (r *WishlistRepo) GetItemsByProfileID(ctx context.Context, profileID string) ([]*models.WishlistItem, error) {
	stmt, args, err := r.goqu.Select(
		"item_id",
		"profile_id",
		"product_id",
		"added_at",
		"is_active",
	).From("wishlist_item").
		Where(goqu.Ex{"profile_id": profileID}).
		Order(goqu.C("added_at").Desc()).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build GetItemsByProfileID query: %w", err)
	}

	rows, err := r.DB.Query(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("exec GetItemsByProfileID query: %w", err)
	}
	defer rows.Close()

	var items []*models.WishlistItem
	for rows.Next() {
		var item models.WishlistItem
		if err := rows.Scan(
			&item.ItemID,
			&item.ProfileID,
			&item.ProductID,
			&item.AddedAt,
			&item.IsActive,
		); err != nil {
			return nil, fmt.Errorf("scan GetItemsByProfileID row: %w", err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate GetItemsByProfileID rows: %w", err)
	}

	return items, nil
}

func (r *WishlistRepo) RemoveItem(ctx context.Context, itemID string) error {
	stmt, args, err := r.goqu.Delete("wishlist_item").
		Where(goqu.Ex{"item_id": itemID}).
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build RemoveItem query: %w", err)
	}

	_, err = r.DB.Exec(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("exec RemoveItem query: %w", err)
	}

	return nil
}

func (r *WishlistRepo) UpdateItem(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error) {
	stmt, args, err := r.goqu.Update("wishlist_item").
		Set(goqu.Record{
			"is_active": item.IsActive,
		}).
		Where(goqu.Ex{"item_id": item.ItemID}).
		Returning(
			"item_id",
			"profile_id",
			"product_id",
			"added_at",
			"is_active",
		).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build UpdateItem query: %w", err)
	}

	var updatedItem models.WishlistItem
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(
		&updatedItem.ItemID,
		&updatedItem.ProfileID,
		&updatedItem.ProductID,
		&updatedItem.AddedAt,
		&updatedItem.IsActive,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrNotFound
		}
		return nil, fmt.Errorf("exec UpdateItem query: %w", err)
	}

	return &updatedItem, nil
}

func (r *WishlistRepo) Exists(ctx context.Context, profileID string, productID string) (bool, error) {
	stmt, args, err := r.goqu.Select(goqu.L("1")).
		From("wishlist_item").
		Where(goqu.Ex{
			"profile_id": profileID,
			"product_id": productID,
		}).
		Limit(1).
		Prepared(true).ToSQL()
	if err != nil {
		return false, fmt.Errorf("build Exists query: %w", err)
	}

	var exists int
	err = r.DB.QueryRow(ctx, stmt, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exec Exists query: %w", err)
	}

	return true, nil
}
