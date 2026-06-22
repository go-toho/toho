package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"go.uber.org/fx"

	"github.com/go-toho/toho"
	"github.com/go-toho/toho/app"
	"github.com/go-toho/toho/config/configfx"
	_ "github.com/go-toho/toho/contrib/log/slogo/slogofx"
	"github.com/go-toho/toho/logger"
	"github.com/go-toho/toho/logger/loggerfx"
	"github.com/go-toho/toho/tohofx"
)

type runtimeConfig struct {
	Environment string          `default:"local"`
	Logger      logger.Config   `default:"{}"`
	Metrics     metricsConfig   `default:"{}"`
	HTTP        httpConfig      `default:"{}"`
	Features    map[string]bool `default:"{\"metrics\":true}"`
}

type metricsConfig struct {
	Namespace string `default:"toho_fx_example"`
}

type httpConfig struct {
	Addr string `default:"127.0.0.1:0"`
}

type metricsBundle struct {
	runs atomic.Uint64
}

func main() {
	cfg := runtimeConfig{
		Environment: "local",
		Logger: logger.Config{
			Level:  "info",
			Format: logger.JSONFormat,
		},
		Metrics: metricsConfig{Namespace: "toho_fx_example"},
		HTTP:    httpConfig{Addr: "127.0.0.1:0"},
		Features: map[string]bool{
			"metrics": true,
		},
	}

	a := toho.NewC[runtimeConfig](
		toho.AppCore(tohofx.NewCore()),
		toho.AppInfo(
			app.Name("orders-api"),
			app.Version("dev"),
			app.Metadata(map[string]string{"component": "example"}),
		),
		toho.Options(
			fx.NopLogger,
			configfx.SupplyConfig(&cfg),
			configfx.ProvideAppConfig[runtimeConfig](),
			loggerfx.SupplyFxSetupConfig(&cfg.Logger),
			fx.Supply(&metricsBundle{}),
			fx.Provide(
				newHTTPServer,
				newHTTPListener,
			),
			fx.Invoke(
				runApplication,
				registerHTTPServer,
				registerExampleShutdown,
			),
		),
	)

	if err := a.Start(); err != nil {
		panic(err)
	}
	<-a.Wait()
	if err := a.Stop(); err != nil {
		panic(err)
	}
}

func runApplication(info app.Info, cfg runtimeConfig, log *slog.Logger, metrics *metricsBundle) {
	metrics.runs.Add(1)
	log.Info(
		"application wired",
		slog.String("app", info.Name()),
		slog.String("version", info.Version()),
		slog.String("environment", cfg.Environment),
		slog.String("http_addr", cfg.HTTP.Addr),
	)

	fmt.Printf(
		"ready app=%s env=%s metrics=%s runs=%d\n",
		info.Name(),
		cfg.Environment,
		cfg.Metrics.Namespace,
		metrics.runs.Load(),
	)
}

func newHTTPServer(cfg runtimeConfig, info app.Info) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "ok app=%s env=%s\n", info.Name(), cfg.Environment)
	})

	return &http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: mux,
	}
}

func newHTTPListener(cfg runtimeConfig) (net.Listener, error) {
	return net.Listen("tcp", cfg.HTTP.Addr)
}

func registerHTTPServer(lifecycle fx.Lifecycle, server *http.Server, listener net.Listener, log *slog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("starting HTTP server", slog.String("addr", listener.Addr().String()))
			go func() {
				if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("HTTP server failed", slog.Any("error", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}

func registerExampleShutdown(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, log *slog.Logger) {
	done := make(chan struct{})

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				timer := time.NewTimer(2 * time.Second)
				defer timer.Stop()

				select {
				case <-timer.C:
					log.Info("requesting application shutdown")
					if err := shutdowner.Shutdown(); err != nil {
						log.Error("application shutdown request failed", slog.Any("error", err))
					}
				case <-done:
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			close(done)
			return nil
		},
	})
}
