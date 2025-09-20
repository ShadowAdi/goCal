package logger

import (
	"context"
	"log/slog"
)

func Error(msg string, args ...any) {
	ErrorLogger.Log(context.Background(), slog.LevelError, msg, args...)
}
