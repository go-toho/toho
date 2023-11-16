package tohofx

import "go.uber.org/fx"

func ProvideRegistered() fx.Option {
	var fxopts []fx.Option
	for _, opts := range Options {
		fxopts = append(fxopts, opts())
	}
	return fx.Options(fxopts...)
}
