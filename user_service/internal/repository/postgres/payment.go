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

func NewPaymentMethodRepo(
	db *psql.Storage,
	goqu *goqu.DialectWrapper,
) *PaymentRepo {
	return &PaymentRepo{
		db,
		goqu,
	}
}

func (r *PaymentRepo) Create(ctx context.Context, method *models.Payment) (uint32, error) {
	panic("implement me!")
}

func (r *PaymentRepo) GetByID(ctx context.Context, id uint32) (*models.Payment, error) {
	panic("implement me!")
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
