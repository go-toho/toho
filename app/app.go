package app

type Info interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

type App struct {
	opts options
}

func New(opts ...Option) *App {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &App{
		opts: o,
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.opts.id }

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *App) Metadata() map[string]string { return a.opts.metadata }

// Endpoint returns endpoints.
func (a *App) Endpoint() []string {
	if len(a.opts.endpoints) > 0 {
		endpoints := make([]string, 0, len(a.opts.endpoints))
		for _, e := range a.opts.endpoints {
			endpoints = append(endpoints, e.String())
		}
		return endpoints
	}
	return nil
}
