package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/review_service/internal/domain/models"
	"github.com/DENFNC/Zappy/review_service/internal/pkg/paginate"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofrs/uuid"
)

type Review struct {
	*postgres.Storage
	*paginate.Paginator[dao.ReviewDAO]
}

func NewReviewRepo(
	db *postgres.Storage,
	coder paginate.TokenCoder,
) *Review {
	paginate, err := paginate.NewPaginator[dao.ReviewDAO](
		db.Client, db.Dialect, coder,
	)
	if err != nil {
		panic(err)
	}

	return &Review{
		Storage:   db,
		Paginator: paginate,
	}
}

func (repo *Review) Create(
	ctx context.Context,
	review *models.Review,
) (string, error) {
	uid, _ := uuid.NewV7()
	dateTimeNow := time.Now().UTC()

	stmt, args, err := repo.Dialect.Insert("review").
		Rows(goqu.Record{
			"review_id":  uid.String(),
			"product_id": review.ProductID,
			"profile_id": review.ProfileID,
			"rating":     review.Rating,
			"created_at": dateTimeNow,
			"updated_at": dateTimeNow,
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

	return uid.String(), nil
}
