package slogofx

import (
	"go.uber.org/fx"

	"github.com/go-toho/toho/tohofx"
)

func init() {
	tohofx.Add("log/slog", func() fx.Option {
		return Module
	})
}
