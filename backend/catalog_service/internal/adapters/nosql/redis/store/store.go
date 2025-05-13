package store

import (
	"context"
	"time"

	rds "github.com/DENFNC/Zappy/catalog_service/internal/adapters/nosql/redis"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	*rds.Client
}

func New(
	client *rds.Client,
) *Store {
	return &Store{
		Client: client,
	}
}

func (store *Store) Get(
	ctx context.Context,
	key string,
) (string, error) {
	cmdTag := store.Client.Get(ctx, key)
	if cmdTag.Err() != nil {
		return "", cmdTag.Err()
	}

	return cmdTag.Result()
}

func (store *Store) Set(
	ctx context.Context,
	key, value string,
) error {
	return store.Client.Set(
		ctx,
		key,
		value,
		redis.KeepTTL,
	).Err()
}

func (store *Store) SetEx(
	ctx context.Context,
	key string,
	value any,
	expire time.Duration,
) error {
	return store.Client.SetEx(
		ctx,
		key,
		value,
		expire,
	).Err()
}
