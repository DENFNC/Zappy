package repo

import (
	"context"
	"errors"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/catalog_service/internal/errors"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/dbutils"
	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	psql "github.com/DENFNC/Zappy/catalog_service/internal/storage/postgres"
	"github.com/doug-martin/goqu/v9"
)

type CategoryRepo struct {
	*psql.Storage
	goqu     goqu.DialectWrapper
	paginate *paginate.Paginator[psql.CategoryDAO]
}

func NewCategoryRepo(
	db *psql.Storage,
	goqu goqu.DialectWrapper,
	coder paginate.TokenCoder,
) *CategoryRepo {
	paginate, err := paginate.NewPaginator[psql.CategoryDAO](
		db.DB,
		goqu,
		coder,
	)
	if err != nil {
		panic(err)
	}

	return &CategoryRepo{
		Storage:  db,
		goqu:     goqu,
		paginate: paginate,
	}
}

func (repo *CategoryRepo) Create(
	ctx context.Context,
	name, parentID string,
) (string, error) {
	uid := dbutils.NewUUIDV7().String()
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
	pageSize uint,
	pageToken string,
) ([]models.Category, string, error) {
	sqlStr := goqu.Select(
		"category_id",
		"category_name",
		"parent_id",
		"created_at",
		"updated_at",
	).From("category")

	paginate := repo.paginate.WithDataset(sqlStr).
		WithColumns("created_at", "category_id").
		WithLimit(pageSize)
	itemsDAO, nextToken, err := paginate.Paginate(ctx, pageToken)
	if err != nil {
		return nil, "", err
	}

	items := make([]models.Category, len(itemsDAO))
	for i, itemDAO := range itemsDAO {
		items[i] = models.Category{
			CategoryID:   itemDAO.CategoryID.String(),
			CategoryName: itemDAO.CategoryName,
			ParentID:     itemDAO.ParentID.String(),
			CreatedAt:    itemDAO.CreatedAt.Time,
			UpdatedAt:    itemDAO.UpdatedAt.Time,
		}
	}

	return items, nextToken, nil
}

func (repo *CategoryRepo) ListByParentID(
	ctx context.Context,
	pageSize uint,
	parentID, pageToken string,
) (*models.Category, any, error) {
	panic("implement me!")
}

func (repo *CategoryRepo) Delete(
	ctx context.Context,
	categoryID string,
) error {
	stmt, args, err := repo.goqu.Delete("category").
		Where(goqu.C("category_id").Eq(categoryID)).
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
