package psql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func New(conn string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, err
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Storage{
		DB: dbpool,
	}, nil
}
