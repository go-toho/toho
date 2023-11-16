package fxevent

import (
	"context"
	"log/slog"
	"strings"

	"go.uber.org/fx/fxevent"
)

// SlogLogger is an Fx event logger that logs events to Zap.
type SlogLogger struct {
	Logger slog.Logger

	logLevel   slog.Level // default: slog.LevelInfo
	errorLevel *slog.Level
}

var _ fxevent.Logger = (*SlogLogger)(nil)

// UseErrorLevel sets the level of error logs emitted by Fx to level.
func (l *SlogLogger) UseErrorLevel(level slog.Level) {
	l.errorLevel = &level
}

// UseLogLevel sets the level of non-error logs emitted by Fx to level.
func (l *SlogLogger) UseLogLevel(level slog.Level) {
	l.logLevel = level
}

func (l *SlogLogger) logEvent(msg string, fields ...slog.Attr) {
	l.Logger.LogAttrs(context.Background(), l.logLevel, msg, fields...)
}

func (l *SlogLogger) logError(msg string, fields ...slog.Attr) {
	lvl := slog.LevelError
	if l.errorLevel != nil {
		lvl = *l.errorLevel
	}
	l.Logger.LogAttrs(context.Background(), lvl, msg, fields...)
}

// LogEvent logs the given event to the provided Zap logger.
func (l *SlogLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logEvent("OnStart hook executing",
			slog.String("callee", e.FunctionName),
			slog.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logError("OnStart hook failed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.Any("err", e.Err),
			)
		} else {
			l.logEvent("OnStart hook executed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		l.logEvent("OnStop hook executing",
			slog.String("callee", e.FunctionName),
			slog.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logError("OnStop hook failed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.Any("err", e.Err),
			)
		} else {
			l.logEvent("OnStop hook executed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.logError("error encountered while applying options",
				slog.String("type", e.TypeName),
				slog.Any("stacktrace", e.StackTrace),
				slog.Any("moduletrace", e.ModuleTrace),
				moduleField(e.ModuleName),
				slog.Any("err", e.Err),
			)
		} else {
			l.logEvent("supplied",
				slog.String("type", e.TypeName),
				slog.Any("stacktrace", e.StackTrace),
				slog.Any("moduletrace", e.ModuleTrace),
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("provided",
				slog.String("constructor", e.ConstructorName),
				slog.Any("stacktrace", e.StackTrace),
				slog.Any("moduletrace", e.ModuleTrace),
				moduleField(e.ModuleName),
				slog.String("type", rtype),
				maybeBool("private", e.Private),
			)
		}
		if e.Err != nil {
			l.logError("error encountered while applying options",
				moduleField(e.ModuleName),
				slog.Any("stacktrace", e.StackTrace),
				slog.Any("moduletrace", e.ModuleTrace),
				slog.Any("err", e.Err),
			)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("replaced",
				slog.Any("stacktrace", e.StackTrace),
				slog.Any("moduletrace", e.ModuleTrace),
				moduleField(e.ModuleName),
				slog.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.logError("error encountered while replacing",
				slog.Any("stacktrace", e.StackTrace),
				slog.Any("moduletrace", e.ModuleTrace),
				moduleField(e.ModuleName),
				slog.Any("err", e.Err),
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("decorated",
				slog.String("decorator", e.DecoratorName),
				slog.Any("stacktrace", e.StackTrace),
				slog.Any("moduletrace", e.ModuleTrace),
				moduleField(e.ModuleName),
				slog.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.logError("error encountered while applying options",
				slog.Any("stacktrace", e.StackTrace),
				slog.Any("moduletrace", e.ModuleTrace),
				moduleField(e.ModuleName),
				slog.Any("err", e.Err),
			)
		}
	case *fxevent.Run:
		if e.Err != nil {
			l.logError("error returned",
				slog.String("name", e.Name),
				slog.String("kind", e.Kind),
				moduleField(e.ModuleName),
				slog.Any("err", e.Err),
			)
		} else {
			l.logEvent("run",
				slog.String("name", e.Name),
				slog.String("kind", e.Kind),
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.logEvent("invoking",
			slog.String("function", e.FunctionName),
			moduleField(e.ModuleName),
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			l.logError("invoke failed",
				slog.Any("err", e.Err),
				slog.String("stack", e.Trace),
				slog.String("function", e.FunctionName),
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Stopping:
		l.logEvent("received signal",
			slog.String("signal", strings.ToUpper(e.Signal.String())),
		)
	case *fxevent.Stopped:
		if e.Err != nil {
			l.logError("stop failed",
				slog.Any("err", e.Err),
			)
		}
	case *fxevent.RollingBack:
		l.logError("start failed, rolling back",
			slog.Any("err", e.StartErr),
		)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logError("rollback failed",
				slog.Any("err", e.Err),
			)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.logError("start failed",
				slog.Any("err", e.Err),
			)
		} else {
			l.logEvent("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logError("custom logger initialization failed",
				slog.Any("err", e.Err),
			)
		} else {
			l.logEvent("initialized custom fxevent.Logger",
				slog.String("function", e.ConstructorName),
			)
		}
	}
}

func moduleField(name string) slog.Attr {
	if len(name) == 0 {
		return slog.Attr{}
	}
	return slog.String("module", name)
}

func maybeBool(name string, b bool) slog.Attr {
	if b {
		return slog.Bool(name, true)
	}
	return slog.Attr{}
}
