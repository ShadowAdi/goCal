package logger

import (
	"context"
	"log/slog"
)

func Info(msg string, args ...any) {
	InfoLogger.Log(context.Background(), slog.LevelInfo, "msg", msg, args)
}
