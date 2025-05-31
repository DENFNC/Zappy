package postgres

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Client  *pgxpool.Pool
	Dialect *goqu.DialectWrapper
}

func NewStorage(conn string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, err
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, err
	}

	dialect := goqu.Dialect("postgres")

	return &Storage{
		Client:  dbpool,
		Dialect: &dialect,
	}, nil
}
func (s *Storage) WithTx(ctx context.Context, f func(tx pgx.Tx) error) (err error) {
	tx, err := s.Client.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		} else if err != nil {
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
