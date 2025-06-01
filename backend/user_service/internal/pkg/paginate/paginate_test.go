package paginate_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"testing"

	"github.com/DENFNC/Zappy/user_service/internal/pkg/paginate"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testModel struct {
	ID    int
	Name  string
	Value float64
}

func TestPaginator_NewPaginator(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		conn := &pgxpool.Pool{}
		coder := &mockCoder{}
		p, err := paginate.NewPaginator[testModel](conn, goqu.Dialect("postgres"), coder)
		require.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("nil connection", func(t *testing.T) {
		_, err := paginate.NewPaginator[testModel](nil, goqu.Dialect("postgres"), &mockCoder{})
		assert.EqualError(t, err, "db connection is nil")
	})

	t.Run("nil coder", func(t *testing.T) {
		conn := &pgxpool.Pool{}
		_, err := paginate.NewPaginator[testModel](conn, goqu.Dialect("postgres"), nil)
		assert.EqualError(t, err, "coder is nil")
	})
}

func TestPaginator_Paginate(t *testing.T) {
	conn := &pgxpool.Pool{}
	coder := &mockCoder{}
	p, err := paginate.NewPaginator[testModel](conn, goqu.Dialect("postgres"), coder)
	require.NoError(t, err)

	ds := goqu.From("test_table").Select("id", "name", "value")
	p.WithDataset(ds).WithColumns("id", "name").WithLimit(10)

	t.Run("first page", func(t *testing.T) {
		items, token, err := p.Paginate(context.Background(), "")
		require.NoError(t, err)
		assert.NotEmpty(t, items)
		assert.NotEmpty(t, token)
	})

	t.Run("with token", func(t *testing.T) {
		cursorData := []interface{}{100, "test"}
		var buf bytes.Buffer
		require.NoError(t, gob.NewEncoder(&buf).Encode(cursorData))
		encrypted, err := coder.Encrypt(buf.Bytes())
		require.NoError(t, err)

		items, token, err := p.Paginate(context.Background(), encrypted)
		require.NoError(t, err)
		assert.NotEmpty(t, items)
		assert.NotEmpty(t, token)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, _, err := p.Paginate(context.Background(), "invalid")
		assert.Error(t, err)
	})

	t.Run("no columns", func(t *testing.T) {
		p := p.WithColumns()
		_, _, err := p.Paginate(context.Background(), "")
		assert.EqualError(t, err, "no columns specified for cursor")
	})

	t.Run("no dataset", func(t *testing.T) {
		p := p.WithDataset(nil)
		_, _, err := p.Paginate(context.Background(), "")
		assert.EqualError(t, err, "dataset not provided")
	})
}

func TestEncryptor(t *testing.T) {
	key := make([]byte, 32)
	enc, err := paginate.NewEncryptor(key, nil)
	require.NoError(t, err)

	t.Run("encrypt/decrypt", func(t *testing.T) {
		data := []byte("test data")
		encrypted, err := enc.Encrypt(data)
		require.NoError(t, err)

		decrypted, err := enc.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, data, decrypted)
	})

	t.Run("invalid key size", func(t *testing.T) {
		_, err := paginate.NewEncryptor([]byte("short"), nil)
		assert.Error(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := enc.Decrypt("invalid")
		assert.Error(t, err)
	})
}

type mockCoder struct{}

func (m *mockCoder) Encrypt(data []byte) (string, error) {
	return "mock_encrypted", nil
}

func (m *mockCoder) Decrypt(token string) ([]byte, error) {
	if token == "invalid" {
		return nil, errors.New("invalid token")
	}
	return []byte("mock_decrypted"), nil
}
