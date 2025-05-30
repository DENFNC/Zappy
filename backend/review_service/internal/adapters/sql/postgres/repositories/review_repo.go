package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/review_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/review_service/internal/domain/models"
	"github.com/DENFNC/Zappy/review_service/internal/pkg/paginate"
	"github.com/DENFNC/Zappy/review_service/internal/utils/dbutils"
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
			"comment":    review.Comment,
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

func (repo *Review) GetByID(
	ctx context.Context,
	uid string,
) (*models.Review, error) {
	stmt, args, err := repo.Dialect.Select(
		"review_id",
		"product_id",
		"profile_id",
		"rating",
		"comment",
		"created_at",
		"updated_at",
	).
		From("review").
		Where(goqu.C("review_id").Eq(uid)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var reviewDAO dao.ReviewDAO
	row := repo.Client.QueryRow(ctx, stmt, args...)
	if err := dbutils.ScanStruct(row, &reviewDAO); err != nil {
		fmt.Println("Ошибка тут")
		return nil, err
	}

	return &models.Review{
		ReviewID:  reviewDAO.ReviewID.String(),
		ProductID: reviewDAO.ProductId.String(),
		ProfileID: reviewDAO.ProfileID.String(),
		Rating:    reviewDAO.Rating,
		Comment:   reviewDAO.Comment.String,
		CreatedAt: reviewDAO.CreatedAt.Time,
		UpdatedAt: reviewDAO.UpdatedAt.Time,
	}, nil
}
