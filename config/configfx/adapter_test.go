package configfx_test

import (
	"strings"
	"testing"

	"go.uber.org/fx"

	"github.com/go-toho/toho/config/configfx"
)

type extractSourceConfig struct {
	Name string
}

type extractTargetConfig extractSourceConfig

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

func TestExtractConfigReturnsDirectValue(t *testing.T) {
	want := extractSourceConfig{Name: "toho"}

	got, err := configfx.ExtractConfig[extractSourceConfig](want)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("extracted config = %+v, want %+v", got, want)
	}
}

func TestExtractConfigDereferencesPointer(t *testing.T) {
	want := extractSourceConfig{Name: "toho"}

	got, err := configfx.ExtractConfig[extractSourceConfig](&want)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("extracted config = %+v, want %+v", got, want)
	}
}

func TestExtractConfigRejectsNilInput(t *testing.T) {
	_, err := configfx.ExtractConfig[extractSourceConfig](nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestExtractConfigRejectsNilPointer(t *testing.T) {
	var cfg *extractSourceConfig

	_, err := configfx.ExtractConfig[extractSourceConfig](cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestExtractConfigAcceptsConvertibleValue(t *testing.T) {
	source := extractSourceConfig{Name: "toho"}

	got, err := configfx.ExtractConfig[extractTargetConfig](source)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != source.Name {
		t.Fatalf("extracted name = %q, want %q", got.Name, source.Name)
	}
}

func TestExtractConfigRejectsIncompatibleValue(t *testing.T) {
	_, err := configfx.ExtractConfig[extractSourceConfig]("toho")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "cannot be used") {
		t.Fatalf("expected incompatible config error, got: %v", err)
	}
}

func TestProvideConfigProvidesNestedConfig(t *testing.T) {
	root := &RootConfig{
		Database: DatabaseConfig{DSN: "postgres://example"},
		HTTP:     HTTPConfig{Addr: ":8080"},
	}

	var got DatabaseConfig
	app := fx.New(
		fx.NopLogger,
		configfx.SupplyConfigPointer(root),
		configfx.SupplyConfig(root),
		configfx.ProvideConfig[DatabaseConfig](),
		fx.Invoke(func(cfg DatabaseConfig) {
			got = cfg
		}),
	)
	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
	if got.DSN != root.Database.DSN {
		t.Fatalf("provided DSN = %q, want %q", got.DSN, root.Database.DSN)
	}
}

func TestProvideAppConfigProvidesResolvedAppConfig(t *testing.T) {
	root := &RootConfig{HTTP: HTTPConfig{Addr: ":8080"}}

	var got RootConfig
	app := fx.New(
		fx.NopLogger,
		configfx.SupplyConfig(root),
		configfx.ProvideAppConfig[RootConfig](),
		fx.Invoke(func(cfg RootConfig) {
			got = cfg
		}),
	)
	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
	if got.HTTP.Addr != root.HTTP.Addr {
		t.Fatalf("provided addr = %q, want %q", got.HTTP.Addr, root.HTTP.Addr)
	}
}
