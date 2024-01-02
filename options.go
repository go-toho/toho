package toho

import (
	"context"
	"time"

	"github.com/go-toho/toho/app"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	appInfoOpts []app.Option

	ctx context.Context

	core    Core
	logger  any
	options []any

	startTimeout time.Duration
	stopTimeout  time.Duration

	// Lifecycle functions
	beforeStart []func(context.Context) error
	afterStart  []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// AppInfo with app info.
func AppInfo(opts ...app.Option) Option {
	return func(o *options) { o.appInfoOpts = opts }
}

// AppCore with app core.
func AppCore(s Core) Option {
	return func(o *options) { o.core = s }
}

// Context with app context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Logger with app logger.
func Logger(l any) Option {
	return func(o *options) { o.logger = l }
}

// Options with any options for the app.
func Options(opts ...any) Option {
	return func(o *options) { o.options = opts }
}

// StartTimeout with app start timeout.
func StartTimeout(t time.Duration) Option {
	return func(o *options) { o.startTimeout = t }
}

// StopTimeout with app stop timeout.
func StopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

// Lifecycle functions

// BeforeStart run functions before app starts.
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

// AfterStart run functions after app starts.
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

// BeforeStop run functions before app stops.
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

// AfterStop run functions after app stops.
func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}
