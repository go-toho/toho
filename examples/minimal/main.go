package main

import (
	"context"
	"fmt"

	"github.com/go-toho/toho"
)

func main() {
	app := toho.New(
		toho.BeforeStart(func(context.Context) error {
			fmt.Println("starting")
			return nil
		}),
		toho.AfterStart(func(context.Context) error {
			fmt.Println("started")
			return nil
		}),
	)

	if err := app.Start(); err != nil {
		panic(err)
	}
	if err := app.Stop(); err != nil {
		panic(err)
	}
}
