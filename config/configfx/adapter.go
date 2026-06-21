package configfx

import (
	"errors"
	"fmt"
	"reflect"

	"go.uber.org/fx"

	"github.com/go-toho/toho/config"
	"github.com/go-toho/toho/pkg/fxtags"
)

// ProvideAppConfig exposes the resolved application config as a concrete type.
func ProvideAppConfig[T any]() fx.Option {
	return provideConfigFromName[T](config.NamedConfigPointerOut)
}

// ProvideConfig exposes a nested config value that Toho named by its Go type.
func ProvideConfig[T any]() fx.Option {
	return provideConfigFromName[T](reflect.TypeFor[T]().String())
}

func provideConfigFromName[T any](name string) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func(c any) (T, error) {
				return ExtractConfig[T](c)
			},
			fx.ParamTags(fxtags.Named(name)),
		),
	)
}

func ExtractConfig[T any](c any) (T, error) {
	var zero T
	if c == nil {
		return zero, errors.New("config is nil")
	}

	if cfg, ok := c.(T); ok {
		return cfg, nil
	}

	if ptr, ok := c.(*T); ok {
		if ptr == nil {
			return zero, errors.New("config pointer is nil")
		}
		return *ptr, nil
	}

	target := reflect.TypeFor[T]()
	val := reflect.ValueOf(c)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return zero, errors.New("config pointer is nil")
		}
		val = val.Elem()
	}

	if val.Type().AssignableTo(target) {
		return val.Interface().(T), nil
	}
	if val.Type().ConvertibleTo(target) {
		return val.Convert(target).Interface().(T), nil
	}

	return zero, fmt.Errorf("config %T cannot be used as %s", c, target)
}
