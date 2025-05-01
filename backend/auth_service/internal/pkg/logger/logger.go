// Package logger предоставляет кастомный логгер, построенный на базе пакета slog.
// В режиме разработки ("dev") используются цветная и удобочитаемая (pretty) раскраска сообщений,
// а в режиме продакшна ("prod") – логирование в формате JSON с информацией об источнике.
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
)

// PrettyHandler расширяет стандартный обработчик slog.Handler для форматированного вывода логов.
// Он оборачивает базовый slog.Handler и переопределяет метод Handle для вывода
// сообщений с цветовой подсветкой в зависимости от уровня лога.
type PrettyHandler struct {
	slog.Handler
}

// Handle обрабатывает запись лога slog.Record.
// Он форматирует время, уровень, сообщение и дополнительные атрибуты.
// В зависимости от уровня (Debug, Info, Warn, Error) применяется соответствующая цветовая подсветка.
// Если присутствуют дополнительные атрибуты, они сериализуются в отформатированный JSON.
// Полученная строка выводится в стандартный вывод.
func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelString := r.Level.String() + " | "
	msgString := r.Message

	switch r.Level {
	case slog.LevelDebug:
		msgString = color.HiMagentaString(msgString)
		levelString = color.HiMagentaString(levelString)
	case slog.LevelInfo:
		msgString = color.HiGreenString(msgString)
		levelString = color.HiGreenString(levelString)
	case slog.LevelWarn:
		msgString = color.HiYellowString(msgString)
		levelString = color.HiYellowString(levelString)
	case slog.LevelError:
		msgString = color.HiRedString(msgString)
		levelString = color.HiRedString(levelString)
	}

	fields := make(map[string]any, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	var fieldJSON []byte
	if len(fields) > 0 {
		// Форматируем атрибуты в отформатированный JSON.
		var err error
		fieldJSON, err = json.MarshalIndent(fields, "", " ")
		if err != nil {
			return err
		}
	}

	var sb strings.Builder
	timeFormat := r.Time.Format("2006-01-02 15:04:05")
	sb.WriteString(fmt.Sprintf("[%s]  ", timeFormat))
	sb.WriteString(levelString)
	sb.WriteString(msgString)
	if len(fields) > 0 {
		sb.WriteString("\n")
		sb.WriteString(color.RGB(112, 128, 144).Sprint(string(fieldJSON)))
	}

	finalString := sb.String()
	fmt.Println(finalString)
	return nil
}

// WithAttrs возвращает новый обработчик с добавленными атрибутами.
// Переданные атрибуты объединяются с базовыми атрибутами обработчика.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

// WithGroup возвращает новый обработчик, который добавляет указанную группу к ключам атрибутов.
// Это позволяет структурировать атрибуты логов по группам.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
	}
}

// New создает новый экземпляр slog.Logger в зависимости от переданного типа логирования.
// Параметр logType должен принимать одно из двух значений: "dev" для режима разработки или "prod" для продакшн-режима.
//
// В режиме "dev":
// - Создается PrettyHandler, который выводит текстовые логи с цветовой подсветкой и уровнем логирования Debug.
// - Для каждого атрибута вызывается ReplaceAttr, объединяющий группы (если присутствуют) с ключом атрибута.
//
// В режиме "prod":
//   - Создается стандартный JSONHandler с дополнительной информацией об источнике вызова
//     и уровнем логирования Info.
//
// Если logType пустой, функция возвращает ошибку.
func New(logType string) (*slog.Logger, error) {
	var handle slog.Handler

	if logType == "" {
		return nil, fmt.Errorf("log type is empty")
	}

	switch logType {
	case "dev":
		handle = &PrettyHandler{
			Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if len(groups) > 0 {
						a.Key = strings.Join(groups, ".") + "." + a.Key
					}
					return a
				},
			}),
		}
	case "prod":
		handle = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
	}

	logger := slog.New(handle)

	return logger, nil
}
