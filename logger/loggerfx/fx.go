package loggerfx

import (
	"go.uber.org/fx"

	"github.com/go-toho/toho/logger"
	"github.com/go-toho/toho/pkg/fxtags"
)

var Module = fx.Module("logger",
	provideFxSetupConfigPointer,
	provideConfigPointer,
	provideConfig,
)

var (
	provideFxSetupConfigPointer = fx.Provide(
		fx.Annotate(
			func(config any) *logger.Config {
				return loggerConfigOrDefault(config)
			},
			fx.ParamTags(fxtags.NamedOptional(logger.NamedFxSetupConfig)),
			fx.ResultTags(fxtags.Named(logger.NamedFxSetupConfig)),
		),
	)

	provideConfigPointer = fx.Provide(
		fx.Annotate(
			func(config any) *logger.Config {
				return loggerConfigOrDefault(config)
			},
			fx.ParamTags(fxtags.NamedOptional(logger.NamedConfig)),
			fx.ResultTags(fxtags.Named(logger.NamedConfig)),
		),
	)

	provideConfig = fx.Provide(
		fx.Annotate(
			func(config *logger.Config) logger.Config {
				return *config
			},
			fx.ParamTags(fxtags.Named(logger.NamedConfig)),
		),
		// also provide named version
		fx.Annotate(
			func(config logger.Config) logger.Config {
				return config
			},
			fx.ResultTags(fxtags.Named(logger.NamedConfig)),
		),
	)
)

func loggerConfigOrDefault(config any) *logger.Config {
	if config != nil {
		switch v := config.(type) {
		case *logger.Config:
			return v
		case logger.Config:
			return &v
		default:
			break
		}
	}
	return logger.DefaultConfig
}

func SupplyFxSetupConfig(config *logger.Config) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() any {
				return config
			},
			fx.ResultTags(fxtags.Named(logger.NamedFxSetupConfig)),
		),
	)
}

func SupplyConfig(config *logger.Config) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() any {
				return config
			},
			fx.ResultTags(fxtags.Named(logger.NamedConfig)),
		),
	)
}
