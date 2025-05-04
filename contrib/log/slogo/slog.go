package slogo

import (
	"log/slog"
	"os"

	"github.com/go-toho/toho/logger"
)

const (
	// nameKey is used to log the `WithName` values as an additional attribute.
	nameKey = "logger"

	// errKey is used to log the error parameter of Error as an additional attribute.
	errKey = "error"
)

// ParseLevel parses a level based on the number or ASCII representation of the
// log level. If the provided representation is invalid an error is returned.
//
// This is particularly useful when dealing with text input to configure log
// levels.
func ParseLevel(level string) (slog.Level, error) {
	var l slog.Level
	err := l.UnmarshalText([]byte(level))
	return l, err
}

func NewDefault(handlers ...slog.Handler) (*slog.Logger, error) {
	return New(*logger.DefaultConfig, handlers...)
}

func New(config logger.Config, handlers ...slog.Handler) (*slog.Logger, error) {
	if len(handlers) == 0 || handlers[0] == nil {
		handler, err := NewHandler(config)
		if err != nil {
			return slog.Default(), err
		}
		return slog.New(handler), nil
	}
	return slog.New(handlers[0]), nil
}

func NewHandler(config logger.Config) (slog.Handler, error) {
	level, err := ParseLevel(config.Level)
	if err != nil {
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

// WithName returns a new Logger instance with the specified name element added
// to the Logger's name.  Successive calls with WithName result in duplicate
// name attributes, and should be avoided.  It's strongly recommended that name
// segments contain only letters, digits, and hyphens.
func WithName(l *slog.Logger, name string) *slog.Logger {
	return l.With(slog.String(nameKey, name))
}

// Name returns an slog.Attr for the specified name.
func Name(name string) slog.Attr {
	return slog.String(nameKey, name)
}

// Err returns an slog.Attr for the supplied error.
func Err(err error) slog.Attr {
	return slog.Any(errKey, err)
}
