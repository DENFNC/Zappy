package repo

import (
	"context"
	"errors"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	psql "github.com/DENFNC/Zappy/catalog_service/internal/storage/postgres"
	"github.com/doug-martin/goqu/v9"
)

type CategoryRepo struct {
	*psql.Storage
	goqu *goqu.DialectWrapper
}

func NewCategoryRepo(
	db *psql.Storage,
	goqu *goqu.DialectWrapper,
) *CategoryRepo {
	return &CategoryRepo{
		Storage: db,
		goqu:    goqu,
	}
}

func (repo *CategoryRepo) Create(
	ctx context.Context,
	name, parentID string,
) (string, error) {

	uid := repo.NewV7().String()
	dateTimeNow := time.Now().UTC()

	stmt, args, err := repo.goqu.Insert("category").
		Rows(goqu.Record{
			"category_id":   uid,
			"category_name": name,
			"parent_id":     goqu.L("NULLIF(?, '')::uuid", parentID),
			"created_at":    dateTimeNow,
			"updated_at":    dateTimeNow,
		}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return "", err
	}

	cmdTags, err := repo.DB.Exec(ctx, stmt, args...)
	if err != nil {
		return "", err
	}
	if cmdTags.RowsAffected() == 0 {
		return "", errors.New("no rows affected")
	}

	return uid, nil
}

func (repo *CategoryRepo) GetByID(
	ctx context.Context,
	categoryID string,
) (*models.Category, error) {
	panic("implement me!")
}

func (repo *CategoryRepo) List(
	ctx context.Context,
	pageSize int32,
	pageToken string,
) ([]models.Category, any, error) {
	panic("implement me!")
}

func (repo *CategoryRepo) ListByParentID(
	ctx context.Context,
	pageSize int32,
	parentID, pageToken string,
) (*models.Category, any, error) {
	panic("implement me!")
}

func (repo *CategoryRepo) Delete(
	ctx context.Context,
	categoryID string,
) error {
	panic("implement me!")
}
