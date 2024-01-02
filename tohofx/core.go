package tohofx

import (
	"context"
	"os"
	"reflect"

	"go.uber.org/fx"

	"github.com/go-toho/toho"
	"github.com/go-toho/toho/app/appfx"
	"github.com/go-toho/toho/config/configfx"
	"github.com/go-toho/toho/logger/loggerfx"
)

var FxCore = &fxCore{}

type fxCore struct {
	opts     *toho.CoreOptions
	instance *fx.App
}

var _ toho.Core = (*fxCore)(nil)

func (s *fxCore) Init(opts *toho.CoreOptions) error {
	s.opts = opts

	var fxOptions []fx.Option

	// with config
	if opts.ConfigPointer != nil {
		if _, ok := opts.ConfigPointer.(*struct{}); !ok {
			fxOptions = append(fxOptions, configfx.Module)
			fxOptions = append(fxOptions, configfx.SupplyConfigPointer(opts.ConfigPointer))
		}
	}

	// with logger
	if opts.LogPointer != nil {
		s := reflect.ValueOf(opts.LogPointer)
		if s.Kind() == reflect.Pointer && reflect.Indirect(s).Kind() == reflect.Pointer {
			s = reflect.Indirect(s)
		}
		if s.IsZero() {
			fxOptions = append(fxOptions, fx.Populate(opts.LogPointer))
		}
	}

	// additional options
	for _, opt := range opts.Options {
		switch v := opt.(type) {
		case fx.Option:
			fxOptions = append(fxOptions, v)
		}
	}

	// setup timeouts
	fxOptions = append(fxOptions, fx.StartTimeout(opts.StartTimeout))
	fxOptions = append(fxOptions, fx.StopTimeout(opts.StopTimeout))

	s.instance = fx.New(
		appfx.ProvideApp(&opts.App),
		loggerfx.Module,
		ProvideRegistered(),
		fx.Options(fxOptions...),
	)

	return s.instance.Err()
}

func (s *fxCore) Start(ctx context.Context) error {
	startCtx, cancel := context.WithTimeout(ctx, s.opts.StartTimeout)
	defer cancel()
	return s.instance.Start(startCtx)
}

func (s *fxCore) Stop(ctx context.Context) error {
	stopCtx, cancel := context.WithTimeout(ctx, s.opts.StopTimeout)
	defer cancel()
	return s.instance.Stop(stopCtx)
}

func (s *fxCore) Wait() os.Signal {
	return <-s.instance.Done()
}
