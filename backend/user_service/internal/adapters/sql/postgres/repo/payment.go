package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type PaymentRepo struct {
	*postgres.Storage
}

func NewPaymentRepo(
	db *postgres.Storage,
) *PaymentRepo {
	return &PaymentRepo{
		Storage: db,
	}
}

func (r *PaymentRepo) Create(ctx context.Context, payment *models.Payment) (string, error) {
	uid, _ := uuid.NewV7()
	stmt, args, err := r.Dialect.Insert("payments").
		Rows(goqu.Record{
			"payment_id":    uid.String(),
			"profile_id":    payment.ProfileID,
			"payment_token": payment.PaymentToken,
			"is_default":    payment.IsDefault,
		}).
		Returning("payment_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	var paymentID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&paymentID); err != nil {
		return "", err
	}

	return paymentID, nil
}

func (r *PaymentRepo) GetByID(ctx context.Context, id string) (*models.Payment, error) {
	stmt, args, err := r.Dialect.Select(
		"payment_id",
		"profile_id",
		"payment_token",
		"is_default",
		"create_at",
		"updated_at",
	).From("payments").
		Where(goqu.Ex{
			"payment_id": id,
		}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	var p models.Payment
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&p.PaymentID,
		&p.ProfileID,
		&p.PaymentToken,
		&p.IsDefault,
	); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PaymentRepo) GetByProfileID(ctx context.Context, profileID string) ([]models.Payment, error) {
	stmt, args, err := r.Dialect.
		From("payments").
		Where(goqu.Ex{"profile_id": profileID}).
		Order(goqu.C("is_default").Desc()).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build query GetByProfileID: %w", err)
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query GetByProfileID: %w", err)
	}
	defer rows.Close()

	var list []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(
			&payment.PaymentID,
			&payment.ProfileID,
			&payment.PaymentToken,
			&payment.IsDefault,
		); err != nil {
			return nil, fmt.Errorf("scan row GetByProfileID: %w", err)
		}
		list = append(list, payment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration GetByProfileID: %w", err)
	}

	return list, nil
}

func (r *PaymentRepo) Update(
	ctx context.Context,
	uid string,
	payment *models.Payment,
) (string, error) {
	// stmt, args, err := r.Dialect.
	// 	Update("payments").
	// 	Set(goqu.Record{
	//		"payment_token":
	// 	}).
	// 	Where(goqu.Ex{
	// 		"payment_id": id,
	// 		"profile_id": payment.ProfileID,
	// 	}).
	// 	Returning("payment_id").
	// 	Prepared(true).
	// 	ToSQL()
	// if err != nil {
	// 	return "", err
	// }

	// var updatedID string
	// if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&updatedID); err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return "", errpkg.ErrNotFound
	// 	}
	// 	return "", err
	// }

	return "", nil
}

func (r *PaymentRepo) SetDefault(ctx context.Context, paymentID, profileID string) error {
	return r.WithTx(ctx, func(tx pgx.Tx) error {
		unsetSQL, unsetArgs, err := r.Dialect.Update("payments").
			Set(goqu.Record{"is_default": false}).
			Where(goqu.Ex{"profile_id": profileID}).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, unsetSQL, unsetArgs...); err != nil {
			return err
		}

		setSQL, setArgs, err := r.Dialect.Update("payments").
			Set(goqu.Record{"is_default": true}).
			Returning(goqu.C("payment_id")).
			Where(goqu.Ex{"payment_id": paymentID, "profile_id": profileID}).
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}

		var id string
		if err := tx.QueryRow(ctx, setSQL, setArgs...).Scan(&id); err != nil {
			return err
		}

		return nil
	})
}

func (r *PaymentRepo) Delete(ctx context.Context, id string) (string, error) {
	stmt, args, err := r.Dialect.Delete("payments").
		Where(goqu.Ex{
			"payment_id": id,
		}).
		Returning("payment_id").
		Prepared(true).ToSQL()
	if err != nil {
		return "", err
	}

	var paymentID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&paymentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", err
	}

	return paymentID, nil
}

func (r *PaymentRepo) List(ctx context.Context, profileID string) ([]models.Payment, error) {
	stmt, args, err := r.Dialect.Select(
		"payment_id",
		"profile_id",
		"payment_token",
		"is_default",
		"create_at",
		"updated_at",
	).From("payments").
		Where(goqu.Ex{"profile_id": profileID}).
		Order(goqu.C("is_default").Desc()).
		Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build query List: %w", err)
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query List: %w", err)
	}
	defer rows.Close()

	var list []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(
			&payment.PaymentID,
			&payment.ProfileID,
			&payment.PaymentToken,
			&payment.IsDefault,
		); err != nil {
			return nil, fmt.Errorf("scan row List: %w", err)
		}
		list = append(list, payment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration List: %w", err)
	}

	return list, nil
}
