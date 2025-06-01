package repo

import (
	"context"
	"errors"

	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres"
	"github.com/DENFNC/Zappy/user_service/internal/adapters/sql/postgres/dao"
	"github.com/DENFNC/Zappy/user_service/internal/domain/models"
	"github.com/DENFNC/Zappy/user_service/internal/pkg/paginate"
	errpkg "github.com/DENFNC/Zappy/user_service/internal/utils/errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type PaymentRepo struct {
	*postgres.Storage
	*paginate.Paginator[dao.PaymentDAO]
}

func NewPaymentRepo(
	db *postgres.Storage,
	coder paginate.TokenCoder,
) *PaymentRepo {
	paginate, err := paginate.NewPaginator[dao.PaymentDAO](
		db.Client, db.Dialect, coder,
	)
	if err != nil {
		panic(err)
	}

	return &PaymentRepo{
		Storage:   db,
		Paginator: paginate,
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
		return "", errpkg.New("PAYMENT_CREATE_SQL_BUILD_ERROR", "failed to build create payment sql query", err)
	}

	var paymentID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&paymentID); err != nil {
		return "", errpkg.New("PAYMENT_CREATE_QUERY_ERROR", "failed to execute create payment query", err)
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
		return nil, errpkg.New("PAYMENT_GETBYID_SQL_BUILD_ERROR", "failed to build get payment by id sql query", err)
	}

	var p models.Payment
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(
		&p.PaymentID,
		&p.ProfileID,
		&p.PaymentToken,
		&p.IsDefault,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errpkg.ErrNotFound
		}
		return nil, errpkg.New("PAYMENT_GETBYID_QUERY_ERROR", "failed to execute get payment by id query", err)
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
		return nil, errpkg.New("PAYMENT_GETBYPROFILEID_SQL_BUILD_ERROR", "failed to build get payment by profile id sql query", err)
	}

	rows, err := r.Client.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errpkg.New("PAYMENT_GETBYPROFILEID_QUERY_ERROR", "failed to execute get payment by profile id query", err)
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
			return nil, errpkg.New("PAYMENT_GETBYPROFILEID_SCAN_ERROR", "failed to scan payment row", err)
		}
		list = append(list, payment)
	}
	if err := rows.Err(); err != nil {
		return nil, errpkg.New("PAYMENT_GETBYPROFILEID_ROWS_ITERATION_ERROR", "failed to iterate payment rows", err)
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
			return errpkg.New("PAYMENT_SETDEFAULT_UNSET_SQL_BUILD_ERROR", "failed to build unset default payment sql query", err)
		}
		if _, err := tx.Exec(ctx, unsetSQL, unsetArgs...); err != nil {
			return errpkg.New("PAYMENT_SETDEFAULT_UNSET_QUERY_ERROR", "failed to execute unset default payment query", err)
		}

		setSQL, setArgs, err := r.Dialect.Update("payments").
			Set(goqu.Record{"is_default": true}).
			Returning(goqu.C("payment_id")).
			Where(goqu.Ex{"payment_id": paymentID, "profile_id": profileID}).
			Prepared(true).
			ToSQL()
		if err != nil {
			return errpkg.New("PAYMENT_SETDEFAULT_SET_SQL_BUILD_ERROR", "failed to build set default payment sql query", err)
		}

		var id string
		if err := tx.QueryRow(ctx, setSQL, setArgs...).Scan(&id); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errpkg.ErrNotFound
			}
			return errpkg.New("PAYMENT_SETDEFAULT_SET_QUERY_ERROR", "failed to execute set default payment query", err)
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
		return "", errpkg.New("PAYMENT_DELETE_SQL_BUILD_ERROR", "failed to build delete payment sql query", err)
	}

	var paymentID string
	if err := r.Client.QueryRow(ctx, stmt, args...).Scan(&paymentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errpkg.ErrNotFound
		}
		return "", errpkg.New("PAYMENT_DELETE_QUERY_ERROR", "failed to execute delete payment query", err)
	}

	return paymentID, nil
}

func (repo *PaymentRepo) List(
	ctx context.Context,
	profileID string,
	pageSize uint,
	pageToken string,
) ([]*models.Payment, string, error) {
	ds := repo.Dialect.Select(
		"payment_id",
		"profile_id",
		"payment_token",
		"is_default",
		"created_at",
		"updated_at",
	).
		From("payments").
		Where(goqu.C("profile_id").Eq(profileID))

	repo.Paginator.WithDataset(ds).WithColumns("created_at", "payment_id").WithLimit(pageSize)

	itemsDAO, nextPageToken, err := repo.Paginator.Paginate(
		ctx,
		pageToken,
	)
	if err != nil {
		return nil, "", errpkg.New("PAYMENT_LIST_PAGINATE_ERROR", "failed to paginate payments", err)
	}

	items := make([]*models.Payment, len(itemsDAO))
	for i, itemDAO := range itemsDAO {
		items[i] = &models.Payment{
			PaymentID:    itemDAO.PaymentID.String(),
			ProfileID:    itemDAO.ProfileID.String(),
			PaymentToken: itemDAO.PaymentToken.String,
			IsDefault:    itemDAO.IsDefault.Bool,
			// CreatedAt:  itemDAO.CreatedAt.Time,
			// UpdatedAt:  itemDAO.UpdatedAt.Time,
		}
	}

	return items, nextPageToken, nil
}
