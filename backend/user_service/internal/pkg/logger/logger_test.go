package logger_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/DENFNC/Zappy/user_service/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger_DevMode(t *testing.T) {
	// Проверяем создание логгера в режиме разработки
	log, err := logger.New("dev")
	require.NoError(t, err)
	assert.NotNil(t, log)

	// Перехватываем вывод для проверки форматирования
	var buf bytes.Buffer
	log.Handler().(*logger.PrettyHandler).Handler = slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})

	log.Debug("test debug")
	log.Info("test info")
	log.Warn("test warn")
	log.Error("test error")

	output := buf.String()
	assert.Contains(t, output, "test debug")
	assert.Contains(t, output, "test info")
	assert.Contains(t, output, "test warn")
	assert.Contains(t, output, "test error")
}

func TestNewLogger_ProdMode(t *testing.T) {
	// Проверяем создание логгера в продакшн режиме
	log, err := logger.New("prod")
	require.NoError(t, err)
	assert.NotNil(t, log)

	// Проверяем что используется JSONHandler
	_, ok := log.Handler().(*slog.JSONHandler)
	assert.True(t, ok)
}

func TestNewLogger_InvalidMode(t *testing.T) {
	// Проверяем обработку невалидного режима логирования
	_, err := logger.New("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid log type")
}

func TestNewLogger_EmptyMode(t *testing.T) {
	// Проверяем обработку пустого режима логирования
	_, err := logger.New("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "log type is empty")
}

func TestPrettyHandler_WithAttrs(t *testing.T) {
	// Проверяем добавление атрибутов в PrettyHandler
	var buf bytes.Buffer
	h := &logger.PrettyHandler{
		Handler: slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}),
	}

	log := slog.New(h.WithAttrs([]slog.Attr{slog.String("key", "value")}))
	log.Info("test")

	assert.Contains(t, buf.String(), "key")
	assert.Contains(t, buf.String(), "value")
}
