package dbutils

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Scanner interface {
	Scan(dest ...any) error
}

func ScanStruct(a Scanner, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("dest must be a non-nil pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("dest must point to a struct")
	}

	var addrs []any
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		addrs = append(addrs, v.Field(i).Addr().Interface())
	}

	if err := a.Scan(addrs...); err != nil {
		return fmt.Errorf("scan into struct failed: %w", err)
	}
	return nil
}

func WithTx(ctx context.Context, conn *pgxpool.Conn, f func(pgx.Tx) error) error {
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

func NewUUIDV7() uuid.UUID {
	uuid, _ := uuid.NewV7()
	return uuid
}
