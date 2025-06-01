package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/paginate"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type WishlistRepo struct {
	*postgres.Storage
	*paginate.Paginator[dao.WishlistDAO]
}

func NewWishlistRepo(
	db *postgres.Storage,
	coder paginate.TokenCoder,
) *WishlistRepo {
	paginate, err := paginate.NewPaginator[dao.WishlistDAO](
		db.Client, db.Dialect, coder,
	)
	if err != nil {
		panic(err)
	}

	return &WishlistRepo{
		Storage:   db,
		Paginator: paginate,
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
		return "", errpkg.New("WISHLIST_CREATE_SQL_BUILD_ERROR", "failed to build create wishlist sql query", err)
	}

	var wishlistID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&wishlistID); err != nil {
		return "", errpkg.New("WISHLIST_CREATE_QUERY_ERROR", "failed to execute create wishlist query", err)
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
		return nil, errpkg.New("WISHLIST_GETBYID_SQL_BUILD_ERROR", "failed to build get wishlist by id sql query", err)
	}

	var wishList models.WishlistItem
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&wishList.ItemID,
		&wishList.ProfileID,
		&wishList.ProductID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrNotFound
		}
		return nil, errpkg.New("WISHLIST_GETBYID_QUERY_ERROR", "failed to execute get wishlist by id query", err)
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
		return nil, errpkg.New("WISHLIST_GETBYPROFILEID_SQL_BUILD_ERROR", "failed to build get wishlist by profile id sql query", err)
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errpkg.New("WISHLIST_GETBYPROFILEID_QUERY_ERROR", "failed to execute get wishlist by profile id query", err)
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
			return nil, errpkg.New("WISHLIST_GETBYPROFILEID_SCAN_ERROR", "failed to scan wishlist row", err)
		}
		wishlists = append(wishlists, wishlist)
	}
	if err := rows.Err(); err != nil {
		return nil, errpkg.New("WISHLIST_GETBYPROFILEID_ROWS_ITERATION_ERROR", "failed to iterate wishlist rows", err)
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
		return "", errpkg.New("WISHLIST_DELETE_SQL_BUILD_ERROR", "failed to build delete wishlist sql query", err)
	}

	var wishlistID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&wishlistID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", errpkg.New("WISHLIST_DELETE_QUERY_ERROR", "failed to execute delete wishlist query", err)
	}

	return wishlistID, nil
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
		return nil, errpkg.New("WISHLIST_GETITEMBYID_SQL_BUILD_ERROR", "failed to build get wishlist item by id sql query", err)
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
		return nil, errpkg.New("WISHLIST_GETITEMBYID_QUERY_ERROR", "failed to execute get wishlist item by id query", err)
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
		return nil, errpkg.New("WISHLIST_GETITEMSBYPROFILEID_SQL_BUILD_ERROR", "failed to build get wishlist items by profile id sql query", err)
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errpkg.New("WISHLIST_GETITEMSBYPROFILEID_QUERY_ERROR", "failed to execute get wishlist items by profile id query", err)
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
			return nil, errpkg.New("WISHLIST_GETITEMSBYPROFILEID_SCAN_ERROR", "failed to scan wishlist item row", err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, errpkg.New("WISHLIST_GETITEMSBYPROFILEID_ROWS_ITERATION_ERROR", "failed to iterate wishlist item rows", err)
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

func (repo *WishlistRepo) List(
	ctx context.Context,
	profileID string,
	pageSize uint,
	pageToken string,
) ([]*models.WishlistItem, string, error) {
	ds := repo.Dialect.Select(
		"item_id",
		"profile_id",
		"product_id",
		"added_at",
		"is_active",
	).From("wishlist_item").
		Where(goqu.Ex{"profile_id": profileID})

	repo.Paginator.WithDataset(ds).WithColumns("added_at", "item_id").WithLimit(pageSize)

	itemsDAO, nextPageToken, err := repo.Paginator.Paginate(
		ctx,
		pageToken,
	)
	if err != nil {
		return nil, "", errpkg.New("WISHLIST_LIST_PAGINATE_ERROR", "failed to paginate wishlist items", err)
	}

	items := make([]*models.WishlistItem, len(itemsDAO))
	for i, itemDAO := range itemsDAO {
		items[i] = &models.WishlistItem{
			ItemID:    itemDAO.ItemID.String(),
			ProfileID: itemDAO.ProfileID.String(),
			ProductID: itemDAO.ProductID.String(),
			AddedAt:   itemDAO.AddedAt.Time,
			IsActive:  itemDAO.IsActive.Bool,
		}
	}

	return items, nextPageToken, nil
}
