package store

import (
	"context"
	"log/slog"
	"time"

	rds "github.com/DENFNC/Zappy/catalog_service/internal/adapters/nosql/redis"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	*rds.Client
	log *slog.Logger
}

func New(
	client *rds.Client,
	log *slog.Logger,
) *Store {
	return &Store{
		Client: client,
		log:    log,
	}
}

func (store *Store) Get(
	ctx context.Context,
	key string,
) (string, error) {
	const op = "redisStore.Get"

	cmdTag := store.Client.Get(ctx, key)
	if cmdTag.Err() != nil {
		store.log.Error(
			"Failed to get key from Redis",
			slog.String("op", op),
			slog.String("key", key),
			slog.String("error", cmdTag.Err().Error()),
		)
		return "", cmdTag.Err()
	}

	return cmdTag.Result()
}

func (store *Store) Set(
	ctx context.Context,
	key, value string,
) error {
	const op = "redisStore.Set"

	err := store.Client.Set(
		ctx,
		key,
		value,
		redis.KeepTTL,
	).Err()
	if err != nil {
		store.log.Error(
			"Failed to set key in Redis",
			slog.String("op", op),
			slog.String("key", key),
			slog.String("error", err.Error()),
		)
	}
	return err
}

func (store *Store) HSet(
	ctx context.Context,
	key string,
	values map[string]any,
) error {
	const op = "redisStore.HSet"

	_, err := store.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, key, values)
		p.Expire(ctx, key, time.Hour)
		return nil
	})
	if err != nil {
		store.log.Error(
			"Failed to HSet in Redis",
			slog.String("op", op),
			slog.String("key", key),
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (store *Store) HGet(
	ctx context.Context,
	key, field string,
) (string, error) {
	const op = "redisStore.HGet"

	cmd := store.Client.HGet(
		ctx,
		key, field,
	)
	if cmd.Err() != nil {
		store.log.Error(
			"Failed to HGet from Redis",
			slog.String("op", op),
			slog.String("key", key),
			slog.String("field", field),
			slog.String("error", cmd.Err().Error()),
		)
		return "", cmd.Err()
	}

	return cmd.Result()
}

func (store *Store) HGetAll(
	ctx context.Context,
	key string,
) (map[string]string, error) {
	const op = "redisStore.HGetAll"

	cmd := store.Client.HGetAll(ctx, key)
	if cmd.Err() != nil {
		store.log.Error(
			"Failed to HGetAll from Redis",
			slog.String("op", op),
			slog.String("key", key),
			slog.String("error", cmd.Err().Error()),
		)
		return nil, cmd.Err()
	}

	return cmd.Result()
}

func (store *Store) SetEx(
	ctx context.Context,
	key string,
	value any,
	expire time.Duration,
) error {
	const op = "redisStore.SetEx"

	err := store.Client.SetEx(
		ctx,
		key,
		value,
		expire,
	).Err()
	if err != nil {
		store.log.Error(
			"Failed to SetEx in Redis",
			slog.String("op", op),
			slog.String("key", key),
			slog.String("error", err.Error()),
		)
	}
	return err
}
