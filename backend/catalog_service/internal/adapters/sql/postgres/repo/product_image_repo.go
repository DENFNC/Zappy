package repo

import (
	"context"
	"errors"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/dbutils"
	"github.com/doug-martin/goqu/v9"
)

type ProductImage struct {
	*postgres.Storage
	paginate *paginate.Paginator[dao.ProductImageDAO]
}

func NewProductImage(
	db *postgres.Storage,
	coder paginate.TokenCoder,
) *ProductImage {
	paginate, err := paginate.NewPaginator[dao.ProductImageDAO](db.Client, db.Dialect, coder)
	if err != nil {
		panic(err)
	}

	return &ProductImage{
		Storage:  db,
		paginate: paginate,
	}
}

func (repo *ProductImage) Create(
	ctx context.Context,
	image *models.ProductImage,
) (string, error) {
	uid := dbutils.NewUUIDV7().String()

	stmt, args, err := repo.Dialect.Insert("product_image").Rows(
		goqu.Record{
			"image_id":   uid,
			"product_id": goqu.L("NULLIF(?, '')::uuid", image.ProductID),
			"url":        image.URL,
			"alt":        image.ALT,
			"object_key": image.ObjectKey,
			"created_at": time.Now().UTC(),
			"updated_at": time.Now().UTC(),
		}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return "", err
	}

	cmdTags, err := repo.Client.Exec(ctx, stmt, args...)
	if err != nil {
		return "", err
	}
	if cmdTags.RowsAffected() == 0 {
		return "", errors.New("no rows affected")
	}

	return uid, nil
}

func (repo *ProductImage) GetByID(
	ctx context.Context,
	uid string,
) (*models.ProductImage, error) {
	stmt, args, err := repo.Dialect.Select(
		"image_id",
		"product_id",
		"url",
		"alt",
		"object_key",
		"created_at",
		"updated_at",
	).
		From("product_image").
		Where(goqu.C("image_id").Eq(uid)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}

	row := repo.Client.QueryRow(ctx, stmt, args...)

	var productImage dao.ProductImageDAO
	if err := dbutils.ScanStruct(row, &productImage); err != nil {
		return nil, err
	}

	product := &models.ProductImage{
		ImageID:   productImage.ImageID.String(),
		ProductID: productImage.ProductID.String(),
		URL:       productImage.URL,
		ALT:       productImage.ALT,
		ObjectKey: productImage.ObjectKey,
		CreatedAt: productImage.CreatedAt.Time,
		UpdatedAt: productImage.UpdatedAt.Time,
	}

	return product, nil
}

func (repo *ProductImage) List(
	ctx context.Context,
	pageSize uint32,
	pageToken string,
	productID string,
) ([]models.ProductImage, string, error) {
	ds := repo.Dialect.Select(
		"image_id",
		"product_id",
		"url",
		"alt",
		"object_key",
		"created_at",
		"updated_at",
	).
		From("product_image").
		Where(goqu.C("product_id").Eq(productID)).
		Prepared(true)

	repo.paginate.WithDataset(ds).
		WithColumns("created_at", "image_id").
		WithLimit(uint(pageSize))

	itemsDAO, nextPageToken, err := repo.paginate.Paginate(
		ctx,
		pageToken,
	)
	if err != nil {
		return nil, "", err
	}

	productImages := make([]models.ProductImage, len(itemsDAO))
	for i, itemDAO := range itemsDAO {
		productImages[i] = models.ProductImage{
			ImageID:   itemDAO.ImageID.String(),
			ProductID: itemDAO.ProductID.String(),
			URL:       itemDAO.URL,
			ALT:       itemDAO.ALT,
			ObjectKey: itemDAO.ObjectKey,
			CreatedAt: itemDAO.CreatedAt.Time,
			UpdatedAt: itemDAO.UpdatedAt.Time,
		}
	}

	return productImages, nextPageToken, nil
}

func (repo *ProductImage) Delete(
	ctx context.Context,
	imageID string,
) (string, error) {
	stmt, args, err := repo.Dialect.Delete("product_image").
		Returning("object_key").
		Where(goqu.C("image_id").Eq(imageID)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return "", err
	}

	var objectKey string
	if err := repo.Client.QueryRow(ctx, stmt, args...).Scan(&objectKey); err != nil {
		return "", err
	}

	return objectKey, nil
}
