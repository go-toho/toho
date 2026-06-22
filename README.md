# toho

Toho is a small Go library for application lifecycle structure.

Use the root package when you want simple start/stop hooks without a dependency
injection container. Use `tohofx` when an application benefits from Uber Fx
modules, lifecycle hooks, and dependency graph composition.

Fx mode can populate typed config and logger values exposed through `Config()`
and `Logger()`.

## Install

```sh
go get github.com/go-toho/toho
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"

	"github.com/go-toho/toho"
	"github.com/go-toho/toho/app"
)

func main() {
	app := toho.New(
		toho.AppInfo(app.Name("hello")),
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
```

## Modes

| Mode | Import path | Use when |
| --- | --- | --- |
| Minimal core | `github.com/go-toho/toho` | You need lightweight lifecycle hooks and a small application shell. |
| Fx core | `github.com/go-toho/toho/tohofx` | You want Fx modules, populated config/logger values, and dependency graph composition. |

## Examples

- `examples/minimal`
- `examples/fx`
- `examples/fx-full`

## Development

```sh
go test ./...
```
