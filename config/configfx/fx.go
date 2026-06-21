package configfx

import (
	"reflect"
	"strings"

	"go.uber.org/fx"

	"github.com/go-toho/toho/config"
	"github.com/go-toho/toho/pkg/fxtags"
)

const configStructSuffix = "config"

var Module = fx.Module("config",
	configStructCheck,
	ensureConfigOutOption,
)

var (
	configStructCheck = fx.Invoke(
		fx.Annotate(
			func(cfg any) error {
				return config.StructCheck(cfg)
			},
			fx.ParamTags(fxtags.NamedOptional(config.NamedConfigPointerIn)),
		),
	)

	ensureConfigOutOption = fx.Invoke(
		fx.Annotate(
			func(_ any) {},
			fx.ParamTags(fxtags.Named(config.NamedConfigPointerOut)),
		),
	)
)

func SupplyConfigPointer(c any) fx.Option {
	var fxopts []fx.Option

	fxopts = append(fxopts, provideConfigPointer(c))
	fxopts = append(fxopts, structInnerConfigProviders(c)...)

	return fx.Options(fxopts...)
}

func SupplyConfigFile(filename string) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() string { return filename },
			fx.ResultTags(fxtags.Group(config.GroupConfigFiles)),
		),
	)
}

func SupplyConfigFiles(files []string) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() []string { return files },
			fx.ResultTags(fxtags.GroupFlatten(config.GroupConfigFiles)),
		),
	)
}

func SupplyConfig(cfg any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() any { return cfg },
			fx.ResultTags(fxtags.Named(config.NamedConfigPointerOut)),
		),
	)
}

func provideConfigPointer(c any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() any { return c },
			fx.ResultTags(fxtags.Named(config.NamedConfigPointerIn)),
		),
	)
}

func structInnerConfigProviders(structure any) []fx.Option {
	var opts []fx.Option
	inputType := reflect.TypeOf(structure)
	if inputType != nil {
		if inputType.Kind() == reflect.Pointer {
			if inputType.Elem().Kind() == reflect.Struct {
				inputValue := reflect.ValueOf(structure)
				if inputValue.IsNil() {
					return opts
				}

				return appendConfigStructProvider(opts, inputValue.Elem(), config.NamedConfigPointerOut)
			}
		}
	}
	return opts
}

func appendConfigStructProvider(opts []fx.Option, s reflect.Value, paramName string) []fx.Option {
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldType := s.Type().Field(i).Type
		typeName := fieldType.String()

		if fieldType.Kind() == reflect.Struct {
			if strings.HasSuffix(strings.ToLower(typeName), configStructSuffix) {
				if !field.CanInterface() {
					continue
				}

				opts = append(opts, innerConfigStructProvider(paramName, typeName, i))
				opts = appendConfigStructProvider(opts, field, typeName)
			}
		} else if fieldType.Kind() == reflect.Pointer {
			if !field.CanInterface() {
				continue
			}
			if !field.IsZero() && field.Elem().Type().Kind() == reflect.Struct {
				opts = appendConfigStructProvider(opts, field.Elem(), paramName)
				continue
			}
		}
	}

	return opts
}

func innerConfigStructProvider(inName, outName string, fieldIndex int) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func(c any) any {
				if c == nil {
					return nil
				}
				s := reflect.ValueOf(c)
				if s.Kind() == reflect.Pointer {
					if !s.IsZero() && s.Elem().Type().Kind() == reflect.Struct {
						s = s.Elem()
					} else {
						return nil
					}
				}
				return s.Field(fieldIndex).Interface()
			},
			fx.ParamTags(fxtags.Named(inName)),
			fx.ResultTags(fxtags.Named(outName)),
		),
	)
}
