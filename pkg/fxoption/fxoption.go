package fxoption

import (
	"reflect"

	"go.uber.org/fx"
)

func SafeAppend(opts []fx.Option, option fx.Option) []fx.Option {
	if option == nil {
		return opts
	}
	return append(opts, option)
}

func NameAnnotatedAny[T any](name string, value T) fx.Option {
	if reflect.ValueOf(value).IsZero() {
		return nil
	}
	return fx.Provide(fx.Annotated{
		Name:   name,
		Target: func() T { return value },
	})
}

func GroupAnnotatedAny[T any](group string, value T) fx.Option {
	if reflect.ValueOf(value).IsZero() {
		return nil
	}
	return fx.Provide(fx.Annotated{
		Group:  group,
		Target: func() T { return value },
	})
}
