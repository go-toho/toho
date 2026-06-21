package configfx

import (
	"strings"
	"testing"

	"go.uber.org/fx"

	"github.com/go-toho/toho/pkg/fxtags"
)

type DatabaseConfig struct {
	DSN string
}

type HTTPConfig struct {
	Addr string
}

type RootConfig struct {
	Database DatabaseConfig
	HTTP     HTTPConfig
}

type PointerRootConfig struct {
	Database *DatabaseConfig
	HTTP     HTTPConfig
}

type hiddenConfig struct {
	Value string
}

type UnexportedRootConfig struct {
	hidden hiddenConfig
	HTTP   HTTPConfig
}

const (
	databaseConfigName = "configfx.DatabaseConfig"
	httpConfigName     = "configfx.HTTPConfig"
	hiddenConfigName   = "configfx.hiddenConfig"
)

func TestSupplyConfigPointerAcceptsPointerToStruct(t *testing.T) {
	app := fx.New(SupplyConfigPointer(&RootConfig{}))
	if err := app.Err(); err != nil {
		t.Fatalf("expected app to build, got error: %v", err)
	}
}

func TestSupplyConfigPointerProvidesNestedConfigFields(t *testing.T) {
	root := &RootConfig{
		Database: DatabaseConfig{DSN: "postgres://toho"},
		HTTP:     HTTPConfig{Addr: ":8080"},
	}

	var database any
	app := fx.New(
		SupplyConfig(root),
		SupplyConfigPointer(root),
		fx.Populate(
			fx.Annotate(
				&database,
				fx.ParamTags(fxtags.Named(databaseConfigName)),
			),
		),
	)
	if err := app.Err(); err != nil {
		t.Fatalf("expected app to build, got error: %v", err)
	}

	got, ok := database.(DatabaseConfig)
	if !ok {
		t.Fatalf("expected DatabaseConfig, got %T", database)
	}
	if got.DSN != root.Database.DSN {
		t.Fatalf("expected DSN %q, got %q", root.Database.DSN, got.DSN)
	}
}

func TestSupplyConfigPointerProvidesSiblingNestedConfigFields(t *testing.T) {
	root := &RootConfig{
		Database: DatabaseConfig{DSN: "postgres://toho"},
		HTTP:     HTTPConfig{Addr: ":8080"},
	}

	var database any
	var httpConfig any
	app := fx.New(
		SupplyConfig(root),
		SupplyConfigPointer(root),
		fx.Populate(
			fx.Annotate(
				&database,
				fx.ParamTags(fxtags.Named(databaseConfigName)),
			),
			fx.Annotate(
				&httpConfig,
				fx.ParamTags(fxtags.Named(httpConfigName)),
			),
		),
	)
	if err := app.Err(); err != nil {
		t.Fatalf("expected app to build, got error: %v", err)
	}

	gotDatabase, ok := database.(DatabaseConfig)
	if !ok {
		t.Fatalf("expected DatabaseConfig, got %T", database)
	}
	if gotDatabase.DSN != root.Database.DSN {
		t.Fatalf("expected DSN %q, got %q", root.Database.DSN, gotDatabase.DSN)
	}

	gotHTTP, ok := httpConfig.(HTTPConfig)
	if !ok {
		t.Fatalf("expected HTTPConfig, got %T", httpConfig)
	}
	if gotHTTP.Addr != root.HTTP.Addr {
		t.Fatalf("expected address %q, got %q", root.HTTP.Addr, gotHTTP.Addr)
	}
}

func TestSupplyConfigPointerNilPointerReturnsValidationError(t *testing.T) {
	var root *RootConfig
	var app *fx.App

	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				t.Fatalf("expected nil config pointer not to panic, recovered: %v", recovered)
			}
		}()
		app = fx.New(SupplyConfigPointer(root), Module)
	}()

	if app == nil {
		t.Fatal("expected app to be created")
	}
	err := app.Err()
	if err == nil {
		t.Fatal("expected config validation error, got nil")
	}
	if !strings.Contains(err.Error(), "config:") {
		t.Fatalf("expected config validation error, got: %v", err)
	}
}

func TestSupplyConfigPointerPointerNestedConfigDoesNotSkipSiblingFields(t *testing.T) {
	root := &PointerRootConfig{
		Database: &DatabaseConfig{DSN: "postgres://toho"},
		HTTP:     HTTPConfig{Addr: ":8080"},
	}

	var httpConfig any
	app := fx.New(
		SupplyConfig(root),
		SupplyConfigPointer(root),
		fx.Populate(
			fx.Annotate(
				&httpConfig,
				fx.ParamTags(fxtags.Named(httpConfigName)),
			),
		),
	)
	if err := app.Err(); err != nil {
		t.Fatalf("expected app to build, got error: %v", err)
	}

	gotHTTP, ok := httpConfig.(HTTPConfig)
	if !ok {
		t.Fatalf("expected HTTPConfig, got %T", httpConfig)
	}
	if gotHTTP.Addr != root.HTTP.Addr {
		t.Fatalf("expected address %q, got %q", root.HTTP.Addr, gotHTTP.Addr)
	}
}

func TestSupplyConfigPointerSkipsUnexportedNestedConfigFields(t *testing.T) {
	root := &UnexportedRootConfig{
		hidden: hiddenConfig{Value: "secret"},
		HTTP:   HTTPConfig{Addr: ":8080"},
	}

	var hidden any
	app := fx.New(
		SupplyConfig(root),
		SupplyConfigPointer(root),
		fx.Populate(
			fx.Annotate(
				&hidden,
				fx.ParamTags(fxtags.Named(hiddenConfigName)),
			),
		),
	)
	err := app.Err()
	if err == nil {
		t.Fatal("expected unexported nested config field not to be provided")
	}
	if strings.Contains(err.Error(), "reflect.Value.Interface") || strings.Contains(err.Error(), "panic") {
		t.Fatalf("expected unexported field to be skipped without panic, got: %v", err)
	}
}
