package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_LoadFromFile(t *testing.T) {
	// Setup test config file
	configContent := `
log_type: test
paginate_secret: test_secret
grpc:
  port: 50051
  timeout: 5s
  reflection: true
http:
  port: 8081
postgres:
  url: postgres://test:test@localhost:5432/test
`
	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	// Test
	cfg := MustLoad(tmpFile.Name())

	assert.Equal(t, "test", cfg.LogType)
	assert.Equal(t, "test_secret", cfg.PaginateSecret)
	assert.Equal(t, 50051, cfg.GRPC.Port)
	assert.Equal(t, 5*time.Second, cfg.GRPC.Timeout)
	assert.True(t, cfg.GRPC.Reflection)
	assert.Equal(t, 8081, cfg.HTTP.Port)
	assert.Equal(t, "postgres://test:test@localhost:5432/test", cfg.Postgres.URL)
}

func TestConfig_LoadFromEnv(t *testing.T) {
	// Setup test environment
	t.Setenv("LOG_TYPE", "env_test")
	t.Setenv("PAGINATE_SECRET", "env_secret")
	t.Setenv("GRPC_PORT", "50052")
	t.Setenv("GRPC_TIMEOUT", "10s")
	t.Setenv("GRPC_REFLECTION", "false")
	t.Setenv("HTTP_PORT", "8082")
	t.Setenv("POSTGRES_URL", "postgres://env:env@localhost:5432/env")

	// Test
	cfg := MustLoad("")

	assert.Equal(t, "env_test", cfg.LogType)
	assert.Equal(t, "env_secret", cfg.PaginateSecret)
	assert.Equal(t, 50052, cfg.GRPC.Port)
	assert.Equal(t, 10*time.Second, cfg.GRPC.Timeout)
	assert.False(t, cfg.GRPC.Reflection)
	assert.Equal(t, 8082, cfg.HTTP.Port)
	assert.Equal(t, "postgres://env:env@localhost:5432/env", cfg.Postgres.URL)
}

func TestConfig_DefaultValues(t *testing.T) {
	// Clear relevant env vars
	t.Setenv("LOG_TYPE", "")
	t.Setenv("HTTP_PORT", "")

	// Test
	cfg := MustLoad("")

	assert.Equal(t, "dev", cfg.LogType)
	assert.Equal(t, 8081, cfg.HTTP.Port)
}

func TestConfig_RequiredFields(t *testing.T) {
	// Clear required env vars
	t.Setenv("PAGINATE_SECRET", "")
	t.Setenv("GRPC_PORT", "")
	t.Setenv("POSTGRES_URL", "")

	// Test
	assert.Panics(t, func() {
		MustLoad("")
	})
}
