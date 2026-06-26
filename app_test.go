package toho_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/fx"

	"github.com/go-toho/toho"
	"github.com/go-toho/toho/app"
)

var errInitFailed = errors.New("init failed")

type fakeCore struct {
	initErr error
	init    func(*toho.CoreOptions)
}

func (c fakeCore) Init(opts *toho.CoreOptions) error {
	if c.init != nil {
		c.init(opts)
	}
	return c.initErr
}

func (fakeCore) Start(context.Context) error {
	return nil
}

func (fakeCore) Stop(context.Context) error {
	return nil
}

func (fakeCore) Wait() <-chan os.Signal {
	return make(chan os.Signal)
}

func TestApp(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	app := toho.New(
		toho.AppInfo(
			app.Name("toho"),
			app.Version("1.0.0"),
		),
		toho.Context(ctx),
		toho.Options(
			fx.NopLogger,
		),
		toho.BeforeStart(func(_ context.Context) error {
			t.Log("BeforeStart")
			return nil
		}),
		toho.AfterStart(func(_ context.Context) error {
			t.Log("AfterStart")
			return nil
		}),
		toho.BeforeStop(func(_ context.Context) error {
			t.Log("BeforeStop")
			return nil
		}),
		toho.AfterStop(func(_ context.Context) error {
			t.Log("AfterStop")
			cancel()
			return nil
		}),
	)

	time.AfterFunc(time.Second, func() {
		_ = app.Stop()
	})

	if err := app.Start(); err != nil {
		t.Fatal(err)
	}

	<-ctx.Done() // wait for stop

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Error(ctx.Err())
	}
}

func TestStartUnlocksAfterInitError(t *testing.T) {
	app := toho.New(toho.AppCore(fakeCore{initErr: errInitFailed}))

	if err := app.Start(); !errors.Is(err, errInitFailed) {
		t.Fatalf("Start() error = %v, want %v", err, errInitFailed)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Start()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, errInitFailed) {
			t.Fatalf("second Start() error = %v, want %v", err, errInitFailed)
		}
	case <-time.After(150 * time.Millisecond):
		t.Fatal("second Start() timed out; mutex was not released after init error")
	}
}

type valueLogger struct{}

type testConfig struct {
	Name string
}

func TestConfigExposesBackingConfig(t *testing.T) {
	app := toho.NewC[testConfig](
		toho.AppCore(fakeCore{
			init: func(opts *toho.CoreOptions) {
				cfg, ok := opts.ConfigPointer.(*testConfig)
				if !ok {
					t.Fatalf("ConfigPointer = %T, want *testConfig", opts.ConfigPointer)
				}
				cfg.Name = "loaded"
			},
		}),
	)

	if got := app.Config().Name; got != "" {
		t.Fatalf("Config().Name = %q, want empty", got)
	}

	if err := app.Start(); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	if got := app.Config().Name; got != "loaded" {
		t.Fatalf("Config().Name = %q, want loaded", got)
	}
}

func TestConfigOptionUsesExternalBackingConfig(t *testing.T) {
	cfg := &testConfig{Name: "initial"}
	app := toho.NewC[testConfig](
		toho.Config(cfg),
		toho.AppCore(fakeCore{
			init: func(opts *toho.CoreOptions) {
				got, ok := opts.ConfigPointer.(*testConfig)
				if !ok {
					t.Fatalf("ConfigPointer = %T, want *testConfig", opts.ConfigPointer)
				}
				if got != cfg {
					t.Fatalf("ConfigPointer = %p, want %p", got, cfg)
				}
				cfg.Name = "loaded"
			},
		}),
	)

	if got := app.Config().Name; got != "initial" {
		t.Fatalf("Config().Name = %q, want initial", got)
	}

	if err := app.Start(); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	if got := app.Config().Name; got != "loaded" {
		t.Fatalf("Config().Name = %q, want loaded", got)
	}
	if got := cfg.Name; got != "loaded" {
		t.Fatalf("external config Name = %q, want loaded", got)
	}
}

type otherConfig struct {
	Name string
}

func TestConfigOptionRejectsWrongType(t *testing.T) {
	app := toho.NewC[testConfig](
		toho.Config(&otherConfig{}),
		toho.AppCore(fakeCore{}),
	)

	err := app.Start()
	if err == nil {
		t.Fatal("Start() error = nil, want config type error")
	}
	if !strings.Contains(err.Error(), "config type") {
		t.Fatalf("Start() error = %q, want config type error", err)
	}
}

func TestConfigOptionRejectsNilPointer(t *testing.T) {
	var cfg *testConfig
	app := toho.NewC[testConfig](
		toho.Config(cfg),
		toho.AppCore(fakeCore{}),
	)

	err := app.Start()
	if err == nil {
		t.Fatal("Start() error = nil, want nil config pointer error")
	}
	if !strings.Contains(err.Error(), "config pointer is nil") {
		t.Fatalf("Start() error = %q, want nil config pointer error", err)
	}
}

func TestStartDoesNotPanicWithValueLogger(t *testing.T) {
	app := toho.NewL[valueLogger](toho.AppCore(fakeCore{}))

	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("Start() panicked: %v", recovered)
		}
	}()

	if err := app.Start(); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
}

func TestLifecycleStopUnlocksAfterErrNotStarted(t *testing.T) {
	app := toho.New(toho.AppCore(fakeCore{}))

	if err := app.Stop(); err == nil {
		t.Fatal("Stop() error = nil, want non-nil")
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Stop()
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("second Stop() error = nil, want non-nil")
		}
	case <-time.After(150 * time.Millisecond):
		t.Fatal("second Stop() timed out; mutex was not released after not-started error")
	}
}

func TestLifecycleHooksRunInOrder(t *testing.T) {
	var calls []string
	app := toho.New(
		toho.AppCore(fakeCore{}),
		toho.BeforeStart(func(context.Context) error {
			calls = append(calls, "beforeStart")
			return nil
		}),
		toho.AfterStart(func(context.Context) error {
			calls = append(calls, "afterStart")
			return nil
		}),
		toho.BeforeStop(func(context.Context) error {
			calls = append(calls, "beforeStop")
			return nil
		}),
		toho.AfterStop(func(context.Context) error {
			calls = append(calls, "afterStop")
			return nil
		}),
	)

	if err := app.Start(); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}
	if err := app.Stop(); err != nil {
		t.Fatalf("Stop() error = %v, want nil", err)
	}

	want := []string{"beforeStart", "afterStart", "beforeStop", "afterStop"}
	if len(calls) != len(want) {
		t.Fatalf("lifecycle calls = %v, want %v", calls, want)
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("lifecycle calls = %v, want %v", calls, want)
		}
	}
}
