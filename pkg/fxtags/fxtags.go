package fxtags

import "fmt"

const (
	Empty    = ``
	Optional = `optional:"true"`
)

const (
	FormatName         = `name:"%s"`
	FormatNameOptional = FormatName + " " + Optional
	FormatGroup        = `group:"%s"`
	FormatGroupFlatten = `group:"%s,flatten"`
	FormatGroupSoft    = `group:"%s,soft"`
)

func Named(name string) string {
	return fmt.Sprintf(FormatName, name)
}

func NamedOptional(name string) string {
	return fmt.Sprintf(FormatNameOptional, name)
}

func Group(group string) string {
	return fmt.Sprintf(FormatGroup, group)
}

func GroupFlatten(group string) string {
	return fmt.Sprintf(FormatGroupFlatten, group)
}

func GroupSoft(group string) string {
	return fmt.Sprintf(FormatGroupSoft, group)
}
