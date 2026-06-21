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
		inputValue := reflect.ValueOf(structure)
		if inputValue.Kind() == reflect.Pointer && inputValue.IsNil() {
			return errors.New("config: struct nil")
		}

		indirectValue := reflect.Indirect(inputValue)
		if indirectValue.Kind() != reflect.Struct {
			return fmt.Errorf("config: expecting struct, got '%s'", indirectValue.Kind().String())
		}
	}

	return nil
}
