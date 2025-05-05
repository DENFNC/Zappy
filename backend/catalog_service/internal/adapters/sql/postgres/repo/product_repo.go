package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/dbutils"
	errpkg "github.com/DENFNC/Zappy/catalog_service/internal/utils/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
)

type ProductRepo struct {
	*postgres.Storage
	goqu     goqu.DialectWrapper
	paginate *paginate.Paginator[dao.ProductDAO]
}

func NewProductRepo(
	db *postgres.Storage,
	goqu goqu.DialectWrapper,
	coder paginate.TokenCoder,
) *ProductRepo {
	paginate, err := paginate.NewPaginator[dao.ProductDAO](db.DB, goqu, coder)
	if err != nil {
		panic(err)
	}

	return &ProductRepo{
		Storage:  db,
		goqu:     goqu,
		paginate: paginate,
	}
}

func (repo *ProductRepo) Create(
	ctx context.Context,
	name, desc string,
	categoryIDs []string,
	price int64,
) (string, error) {
	conn, err := repo.DB.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	productID := dbutils.NewUUIDV7().String()
	dateNow := time.Now().UTC()

	if err := dbutils.WithTx(ctx, conn, func(tx pgx.Tx) error {
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

	var productDAO dao.ProductDAO
	row := repo.DB.QueryRow(ctx, stmt, args...)
	if err := dbutils.ScanStruct(row, &productDAO); err != nil {
		return nil, err
	}

	product := models.Product{
		ProductID:   productDAO.ProductID.String(),
		ProductName: productDAO.ProductName,
		Description: productDAO.Description,
		Price:       productDAO.Price.Int.Int64(),
		CreatedAt:   productDAO.CreatedAt.Time,
		UpdatedAt:   productDAO.UpdatedAt.Time,
	}

	return &product, nil
}

func (repo *ProductRepo) List(
	ctx context.Context,
	pageSize uint32,
	pageToken string,
) ([]models.Product, string, error) {
	ds := repo.goqu.Select(
		"product_id",
		"product_name",
		"description",
		"price",
		"created_at",
		"updated_at",
	).
		From("product").
		Prepared(true)

	repo.paginate.WithDataset(ds).
		WithColumns("created_at", "product_id").
		WithLimit(uint(pageSize))

	itemsDAO, nextPageToken, err := repo.paginate.Paginate(
		ctx,
		pageToken,
	)
	if err != nil {
		return nil, "", err
	}

	products := make([]models.Product, len(itemsDAO))
	for i, itemDAO := range itemsDAO {
		products[i] = models.Product{
			ProductID:   itemDAO.ProductID.String(),
			ProductName: itemDAO.ProductName,
			Description: itemDAO.Description,
			Price:       itemDAO.Price.Int.Int64(),
			CreatedAt:   itemDAO.CreatedAt.Time,
			UpdatedAt:   itemDAO.UpdatedAt.Time,
		}
	}

	return products, nextPageToken, nil
}

func (repo *ProductRepo) Update(
	ctx context.Context,
	uid string,
	desc, name string,
	categoryIDs []string,
	price int64,
) error {
	conn, err := repo.DB.Acquire(ctx)
	if err != nil {
		return errors.New("acquire conn err")
	}
	defer conn.Release()

	now := time.Now().UTC()
	return dbutils.WithTx(ctx, conn, func(tx pgx.Tx) error {
		updSQL, updArgs, err := repo.goqu.Update("product").
			Set(goqu.Record{
				"product_name": name,
				"description":  desc,
				"price":        price,
				"updated_at":   now,
			}).
			Where(goqu.C("product_id").Eq(uid)).
			Prepared(true).
			ToSQL()
		if err != nil {
			return errors.New("build update product SQL error")
		}

		tag, err := tx.Exec(ctx, updSQL, updArgs...)
		if err != nil {
			return errors.New("exec update product error")
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("no product updated for id=%s", uid)
		}

		if len(categoryIDs) > 0 {
			delSQL, delArgs, err := repo.goqu.Delete("product_category").
				Where(goqu.C("product_id").Eq(uid)).
				Prepared(true).
				ToSQL()
			if err != nil {
				return errors.New("build delete categories SQL error")
			}

			if _, err := tx.Exec(ctx, delSQL, delArgs...); err != nil {
				return errors.New("exec delete categories error")
			}

			recs := make([]interface{}, len(categoryIDs))
			for i, cid := range categoryIDs {
				recs[i] = goqu.Record{
					"product_id":  uid,
					"category_id": cid,
					"assigned_at": now,
				}
			}

			insSQL, insArgs, err := repo.goqu.Insert("product_category").
				Rows(recs...).
				Prepared(true).
				ToSQL()
			if err != nil {
				return errors.New("build insert categories SQL")
			}

			tag, err := tx.Exec(ctx, insSQL, insArgs...)
			if err != nil {
				return errors.New("exec insert categories")
			}
			if tag.RowsAffected() == 0 {
				return fmt.Errorf("no categories inserted for id=%s", uid)
			}
		}

		return nil
	})
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
