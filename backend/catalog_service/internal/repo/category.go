package repo

import (
	"context"
	"errors"
	"time"

	"github.com/DENFNC/Zappy/catalog_service/internal/domain/models"
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
) (*CategoryRepo, error) {
	paginate, err := paginate.NewPaginator[psql.CategoryDAO](
		db.DB,
		goqu,
		coder,
	)
	if err != nil {
		return nil, err
	}

	return &CategoryRepo{
		Storage:  db,
		goqu:     goqu,
		paginate: paginate,
	}, nil
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
	pageSize uint,
	pageToken string,
) ([]models.Category, string, error) {
	paginate := repo.paginate.WithTable("category").
		WithColumns("created_at", "category_id").
		WithLimit(pageSize)
	// TODO: пофиксить проблему с nil элементами
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
	panic("implement me!")
}
