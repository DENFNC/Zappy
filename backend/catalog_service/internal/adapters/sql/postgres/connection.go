package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Client  *pgxpool.Pool
	Dialect goqu.DialectWrapper
	log     *slog.Logger
}

func NewStorage(conn string, log *slog.Logger) (*Storage, error) {
	const op = "postgres.NewStorage"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, conn)
	if err != nil {
		log.Error(
			"Failed to create connection pool",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Error(
			"Failed to ping database",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	dialect := goqu.Dialect("postgres")

	return &Storage{
			Client:  dbpool,
			Dialect: dialect,
			log:     log,
		},
		nil
}
func (s *Storage) WithTx(ctx context.Context, f func(tx pgx.Tx) error) (err error) {
	const op = "postgres.WithTx"

	tx, err := s.Client.Begin(ctx)
	if err != nil {
		s.log.Error(
			"Failed to begin transaction",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			s.log.Error(
				"Transaction panicked, rolling back",
				slog.String("op", op),
				slog.Any("panic", r),
			)
			_ = tx.Rollback(ctx)
			panic(r)
		} else if err != nil {
			s.log.Error(
				"Transaction failed, rolling back",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			_ = tx.Rollback(ctx)
		}
	}()

	if err = f(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (s *Storage) Stop() {
	s.Client.Close()
}
