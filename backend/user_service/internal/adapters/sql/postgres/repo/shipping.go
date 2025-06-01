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

type ShippingRepo struct {
	*postgres.Storage
	*paginate.Paginator[dao.ShippingDAO]
}

func NewShippingRepo(
	db *postgres.Storage,
	coder paginate.TokenCoder,
) *ShippingRepo {
	paginate, err := paginate.NewPaginator[dao.ShippingDAO](
		db.Client, db.Dialect, coder,
	)
	if err != nil {
		panic(err)
	}

	return &ShippingRepo{
		Storage:   db,
		Paginator: paginate,
	}
}

func (r *ShippingRepo) Create(ctx context.Context, address *models.Shipping) (string, error) {
	uid, _ := uuid.NewV7()
	stmt, args, err := r.Dialect.Insert("shipping_address").
		Rows(goqu.Record{
			"address_id":  uid.String(),
			"profile_id":  address.ProfileID,
			"country":     address.Country,
			"city":        address.City,
			"street":      address.Street,
			"postal_code": address.PostalCode,
			"is_default":  address.IsDefault,
		}).
		Returning("address_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	var addrID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&addrID); err != nil {
		return "", err
	}

	return addrID, nil
}

func (r *ShippingRepo) GetByID(ctx context.Context, id string) (*models.Shipping, error) {
	stmt, args, err := r.Dialect.Select(
		"address_id",
		"profile_id",
		"country",
		"city",
		"street",
		"postal_code",
		"is_default",
	).From("shipping_address").
		Where(goqu.Ex{
			"address_id": id,
		}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	fmt.Println(stmt)

	var s models.Shipping
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&s.AddressID,
		&s.ProfileID,
		&s.Country,
		&s.City,
		&s.Street,
		&s.PostalCode,
		&s.IsDefault,
	); err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *ShippingRepo) GetByProfileID(ctx context.Context, profileID string) ([]models.Shipping, error) {
	stmt, args, err := r.Dialect.
		From("shipping_address").
		Where(goqu.Ex{"profile_id": profileID}).
		Order(goqu.C("is_default").Desc()).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build query GetByProfileID: %w", err)
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query GetByProfileID: %w", err)
	}
	defer rows.Close()

	var list []models.Shipping
	for rows.Next() {
		var ship models.Shipping
		if err := rows.Scan(
			&ship.AddressID,
			&ship.ProfileID,
			&ship.Country,
			&ship.City,
			&ship.Street,
			&ship.PostalCode,
			&ship.IsDefault,
		); err != nil {
			return nil, fmt.Errorf("scan row GetByProfileID: %w", err)
		}
		list = append(list, ship)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration GetByProfileID: %w", err)
	}

	return list, nil
}

func (r *ShippingRepo) UpdateAddress(ctx context.Context, id string, address *models.Shipping) (string, error) {
	stmt, args, err := r.Dialect.
		Update("shipping_address").
		Set(goqu.Record{
			"country":     address.Country,
			"city":        address.City,
			"street":      address.Street,
			"postal_code": address.PostalCode,
		}).
		Where(goqu.Ex{
			"address_id": id,
			"profile_id": address.ProfileID,
		}).
		Returning("address_id").
		Prepared(true).
		ToSQL()
	if err != nil {
		return "", err
	}

	var updatedID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", err
	}

	return updatedID, nil
}

func (r *ShippingRepo) SetDefault(ctx context.Context, addressID, profileID string) error {
	return r.WithTx(ctx, func(tx pgx.Tx) error {
		unsetSQL, unsetArgs, err := r.Dialect.Update("shipping_address").
			Set(goqu.Record{"is_default": false}).
			Where(goqu.Ex{"profile_id": profileID}).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, unsetSQL, unsetArgs...); err != nil {
			return err
		}

		setSQL, setArgs, err := r.Dialect.Update("shipping_address").
			Set(goqu.Record{"is_default": true}).
			Returning(goqu.C("address_id")).
			Where(goqu.Ex{"address_id": addressID, "profile_id": profileID}).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}

		var id string
		if err := tx.QueryRow(ctx, setSQL, setArgs...).Scan(&id); err != nil {
			return err
		}

		return nil
	})
}

func (r *ShippingRepo) Delete(ctx context.Context, id string) (string, error) {
	stmt, args, err := r.Dialect.Delete("shipping_address").
		Where(goqu.Ex{
			"address_id": id,
		}).
		Returning("address_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", errpkg.New("SHIPPING_DELETE_SQL_BUILD_ERROR", "failed to build delete shipping sql query", err)
	}

	var addrID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&addrID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", errpkg.New("SHIPPING_DELETE_QUERY_ERROR", "failed to execute delete shipping query", err)
	}

	return addrID, nil
}

func (repo *ShippingRepo) List(
	ctx context.Context,
	profileID string,
	pageSize uint,
	pageToken string,
) ([]*models.Shipping, string, error) {
	ds := repo.Dialect.Select(
		"address_id",
		"profile_id",
		"country",
		"city",
		"street",
		"postal_code",
		"is_default",
		"created_at",
		"updated_at",
	).From("shipping_address").
		Where(goqu.Ex{"profile_id": profileID})

	repo.Paginator.WithDataset(ds).WithColumns("created_at", "address_id").WithLimit(pageSize)

	itemsDAO, nextPageToken, err := repo.Paginator.Paginate(
		ctx,
		pageToken,
	)
	if err != nil {
		return nil, "", err
	}

	items := make([]*models.Shipping, len(itemsDAO))
	for i, itemDAO := range itemsDAO {
		items[i] = &models.Shipping{
			AddressID:  itemDAO.AddressID.String(),
			ProfileID:  itemDAO.ProfileID.String(),
			Country:    itemDAO.Country.String,
			City:       itemDAO.City.String,
			Street:     itemDAO.Street.String,
			PostalCode: itemDAO.PostalCode.String,
			IsDefault:  itemDAO.IsDefault.Bool,
			// CreatedAt:  itemDAO.CreatedAt.Time,
			// UpdatedAt:  itemDAO.UpdatedAt.Time,
		}
	}

	return items, nextPageToken, nil
}
