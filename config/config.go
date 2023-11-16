package config

import (
	"errors"
	"fmt"
	"reflect"
)

func StructCheck(structure any) error {
	if structure == nil {
		return errors.New("config: struct nil")
	}

	inputType := reflect.TypeOf(structure)

	if inputType != nil {
		indirectValue := reflect.Indirect(reflect.ValueOf(structure))
		if indirectValue.Kind() != reflect.Struct {
			return fmt.Errorf("config: expecting struct, got '%s'", indirectValue.Kind().String())
		}
	}

	return nil
}
