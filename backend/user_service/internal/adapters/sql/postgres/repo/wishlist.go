package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type WishlistRepo struct {
	*postgres.Storage
}

func NewWishlistRepo(
	db *postgres.Storage,
) *WishlistRepo {
	return &WishlistRepo{
		Storage: db,
	}
}

func (r *WishlistRepo) Create(ctx context.Context, wishlist *models.WishlistItem) (string, error) {
	uid, _ := uuid.NewV7()
	stmt, args, err := r.Dialect.Insert("wishlists").
		Rows(goqu.Record{
			"wishlist_id": uid.String(),
			"profile_id":  wishlist.ProfileID,
			"product_id":  wishlist.ProductID,
		}).
		Returning("wishlist_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	var wishlistID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&wishlistID); err != nil {
		return "", err
	}

	return wishlistID, nil
}

func (r *WishlistRepo) GetByID(ctx context.Context, id string) (*models.WishlistItem, error) {
	stmt, args, err := r.Dialect.Select(
		"wishlist_id",
		"profile_id",
		"product_id",
	).From("wishlists").
		Where(goqu.Ex{
			"wishlist_id": id,
		}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	var wishList models.WishlistItem
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&wishList.ItemID,
		&wishList.ProfileID,
		&wishList.ProductID,
	); err != nil {
		return nil, err
	}

	return &wishList, nil
}

func (r *WishlistRepo) GetByProfileID(ctx context.Context, profileID string) ([]models.WishlistItem, error) {
	stmt, args, err := r.Dialect.Select(
		"wishlist_id",
		"profile_id",
		"product_id",
	).From("wishlists").
		Where(goqu.Ex{
			"profile_id": profileID,
		}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []models.WishlistItem
	for rows.Next() {
		var wishlist models.WishlistItem
		if err := rows.Scan(
			&wishlist.ItemID,
			&wishlist.ProfileID,
			&wishlist.ProductID,
		); err != nil {
			return nil, err
		}
		wishlists = append(wishlists, wishlist)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlists, nil
}

func (r *WishlistRepo) Delete(ctx context.Context, id string) (string, error) {
	stmt, args, err := r.Dialect.Delete("wishlists").
		Where(goqu.Ex{
			"wishlist_id": id,
		}).
		Returning("wishlist_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	var wishlistID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&wishlistID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", err
	}

	return wishlistID, nil
}

func (r *WishlistRepo) AddItem(ctx context.Context, item *models.WishlistItem) (string, error) {
	uid, _ := uuid.NewV7()

	stmt, args, err := r.Dialect.Insert("wishlist_item").
		Rows(goqu.Record{
			"item_id":    uid.String(),
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
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&itemID); err != nil {
		return "", fmt.Errorf("exec AddItem query: %w", err)
	}

	return itemID, nil
}

func (r *WishlistRepo) GetItemByID(ctx context.Context, itemID string) (*models.WishlistItem, error) {
	stmt, args, err := r.Dialect.Select(
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
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
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
	stmt, args, err := r.Dialect.Select(
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

	rows, err := r.Client.Query(ctx, stmt, args...)
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
	stmt, args, err := r.Dialect.Delete("wishlist_item").
		Where(goqu.Ex{"item_id": itemID}).
		Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("build RemoveItem query: %w", err)
	}

	_, err = r.Client.Exec(ctx, stmt, args...)
	if err != nil {
		return fmt.Errorf("exec RemoveItem query: %w", err)
	}

	return nil
}

func (r *WishlistRepo) UpdateItem(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error) {
	stmt, args, err := r.Dialect.Update("wishlist_item").
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
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
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
	stmt, args, err := r.Dialect.Select(goqu.L("1")).
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
	err = r.Client.QueryRow(ctx, stmt, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exec Exists query: %w", err)
	}

	return true, nil
}

func (r *WishlistRepo) List(ctx context.Context, params []any) ([]*models.WishlistItem, error) {
	stmt, args, err := r.Dialect.Select(
		"wishlist_id",
		"profile_id",
		"product_id",
	).From("wishlists").
		Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []*models.WishlistItem
	for rows.Next() {
		var wishlist models.WishlistItem
		if err := rows.Scan(
			&wishlist.ItemID,
			&wishlist.ProfileID,
			&wishlist.ProductID,
		); err != nil {
			return nil, err
		}
		wishlists = append(wishlists, &wishlist)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlists, nil
}
