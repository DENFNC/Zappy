package postgres

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Client  *pgxpool.Pool
	Dialect goqu.DialectWrapper
}

func New(conn string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	dialect := goqu.Dialect("postgres")

	return &Storage{
		Client:  dbpool,
		Dialect: dialect,
	}, nil
}

func (s *Storage) Stop() {
	s.Client.Close()
}
