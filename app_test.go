package toho_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/fx"

	"github.com/go-toho/toho"
	"github.com/go-toho/toho/app"
)

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
