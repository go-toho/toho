package aconfigofx

import (
	"go.uber.org/fx"

	"github.com/go-toho/toho/tohofx"
)

func init() {
	tohofx.Add("github.com/cristalhq/aconfig", func() fx.Option {
		return Module
	})
}
