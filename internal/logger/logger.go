package logger

import (
	"log/slog"
	"os"
)

var (
	InfoLogger    *slog.Logger
	ErrorLogger   *slog.Logger
	WarningLogger *slog.Logger
)

func InitLogger() {
	infoFile, _ := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	warnFile, _ := os.OpenFile("logs/warn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	errorFile, _ := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	InfoLogger = slog.New(slog.NewJSONHandler(infoFile, nil))
	WarningLogger = slog.New(slog.NewJSONHandler(warnFile, nil))
	ErrorLogger = slog.New(slog.NewJSONHandler(errorFile, nil))

	slog.SetDefault(InfoLogger)
}
