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
