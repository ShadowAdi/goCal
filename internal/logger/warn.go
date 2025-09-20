package logger

import (
	"context"
	"log/slog"
)

func Warn(msg string, args ...any) {
	WarningLogger.Log(context.Background(), slog.LevelDebug, msg, args...)
}
