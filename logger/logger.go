package logger

// The log format can either be text or JSON.
const (
	JSONFormat = "json"
	TextFormat = "text"
)

// Config stores the config for the logger.
type Config struct {
	Level  string `default:"info"`
	Format string `default:"json"`
	Caller bool   `default:"false"`
}

var (
	DefaultConfig = &Config{
		Level:  "info",
		Format: JSONFormat,
		Caller: false,
	}

	DebugTextConfig = &Config{
		Level:  "debug",
		Format: TextFormat,
		Caller: true,
	}
)

// WithLevel returns a new config with overridden value.
func (c *Config) WithLevel(level string) *Config {
	c.Level = level
	return c
}

// WithFormat returns a new config with overridden value.
func (c *Config) WithFormatJSON() *Config {
	c.Format = JSONFormat
	return c
}

// WithFormat returns a new config with overridden value.
func (c *Config) WithFormatText() *Config {
	c.Format = TextFormat
	return c
}

// WithCaller returns a new config with overridden value.
func (c *Config) WithCaller(caller bool) *Config {
	c.Caller = caller
	return c
}
