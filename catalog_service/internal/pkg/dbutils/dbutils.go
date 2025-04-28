package dbutils

import (
	"errors"
	"fmt"
	"reflect"
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
