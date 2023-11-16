package aconfigo

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cristalhq/aconfig"
)

type Loader struct {
	AppName string
	Files   []string

	Config       aconfig.Config
	WalkFn       func(f aconfig.Field) bool
	FileDecoders map[string]aconfig.FileDecoder
}

func NewLoader() *Loader {
	return &Loader{
		FileDecoders: make(map[string]aconfig.FileDecoder),
	}
}

func (l *Loader) WithAppName(appName string) *Loader {
	l.AppName = appName
	return l
}

func (l *Loader) WithFiles(files []string) *Loader {
	l.Files = files
	return l
}

func (l *Loader) WithFile(file string) *Loader {
	l.Files = append(l.Files, file)
	return l
}

func (l *Loader) WithConfig(config aconfig.Config) *Loader {
	l.Config = config
	return l
}

func (l *Loader) WithWalkFn(fn func(f aconfig.Field) bool) *Loader {
	l.WalkFn = fn
	return l
}

func (l *Loader) WithFileDecoder(decoder aconfig.FileDecoder) *Loader {
	fileFormat := fmt.Sprintf(".%s", decoder.Format())
	l.FileDecoders[fileFormat] = decoder
	return l
}

func (l *Loader) For(cfg any) *aconfig.Loader {
	// remove duplicate files
	slices.Sort(l.Files)
	slices.Compact(l.Files)

	c := l.Config

	if c.EnvPrefix == "" && l.AppName != "" {
		c.EnvPrefix = strings.ToUpper(l.AppName)
	} else {
		c.EnvPrefix = strings.TrimSpace(c.EnvPrefix)
	}

	if c.FlagPrefix == "" && l.AppName != "" {
		c.FlagPrefix = strings.ToLower(l.AppName)
	} else {
		c.FlagPrefix = strings.TrimSpace(c.FlagPrefix)
	}

	if len(l.Files) > 0 {
		c.Files = append(c.Files, l.Files...)
	}

	if len(l.FileDecoders) > 0 {
		c.FileDecoders = l.FileDecoders
	}

	loader := aconfig.LoaderFor(cfg, c)

	if l.WalkFn != nil {
		loader.WalkFields(l.WalkFn)
	}

	return loader
}
