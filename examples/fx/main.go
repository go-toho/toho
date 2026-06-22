package main

import (
	"github.com/go-toho/toho"
	"github.com/go-toho/toho/logger"
	"github.com/go-toho/toho/logger/loggerfx"
	"github.com/go-toho/toho/tohofx"
)

func main() {
	app := toho.New(
		toho.AppCore(tohofx.NewCore()),
		toho.Options(
			loggerfx.SupplyFxSetupConfig(logger.DebugTextConfig),
		),
	)

	if err := app.Start(); err != nil {
		panic(err)
	}
	if err := app.Stop(); err != nil {
		panic(err)
	}
}
