package repo

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	psql "github.com/DENFNC/Zappy/user_service/internal/storage/postgres"
	"github.com/doug-martin/goqu/v9"
)

type PaymentRepo struct {
	*psql.Storage
	goqu *goqu.DialectWrapper
}

func NewPaymentRepo(
	db *psql.Storage,
	goqu *goqu.DialectWrapper,
) *PaymentRepo {
	return &PaymentRepo{
		db,
		goqu,
	}
}

func (r *PaymentRepo) Create(ctx context.Context, method *models.Payment) (uint32, error) {
	stmt, args, err := r.goqu.Insert("payment_method").
		Returning(goqu.C("payment_id")).
		Rows(
			goqu.Record{
				"profile_id":    method.ProfileID,
				"payment_token": method.PaymentToken,
				"is_default":    true,
			},
		).
		Prepared(true).
		ToSQL()
	if err != nil {
		return 0, err
	}

	var payID uint32
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(&payID); err != nil {
		return 0, err
	}

	return payID, nil
}

func (r *PaymentRepo) GetByID(ctx context.Context, id uint32) (*models.Payment, error) {
	stmt, args, err := r.goqu.Select(
		"payment_id",
		"profile_id",
		"payment_token",
		"is_default",
	).From("payment_method").Where(
		goqu.C("payment_id").Eq(id),
	).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	var payment models.Payment
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(
		&payment.PaymentID,
		&payment.ProfileID,
		&payment.PaymentToken,
		&payment.IsDefault,
	); err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepo) GetByProfileID(ctx context.Context, profileID uint32) ([]*models.Payment, error) {
	panic("implement me!")
}

func (r *PaymentRepo) SetDefault(ctx context.Context, methodID uint32, profileID uint32) (uint32, error) {
	panic("implement me!")
}

func (r *PaymentRepo) Update(ctx context.Context, method *models.Payment) (uint32, error) {
	panic("implement me!")
}

func (r *PaymentRepo) Delete(ctx context.Context, id uint32) (uint32, error) {
	panic("implement me!")
}
