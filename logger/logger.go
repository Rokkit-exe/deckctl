package logger

import (
	"log/slog"
	"os"
)

var (
	Logger *slog.Logger
)

func Info(msg string) {
	if Logger == nil {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	Logger.Info(msg)
}

func Error(msg string) {
	if Logger == nil {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	Logger.Error(msg)
}
