package debugfx

import (
	"go.uber.org/fx"

	"github.com/go-toho/toho/tohofx"
)

func init() {
	tohofx.Add("debug", func() fx.Option {
		return Module
	})
}
