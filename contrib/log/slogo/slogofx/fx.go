package slogofx

import (
	"log/slog"

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

var FxEventLogger = fx.WithLogger(newFxEventLogger)

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

type setupLoggerWrapper struct {
	slog.Logger
}

func newSetupLoggerWrapper(config *logger.Config) (*setupLoggerWrapper, error) {
	logger, err := slogo.New(*config, nil)
	if err != nil {
		return nil, err
	}
	return &setupLoggerWrapper{Logger: logger}, nil
}

func newSlogFxEventLogger(logger slog.Logger) fxevent.Logger {
	slogLogger := &slogofxevent.SlogLogger{Logger: logger}
	slogLogger.UseLogLevel(slog.LevelDebug)
	return slogLogger
}

func newFxEventLogger(logger *setupLoggerWrapper) fxevent.Logger {
	return newSlogFxEventLogger(logger.Logger)
}