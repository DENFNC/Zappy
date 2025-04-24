package repo

import (
	"context"

	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	psql "github.com/DENFNC/Zappy/user_service/internal/storage/postgres"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UUID = uuid.UUID

type PaymentRepo struct {
	*psql.Storage
	goqu *goqu.DialectWrapper
}

func NewPaymentRepo(db *psql.Storage, g *goqu.DialectWrapper) *PaymentRepo {
	return &PaymentRepo{
		Storage: db,
		goqu:    g,
	}
}

func (r *PaymentRepo) Create(ctx context.Context, method *models.Payment) (UUID, error) {
	id := r.NewV7()
	stmt, args, err := r.goqu.Insert("payment_method").
		Returning(goqu.C("payment_id")).
		Rows(goqu.Record{
			"payment_id":    id,
			"profile_id":    method.ProfileID,
			"payment_token": method.PaymentToken,
			"is_default":    false,
		}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return uuid.Nil, err
	}

	var payID UUID
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(&payID); err != nil {
		return uuid.Nil, err
	}
	return payID, nil
}

func (r *PaymentRepo) GetByID(ctx context.Context, id uint32) (*models.Payment, error) {
	stmt, args, err := r.goqu.Select(
		"payment_id",
		"profile_id",
		"payment_token",
		"is_default",
	).
		From("payment_method").
		Where(goqu.C("payment_id").Eq(id)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var p models.Payment
	if err := r.DB.QueryRow(ctx, stmt, args...).Scan(
		&p.PaymentID,
		&p.ProfileID,
		&p.PaymentToken,
		&p.IsDefault,
	); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepo) GetByProfileID(ctx context.Context, profileID uint32) ([]models.Payment, error) {
	stmt, args, err := r.goqu.Select(
		"payment_id",
		"profile_id",
		"payment_token",
		"is_default",
	).
		From("payment_method").
		Where(goqu.C("profile_id").Eq(profileID)).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var m models.Payment
		if err := rows.Scan(
			&m.PaymentID,
			&m.ProfileID,
			&m.PaymentToken,
			&m.IsDefault,
		); err != nil {
			return nil, err
		}
		payments = append(payments, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepo) SetDefault(ctx context.Context, methodID, profileID uint32) (uint32, error) {
	conn, err := r.DB.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var payID uint32
	err = r.WithTx(ctx, conn, func(tx pgx.Tx) error {
		stmt, args, err := r.goqu.Update("payment_method").
			Set(goqu.Record{"is_default": false}).
			Where(goqu.Ex{"profile_id": profileID}).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, stmt, args...); err != nil {
			return err
		}

		stmt, args, err = r.goqu.Update("payment_method").
			Returning("payment_id").
			Set(goqu.Record{"is_default": true}).
			Where(goqu.Ex{"profile_id": profileID, "payment_id": methodID}).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}
		return tx.QueryRow(ctx, stmt, args...).Scan(&payID)
	})
	if err != nil {
		return 0, err
	}
	return payID, nil
}

func (r *PaymentRepo) Update(ctx context.Context, method *models.Payment) (uint32, error) {
	stmt, args, err := r.goqu.Update("payment_method").
		Returning("profile_id").
		Set(goqu.Record{"payment_token": method.PaymentToken}).
		Where(goqu.Ex{"profile_id": method.ProfileID, "payment_id": method.PaymentID}).
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

func (r *PaymentRepo) Delete(ctx context.Context, id uint32) (uint32, error) {
	stmt, args, err := r.goqu.Delete("payment_method").
		Returning("payment_id").
		Where(goqu.C("payment_id").Eq(id)).
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
