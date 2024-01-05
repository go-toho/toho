package toho

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/go-toho/toho/app"
)

// CoreOptions struct holds the configuration options for the Core interface.
type CoreOptions struct {
	App app.App

	ConfigPointer any
	LogPointer    any

	Options []any

	StartTimeout time.Duration
	StopTimeout  time.Duration
}

// Core interface defines the methods for initializing and starting
// the application.
type Core interface {
	// Init is used to configure application core.
	Init(*CoreOptions) error

	// Start starts the application.
	Start(context.Context) error

	// Stop gracefully stops the application.
	Stop(context.Context) error

	// Wait for termination of interrupt signals.
	Wait() <-chan os.Signal
}

// defaultCore struct is the default implementation of the Core interface.
type defaultCore struct{}

// verify that defaultCore implements the Core interface.
var _ Core = (*defaultCore)(nil)

func (defaultCore) Init(opts *CoreOptions) error {
	if opts.ConfigPointer != nil {
		if _, ok := opts.ConfigPointer.(*struct{}); !ok {
			s := reflect.ValueOf(opts.ConfigPointer)
			return fmt.Errorf("unsupported config type: %s", s.Type().String())
		}
	}

	return nil
}

func (defaultCore) Start(ctx context.Context) error {
	return ctx.Err()
}

func (defaultCore) Stop(ctx context.Context) error {
	return ctx.Err()
}

func (defaultCore) Wait() <-chan os.Signal {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	return ch
}
