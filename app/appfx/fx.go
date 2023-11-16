package appfx

import (
	"go.uber.org/fx"

	"github.com/go-toho/toho/app"
	"github.com/go-toho/toho/pkg/fxoption"
)

func ProvideApp(a *app.App) fx.Option {
	var fxopts []fx.Option

	fxopts = append(fxopts, fx.Provide(func() app.Info { return a }))
	fxopts = fxoption.SafeAppend(fxopts, fxoption.NameAnnotatedAny(app.NamedAppID, a.ID()))
	fxopts = fxoption.SafeAppend(fxopts, fxoption.NameAnnotatedAny(app.NamedAppName, a.Name()))
	fxopts = fxoption.SafeAppend(fxopts, fxoption.NameAnnotatedAny(app.NamedAppVersion, a.Version()))
	fxopts = fxoption.SafeAppend(fxopts, fxoption.NameAnnotatedAny(app.NamedAppMetadata, a.Metadata()))
	fxopts = fxoption.SafeAppend(fxopts, fxoption.NameAnnotatedAny(app.NamedAppEndpoint, a.Endpoint()))

	return fx.Options(fxopts...)
}

func ProvideAppInfo(opts ...app.Option) fx.Option {
	return ProvideApp(app.New(opts...))
}
