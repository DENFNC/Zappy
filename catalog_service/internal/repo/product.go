package repo

import (
	"context"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/catalog_service/internal/errors"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/dbutils"
	psql "github.com/DENFNC/Zappy/catalog_service/internal/storage/postgres"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductRepo struct {
	*psql.Storage
	goqu goqu.DialectWrapper
}

func NewProductRepo(
	db *psql.Storage,
	goqu goqu.DialectWrapper,
) *ProductRepo {
	return &ProductRepo{
		Storage: db,
		goqu:    goqu,
	}
}

func (repo *ProductRepo) Create(
	ctx context.Context,
	name, desc string,
	categoryIDs []string,
	price pgtype.Numeric,
) (string, error) {
	conn, err := repo.DB.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	productID := repo.NewV7().String()
	dateNow := time.Now().UTC()

	if err := repo.WithTx(ctx, conn, func(tx pgx.Tx) error {
		stmt, args, err := repo.goqu.Insert("product").
			Rows(
				goqu.Record{
					"product_id":   productID,
					"product_name": name,
					"description":  desc,
					"price":        price,
					"created_at":   dateNow,
					"updated_at":   dateNow,
				},
			).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}

		cmdTags, err := tx.Exec(ctx, stmt, args...)
		if err != nil {
			return err
		}
		if cmdTags.RowsAffected() == 0 {
			return errpkg.New("NO_ROWS", "no rows affected", err)
		}

		records := make([]interface{}, len(categoryIDs))
		for i, cid := range categoryIDs {
			records[i] = goqu.Record{
				"product_id":  productID,
				"category_id": cid,
				"assigned_at": dateNow,
			}
		}

		stmt, args, err = repo.goqu.Insert("product_category").
			Rows(records...).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}

		cmdTags, err = tx.Exec(ctx, stmt, args...)
		if err != nil {
			return err
		}
		if cmdTags.RowsAffected() == 0 {
			return errpkg.New("NO_ROWS", "no rows affected", err)
		}

		return nil
	}); err != nil {
		return "", err
	}

	return productID, nil
}

func (repo *ProductRepo) GetByID(
	ctx context.Context,
	productID string,
) (*models.Product, error) {
	stmt, args, err := repo.goqu.Select(
		"product_id",
		"product_name",
		"description",
		"price",
		"created_at",
		"updated_at",
	).
		From("product").
		Where(goqu.C("product_id").Eq(productID)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var product models.Product
	row := repo.DB.QueryRow(ctx, stmt, args...)
	if err := dbutils.ScanStruct(row, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (repo *ProductRepo) List(
	ctx context.Context,
) ([]models.Product, error) {
	panic("implement me")
}

func (repo *ProductRepo) Update(
	ctx context.Context,
	uid string,
	desc, name string,
	price pgtype.Numeric,
) error {
	panic("implement me")
}

func (repo *ProductRepo) Delete(
	ctx context.Context,
	productID string,
) error {
	stmt, args, err := repo.goqu.Delete("product").
		Where(goqu.C("product_id").Eq(productID)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	cmdTags, err := repo.DB.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}
	if cmdTags.RowsAffected() == 0 {
		return errpkg.New("NO_ROWS", "no rows affected", err)
	}

	return nil
}
