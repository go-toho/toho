package slogo

import (
	"log/slog"
	"os"

	"github.com/go-toho/toho/logger"
)

func New(config logger.Config, handlers ...slog.Handler) (slog.Logger, error) {
	if len(handlers) == 0 || handlers[0] == nil {
		handler, err := NewHandler(config)
		if err != nil {
			return *slog.Default(), err
		}
		return *slog.New(handler), nil
	}
	return *slog.New(handlers[0]), nil
}

func NewHandler(config logger.Config) (slog.Handler, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{
		AddSource: config.Caller,
		Level:     level,
	}

	if config.Format == logger.TextFormat {
		return slog.NewTextHandler(os.Stdout, opts), nil
	}

	return slog.NewJSONHandler(os.Stdout, opts), nil
}
