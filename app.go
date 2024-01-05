package toho

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/go-toho/toho/app"
	"github.com/go-toho/toho/pkg/xos"
)

// DefaultTimeout is the default timeout for starting or stopping an
// application. It can be configured with the StartTimeout and StopTimeout
// options.
const DefaultTimeout = 15 * time.Second

var (
	errAlreadyStarted = errors.New("already started")
	errNotStarted     = errors.New("not started")
)

// TohoApp defines a application components lifecycle manager.
// [C] defines type of the configuration struct for the application.
// [L] defines type of the logger for the application.
type TohoApp[C, L any] struct {
	opts    options
	ctx     context.Context
	cancel  func()
	appInfo *app.App
	mu      sync.Mutex

	core      Core
	bootstrap bool

	config C
	log    L
}

func New(opts ...Option) *TohoApp[struct{}, *slog.Logger] {
	return NewCL[struct{}, *slog.Logger](opts...)
}

func NewC[C any](opts ...Option) *TohoApp[C, *slog.Logger] {
	return NewCL[C, *slog.Logger](opts...)
}

func NewL[L any](opts ...Option) *TohoApp[struct{}, L] {
	return NewCL[struct{}, L](opts...)
}

func NewCL[C, L any](opts ...Option) *TohoApp[C, L] {
	o := options{
		ctx:          context.Background(),
		core:         &defaultCore{},
		startTimeout: DefaultTimeout,
		stopTimeout:  DefaultTimeout,
	}

	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	return &TohoApp[C, L]{
		opts:    o,
		ctx:     ctx,
		cancel:  cancel,
		appInfo: app.New(o.appInfoOpts...),
		core:    o.core,
	}
}

func (a *TohoApp[C, L]) AppInfo() app.Info {
	return a.appInfo
}

func (a *TohoApp[C, L]) Config() C {
	return a.config
}

func (a *TohoApp[C, L]) Logger() L {
	return a.log
}

// Start executes all OnStart hooks registered with the application's Lifecycle.
func (a *TohoApp[C, L]) Start() error {
	a.mu.Lock()
	if a.bootstrap {
		return errAlreadyStarted
	}

	if a.opts.logger != nil {
		// handle manually configured logger
		if log, ok := a.opts.logger.(L); ok {
			a.log = log
		}
	}

	coreOpts := &CoreOptions{
		App:           *a.appInfo,
		ConfigPointer: &a.config,
		LogPointer:    &a.log,
		Options:       a.opts.options,
		StartTimeout:  a.opts.startTimeout,
		StopTimeout:   a.opts.stopTimeout,
	}

	if err := a.core.Init(coreOpts); err != nil {
		return fmt.Errorf("%s: %w", reflect.TypeOf(a.core), err)
	}

	// fallback logger
	if reflect.ValueOf(a.log).IsNil() {
		if _, ok := any(a.log).(*slog.Logger); ok {
			a.log = any(slog.Default()).(L)
		}
	}

	a.bootstrap = true
	a.mu.Unlock()

	if err := callLifecycleFn(a.ctx, a.opts.beforeStart); err != nil {
		return err
	}

	if err := a.core.Start(a.ctx); err != nil {
		return err
	}

	if err := callLifecycleFn(a.ctx, a.opts.afterStart); err != nil {
		return err
	}

	return nil
}

// Stop gracefully stops the application.
func (a *TohoApp[C, L]) Stop() error {
	a.mu.Lock()
	if !a.bootstrap {
		return errNotStarted
	}
	a.mu.Unlock()

	if err := callLifecycleFn(a.ctx, a.opts.beforeStop); err != nil {
		return err
	}

	if err := a.core.Stop(a.ctx); err != nil {
		return err
	}

	if a.cancel != nil {
		a.cancel()
	}

	if err := callLifecycleFn(a.ctx, a.opts.afterStop); err != nil {
		return err
	}

	return nil
}

// Wait blocks application until termination.
func (a *TohoApp[C, L]) Wait() <-chan error {
	a.mu.Lock()
	defer a.mu.Unlock()

	ch := make(chan error, 1)

	if !a.bootstrap {
		ch <- errNotStarted
		return ch
	}

	go func() {
		signal := <-a.core.Wait()
		ch <- xos.SignalError{Signal: signal}
	}()

	return ch
}

func callLifecycleFn(ctx context.Context, fns []func(context.Context) error) error {
	for _, fn := range fns {
		if fn == nil {
			continue
		}
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}
