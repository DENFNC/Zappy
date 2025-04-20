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

type ShippingRepo struct {
	*psql.Storage
	goqu *goqu.DialectWrapper
}

func NewShippingRepo(
	db *psql.Storage,
	goqu *goqu.DialectWrapper,
) *ShippingRepo {
	return &ShippingRepo{
		Storage: db,
		goqu:    goqu,
	}
}

func (r *ShippingRepo) Create(ctx context.Context, address *models.Shipping) (uint32, error) {
	stmt, args, err := r.goqu.Insert("shipping_address").
		Rows(goqu.Record{
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
		return 0, err
	}

	var addrID uint32
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(&addrID); err != nil {
		return 0, err
	}

	return addrID, nil
}

func (r *ShippingRepo) GetByID(ctx context.Context, id uint32) (*models.Shipping, error) {
	stmt, args, err := r.goqu.Select(
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
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(
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

func (r *ShippingRepo) GetByProfileID(ctx context.Context, profileID uint32) ([]models.Shipping, error) {
	stmt, args, err := r.goqu.
		From("shipping_address").
		Where(goqu.Ex{"profile_id": profileID}).
		Order(goqu.C("is_default").Desc()).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build query GetByProfileID: %w", err)
	}

	rows, err := r.DB.Query(ctx, stmt, args...)
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

func (r *ShippingRepo) UpdateAddress(ctx context.Context, id uint32, address *models.Shipping) (uint32, error) {
	stmt, args, err := r.goqu.
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
		return 0, err
	}

	var updatedID uint32
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errpkg.ErrNotFound
		}
		return 0, err
	}

	return updatedID, nil
}

func (r *ShippingRepo) SetDefault(ctx context.Context, addressID, profileID uint32) error {
	conn, err := r.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return r.Storage.WithTx(ctx, conn, func(tx pgx.Tx) error {
		unsetSQL, unsetArgs, err := r.goqu.Update("shipping_address").
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

		setSQL, setArgs, err := r.goqu.Update("shipping_address").
			Set(goqu.Record{"is_default": true}).
			Returning(goqu.C("address_id")).
			Where(goqu.Ex{"address_id": addressID, "profile_id": profileID}).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}

		var id uint32
		if err := tx.QueryRow(ctx, setSQL, setArgs...).Scan(&id); err != nil {
			return err
		}

		return nil
	})
}

func (r *ShippingRepo) Delete(ctx context.Context, id uint32) (uint32, error) {
	stmt, args, err := r.goqu.Delete("shipping_address").
		Returning(goqu.C("address_id")).
		Where(goqu.C("address_id").Eq(id)).
		Prepared(true).ToSQL()
	if err != nil {
		return 0, err
	}

	var addrID uint32
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(&addrID); err != nil {
		return 0, err
	}

	return addrID, nil
}
