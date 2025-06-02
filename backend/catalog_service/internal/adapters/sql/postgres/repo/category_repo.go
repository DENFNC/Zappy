package repo

import (
	"context"
	"errors"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
	"github.com/DENFNC/Zappy/catalog_service/internal/utils/dbutils"
	errpkg "github.com/DENFNC/Zappy/catalog_service/internal/utils/errors"

	"github.com/DENFNC/Zappy/catalog_service/internal/pkg/paginate"
	"github.com/doug-martin/goqu/v9"
)

type Category struct {
	*postgres.Storage
	paginate *paginate.Paginator[dao.CategoryDAO]
}

func NewCategoryRepo(
	db *postgres.Storage,
	coder paginate.TokenCoder,
) *Category {
	paginate, err := paginate.NewPaginator[dao.CategoryDAO](db.Client, db.Dialect, coder)
	if err != nil {
		panic(err)
	}

	return &Category{
		Storage:  db,
		paginate: paginate,
	}
}

func (repo *Category) Create(
	ctx context.Context,
	name, parentID string,
) (string, error) {
	uid := dbutils.NewUUIDV7().String()
	dateTimeNow := time.Now().UTC()

	stmt, args, err := repo.Dialect.Insert("category").
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

	cmdTags, err := repo.Client.Exec(ctx, stmt, args...)
	if err != nil {
		return "", err
	}
	if cmdTags.RowsAffected() == 0 {
		return "", errors.New("no rows affected")
	}

	return uid, nil
}

func (repo *Category) GetByID(
	ctx context.Context,
	categoryID string,
) (*models.Category, error) {
	panic("implement me!")
}

func (repo *Category) List(
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

func (repo *Category) ListByParentID(
	ctx context.Context,
	pageSize uint,
	parentID, pageToken string,
) (*models.Category, any, error) {
	panic("implement me!")
}

func (repo *Category) Delete(
	ctx context.Context,
	categoryID string,
) error {
	stmt, args, err := repo.Dialect.Delete("category").
		Where(goqu.C("category_id").Eq(categoryID)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	cmdTags, err := repo.Client.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}
	if cmdTags.RowsAffected() == 0 {
		return errpkg.New("NO_ROWS", "no rows affected", err)
	}

	return nil
}
