package logger

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "local", "dev":
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	default:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
