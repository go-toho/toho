package aconfigofx

import (
	"flag"

	"github.com/cristalhq/aconfig"
	"go.uber.org/fx"

	"github.com/go-toho/toho/app"
	"github.com/go-toho/toho/config"
	"github.com/go-toho/toho/contrib/config/aconfigo"
	"github.com/go-toho/toho/pkg/fxtags"
)

var Module = fx.Module("aconfig",
	provideLoader,
	provideConfigLoader,
	provideConfigLoaderFlags,
	provideConfig,
)

var (
	provideLoader = fx.Provide(
		fx.Annotate(
			func(
				appName string,
				files []string,
				config aconfig.Config,
				walkFn func(f aconfig.Field) bool,
				fileDecoders []aconfig.FileDecoder,
			) *aconfigo.Loader {
				loader := aconfigo.NewLoader().
					WithAppName(appName).
					WithFiles(files).
					WithConfig(config).
					WithWalkFn(walkFn)

				for _, decoder := range fileDecoders {
					loader = loader.WithFileDecoder(decoder)
				}

				return loader
			},
			fx.ParamTags(
				fxtags.NamedOptional(app.NamedAppName),
				fxtags.Group(config.GroupConfigFiles),
				fxtags.NamedOptional(aconfigo.NamedConfig),
				fxtags.NamedOptional(aconfigo.NamedWalkFn),
				fxtags.Group(aconfigo.GroupFileDecoders),
			),
		),
	)

	provideConfigLoader = fx.Provide(
		fx.Annotate(
			func(loader *aconfigo.Loader, cfg any) *aconfig.Loader {
				return loader.For(cfg)
			},
			fx.ParamTags(
				fxtags.Empty,
				fxtags.Named(config.NamedConfigPointerIn),
			),
		),
	)

	provideConfigLoaderFlags = fx.Provide(
		fx.Annotate(
			func(aloader *aconfig.Loader) *flag.FlagSet {
				return aloader.Flags()
			},
			fx.ResultTags(fxtags.Named(aconfigo.NamedConfigLoaderFlags)),
		),
	)

	provideConfig = fx.Provide(
		fx.Annotate(
			func(aloader *aconfig.Loader, cfg any) (any, error) {
				if err := aloader.Load(); err != nil {
					return nil, err
				}
				return cfg, nil
			},
			fx.ParamTags(
				fxtags.Empty,
				fxtags.Named(config.NamedConfigPointerIn),
			),
			fx.ResultTags(fxtags.Named(config.NamedConfigPointerOut)),
		),
	)
)

func SupplyConfig(config aconfig.Config) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() aconfig.Config {
				return config
			},
			fx.ResultTags(fxtags.Named(aconfigo.NamedConfig)),
		),
	)
}

func SupplyWalkFn(fn func(f aconfig.Field) bool) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() func(f aconfig.Field) bool {
				return fn
			},
			fx.ResultTags(fxtags.Named(aconfigo.NamedWalkFn)),
		),
	)
}

func SupplyFileDecoder(decoder aconfig.FileDecoder) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() aconfig.FileDecoder {
				return decoder
			},
			fx.ResultTags(fxtags.Group(aconfigo.GroupFileDecoders)),
		),
	)
}
