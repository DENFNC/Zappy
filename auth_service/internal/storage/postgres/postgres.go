package psql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func New(conn string) *Storage {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, conn)
	if err != nil {
		panic("couldn't connect to the database")
	}

	if err := dbpool.Ping(ctx); err != nil {
		panic("couldn't ping the database")
	}

	return &Storage{
		DB: dbpool,
	}
}

func (s *Storage) WithTx(ctx context.Context, conn *pgxpool.Conn, f func(pgx.Tx) error) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = f(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
