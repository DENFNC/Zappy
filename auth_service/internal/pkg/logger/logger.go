package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	slog.Handler
}

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

	fieldJSON, err := json.MarshalIndent(fields, "", " ")
	if err != nil {
		return err
	}

	var sb strings.Builder
	timeFormat := r.Time.Format("2006-01-02 15:04:05")
	sb.WriteString(fmt.Sprintf("[%s]  ", timeFormat))
	sb.WriteString(levelString)
	sb.WriteString(fmt.Sprintf("%s\n", msgString))
	sb.WriteString(color.RGB(112, 128, 144).Sprint(string(fieldJSON)))

	finalString := sb.String()
	fmt.Println(finalString)
	return nil
}

func New() *slog.Logger { return nil }

func main() {
	var handle PrettyHandler = PrettyHandler{
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}),
	}

	logger := slog.New(&handle)

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn(
		"Warn message",
		slog.String("key1", "value1"),
	)
	logger.Error(
		"Error message",
		slog.String("key1", "value1"),
		slog.String("key2", "value2"),
		slog.String("key3", "value2"),
		slog.String("key4", "value2"),
		slog.String("key5", "value2"),
	)
}
