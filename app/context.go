package app

import "context"

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, v Info) context.Context {
	return context.WithValue(ctx, appKey{}, v)
}

// FromContext returns the value stored in ctx, if any.
func FromContext(ctx context.Context) (v Info, ok bool) {
	v, ok = ctx.Value(appKey{}).(Info)
	return
}
