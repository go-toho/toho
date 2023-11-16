package tohofx

import "go.uber.org/fx"

type Creator func() fx.Option

var Options = map[string]Creator{}

func Add(name string, creator Creator) {
	Options[name] = creator
}
