package debugfx

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/fx"

	"github.com/go-toho/toho/contrib/core/debug"
	"github.com/go-toho/toho/pkg/fxtags"
)

var Module = fx.Module("debug",
	provideConfigPointer,
	provideConfig,
	invokeServer,
)

var (
	provideConfigPointer = fx.Provide(
		fx.Annotate(
			func(config any) *debug.Config {
				if config != nil {
					switch v := config.(type) {
					case *debug.Config:
						return v
					case debug.Config:
						return &v
					default:
						break
					}
				}
				return &debug.Config{}
			},
			fx.ParamTags(fxtags.NamedOptional(debug.NamedConfig)),
			fx.ResultTags(fxtags.Named(debug.NamedConfig)),
		),
	)

	provideConfig = fx.Provide(
		fx.Annotate(
			func(config *debug.Config) debug.Config {
				return *config
			},
			fx.ParamTags(fxtags.Named(debug.NamedConfig)),
		),
	)

	invokeServer = fx.Invoke(NewDebugServer)
)

func NewDebugServer(
	config debug.Config,
	log fx.Printer,
	lifecycle fx.Lifecycle,
) error {
	if !config.Enabled {
		log.Printf("debug server not enabled")
		return nil
	}

	server, err := debug.NewHTTPServer(config)
	if err != nil {
		return err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Printf("starting debug server",
					"address", fmt.Sprintf("http://%s/debug/pprof/", server.Addr),
				)

				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("unable to start debug server", "err", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("stopping debug server")
			return server.Shutdown(ctx)
		},
	})

	return nil
}
