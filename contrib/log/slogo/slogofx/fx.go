package slogofx

import (
	"log/slog"
	"reflect"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/go-toho/toho/contrib/log/slogo"
	slogofxevent "github.com/go-toho/toho/contrib/log/slogo/fxevent"
	"github.com/go-toho/toho/logger"
	"github.com/go-toho/toho/pkg/fxtags"
)

var Module = fx.Module("slog",
	invokeHandlersCountCheck,
	provideDefaultHandler,
	provideLogger,
	provideFxEventLogger,
)

var FxPrinterLogger = fx.Provide(newLoggerPrinter)

var FxEventLogger = fx.WithLogger(newFxEventLogger)

var TrimDefaultHandler = trimDefaultHandler

var SetAsDefaultLogger = fx.Invoke(setAsDefaultLogger)

var (
	invokeHandlersCountCheck = fx.Invoke(
		fx.Annotate(
			func(handlers []slog.Handler) {
				if len(handlers) > 1 {
					slog.Warn("found more that one slog handlers",
						slog.Int("count", len(handlers)),
					)
				}
			},
			fx.ParamTags(fxtags.Group(slogo.GroupSlogHandler)),
		),
	)

	provideDefaultHandler = fx.Provide(
		fx.Annotate(
			func(config logger.Config) (slog.Handler, error) {
				return slogo.NewHandler(config)
			},
			fx.ResultTags(fxtags.Group(slogo.GroupSlogHandler)),
		),
	)

	trimDefaultHandler = fx.Decorate(
		fx.Annotate(
			func(handlers []slog.Handler) []slog.Handler {
				if len(handlers) > 1 {
					for _, h := range handlers {
						if !strings.Contains(reflect.ValueOf(h).Type().String(), "*slog.") {
							return []slog.Handler{h}
						}
					}
				}
				return handlers
			},
			fx.ParamTags(fxtags.Group(slogo.GroupSlogHandler)),
			fx.ResultTags(fxtags.Group(slogo.GroupSlogHandler)),
		),
	)

	provideLogger = fx.Provide(
		fx.Annotate(
			slogo.New,
			fx.ParamTags(
				fxtags.Empty,
				fxtags.Group(slogo.GroupSlogHandler),
			),
		),
	)

	provideFxEventLogger = fx.Provide(
		fx.Annotate(
			newSetupLoggerWrapper,
			fx.ParamTags(fxtags.Named(logger.NamedFxSetupConfig)),
		),
	)
)

func setAsDefaultLogger(logger *slog.Logger) {
	slog.SetDefault(logger)
}

type loggerPrinter struct {
	l *slog.Logger
}

func newLoggerPrinter(logger *slog.Logger) fx.Printer {
	return loggerPrinter{l: logger}
}

func (p loggerPrinter) Printf(msg string, args ...interface{}) {
	log := p.l.Info
	for i := 0; i < len(args); i = i + 2 {
		if k, ok := args[i].(string); ok && k == "error" {
			log = p.l.Error
			break
		}
	}
	log(msg, args...)
}

type setupLoggerWrapper struct {
	*slog.Logger
}

func newSetupLoggerWrapper(config *logger.Config) (*setupLoggerWrapper, error) {
	logger, err := slogo.New(*config, nil)
	if err != nil {
		return nil, err
	}
	return &setupLoggerWrapper{Logger: logger}, nil
}

func newSlogFxEventLogger(logger *slog.Logger) fxevent.Logger {
	slogLogger := &slogofxevent.SlogLogger{Logger: logger}
	slogLogger.UseLogLevel(slog.LevelDebug)
	return slogLogger
}

func newFxEventLogger(logger *setupLoggerWrapper) fxevent.Logger {
	return newSlogFxEventLogger(logger.Logger)
}
