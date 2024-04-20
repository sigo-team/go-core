package _slg

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func New(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		log = slog.New(handler)
	case envProd:
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
		log = slog.New(handler)
	}

	return log
}

func Err(err error) slog.Attr {
	return slog.String("err", err.Error())
}
